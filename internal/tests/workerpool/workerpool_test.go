package workerpool_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"marketplace/internal/workerpool"
)

func TestSubmitExecutesAllTasks(t *testing.T) {
	p := workerpool.New(3, 100)

	var mu sync.Mutex
	counter := 0
	for i := 0; i < 50; i++ {
		if ok := p.Submit(func() {
			mu.Lock()
			counter++
			mu.Unlock()
		}); !ok {
			t.Fatalf("Submit вернул false, очередь не должна была переполниться")
		}
	}

	if err := p.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}

	if counter != 50 {
		t.Errorf("выполнено %d задач из 50", counter)
	}
}

func TestConcurrencyLimitedByWorkers(t *testing.T) {
	const workers = 2
	p := workerpool.New(workers, 100)

	var cur, max int32
	for i := 0; i < 20; i++ {
		p.Submit(func() {
			c := atomic.AddInt32(&cur, 1)
			for {
				m := atomic.LoadInt32(&max)
				if c <= m || atomic.CompareAndSwapInt32(&max, m, c) {
					break
				}
			}
			time.Sleep(5 * time.Millisecond)
			atomic.AddInt32(&cur, -1)
		})
	}

	if err := p.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}

	if m := atomic.LoadInt32(&max); m > workers {
		t.Errorf("параллелизм достиг %d, максимум %d", m, workers)
	}
}

func TestSubmitReturnsFalseWhenQueueFull(t *testing.T) {
	gate := make(chan struct{})
	p := workerpool.New(1, 1)

	started := make(chan struct{})
	if !p.Submit(func() {
		close(started)
		<-gate
	}) {
		t.Fatal("первая задача должна быть принята")
	}
	<-started

	if !p.Submit(func() {}) {
		t.Fatal("вторая задача должна занять буфер")
	}
	if p.Submit(func() {}) {
		t.Error("Submit должен вернуть false при полной очереди")
	}

	close(gate)
	if err := p.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
}

func TestSubmitBlockingWaitsForSpace(t *testing.T) {
	gate := make(chan struct{})
	p := workerpool.New(1, 1)

	started := make(chan struct{})
	p.Submit(func() {
		close(started)
		<-gate
	})
	<-started

	if !p.Submit(func() {}) {
		t.Fatal("вторая задача должна занять буфер")
	}

	done := make(chan struct{})
	go func() {
		p.SubmitBlocking(func() {})
		close(done)
	}()

	select {
	case <-done:
		t.Fatal("SubmitBlocking не должен выполняться сразу при полной очереди")
	case <-time.After(50 * time.Millisecond):
	}

	close(gate)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("SubmitBlocking так и не выполнился после освобождения очереди")
	}

	if err := p.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
}

func TestShutdownDrainsQueue(t *testing.T) {
	p := workerpool.New(1, 64)

	var mu sync.Mutex
	counter := 0
	for i := 0; i < 10; i++ {
		p.Submit(func() {
			time.Sleep(5 * time.Millisecond)
			mu.Lock()
			counter++
			mu.Unlock()
		})
	}

	if err := p.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if counter != 10 {
		t.Errorf("очередь дренирована не полностью: выполнено %d из 10", counter)
	}
}

func TestShutdownTimeout(t *testing.T) {
	gate := make(chan struct{})
	p := workerpool.New(1, 8)

	started := make(chan struct{})
	p.Submit(func() {
		close(started)
		<-gate
	})
	<-started

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	if err := p.Shutdown(ctx); err == nil {
		t.Error("ожидали ошибку таймаута от Shutdown")
	}

	close(gate)
}
