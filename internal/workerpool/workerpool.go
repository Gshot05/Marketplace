package workerpool

import (
	"context"
	"sync"
)

type Pool struct {
	tasks chan func()
	wg    sync.WaitGroup
}

func New(workers, queueSize int) *Pool {
	if workers < 1 {
		workers = 1
	}
	if queueSize < 0 {
		queueSize = 0
	}

	p := &Pool{
		tasks: make(chan func(), queueSize),
	}

	p.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go p.worker()
	}

	return p
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for task := range p.tasks {
		task()
	}
}

func (p *Pool) Submit(task func()) bool {
	select {
	case p.tasks <- task:
		return true
	default:
		return false
	}
}

func (p *Pool) SubmitBlocking(task func()) {
	p.tasks <- task
}

func (p *Pool) Shutdown(ctx context.Context) error {
	close(p.tasks)

	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
