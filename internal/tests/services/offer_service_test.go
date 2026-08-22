package services_test

import (
	"context"
	errors2 "marketplace/internal/error"
	"marketplace/internal/model"
	"marketplace/internal/service"
	"testing"
)

func TestOfferServiceCreateValidation(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name        string
		title       string
		description string
		price       float64
		wantErr     error
	}{
		{"валидные данные", "Title", "Desc", 100, nil},
		{"пустой заголовок", "", "Desc", 100, errors2.ErrEmptyTitle},
		{"пустое описание", "Title", "", 100, errors2.ErrEmptyDescription},
		{"нулевая цена", "Title", "Desc", 0, errors2.ErrEmptyPrice},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeOfferRepo{}
			s := service.NewOfferService(repo)

			offer, err := s.CreateOffer(ctx, 1, tc.title, tc.description, tc.price)
			if err != tc.wantErr {
				t.Fatalf("ожидали %v, получили %v", tc.wantErr, err)
			}
			if tc.wantErr != nil {
				if offer != nil {
					t.Error("при ошибке offer должен быть nil")
				}
				repo.mu.Lock()
				called := repo.createCalled
				repo.mu.Unlock()
				if called {
					t.Error("репозиторий не должен вызываться при невалидных данных")
				}
				return
			}
			if offer == nil || offer.ID != 1 || offer.CustomerID != 1 {
				t.Errorf("неожиданный результат: %+v", offer)
			}
		})
	}
}

func TestOfferServiceUpdateAndDelete(t *testing.T) {
	ctx := context.Background()

	s := service.NewOfferService(&fakeOfferRepo{})
	if _, err := s.UpdateOffer(ctx, 1, 1, "T", "D", 10); err != nil {
		t.Errorf("UpdateOffer: %v", err)
	}

	repo := &fakeOfferRepo{deleteRes: true}
	s = service.NewOfferService(repo)
	deleted, err := s.DeleteOffer(ctx, 1, 1)
	if err != nil || !deleted {
		t.Errorf("DeleteOffer: deleted=%v err=%v", deleted, err)
	}

	repo = &fakeOfferRepo{deleteRes: false}
	s = service.NewOfferService(repo)
	deleted, err = s.DeleteOffer(ctx, 99, 1)
	if err != nil || deleted {
		t.Errorf("DeleteOffer несуществующего: deleted=%v err=%v", deleted, err)
	}
}

func TestOfferServiceList(t *testing.T) {
	ctx := context.Background()

	empty := &fakeOfferRepo{}
	s := service.NewOfferService(empty)
	if _, err := s.ListOffers(ctx); err != errors2.ErrEmptyOffers {
		t.Errorf("ожидали ErrEmptyOffers, получили %v", err)
	}

	full := &fakeOfferRepo{list: []model.Offer{{ID: 1}}}
	s = service.NewOfferService(full)
	offers, err := s.ListOffers(ctx)
	if err != nil || len(offers) != 1 {
		t.Errorf("ListOffers: len=%d err=%v", len(offers), err)
	}
}

func TestServiceServiceCreateAndList(t *testing.T) {
	ctx := context.Background()
	s := service.NewServiceService(&fakeServiceRepo{})

	if _, err := s.CreateService(ctx, 1, "", "D", 10); err != errors2.ErrEmptyTitle {
		t.Errorf("ожидали ErrEmptyTitle, получили %v", err)
	}

	svc, err := s.CreateService(ctx, 5, "Title", "Desc", 50)
	if err != nil || svc.PerformerID != 5 {
		t.Errorf("CreateService: %+v err=%v", svc, err)
	}

	empty := &fakeServiceRepo{}
	s = service.NewServiceService(empty)
	if _, err := s.ListServices(ctx); err != errors2.ErrEmptyServices {
		t.Errorf("ожидали ErrEmptyServices, получили %v", err)
	}
}

func TestFavoriteService(t *testing.T) {
	ctx := context.Background()

	fav, err := service.NewFavoriteService(&fakeFavoriteRepo{}).AddFavorite(ctx, 1, 2)
	if err != nil || fav.ServiceID != 2 || fav.CustomerID != 1 {
		t.Errorf("AddFavorite: %+v err=%v", fav, err)
	}

	empty := &fakeFavoriteRepo{}
	s := service.NewFavoriteService(empty)
	if _, err := s.ListFavorites(ctx, 1); err != errors2.ErrEmptyFav {
		t.Errorf("ожидали ErrEmptyFav, получили %v", err)
	}

	notFound := &fakeFavoriteRepo{deleteRes: false}
	s = service.NewFavoriteService(notFound)
	deleted, err := s.DeleteFavorite(ctx, 1, 99)
	if err != nil || deleted {
		t.Errorf("DeleteFavorite: deleted=%v err=%v", deleted, err)
	}
}
