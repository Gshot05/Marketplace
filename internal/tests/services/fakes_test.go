package services_test

import (
	"context"
	"marketplace/internal/model"
	"sync"
)

type fakeOfferRepo struct {
	mu        sync.Mutex
	createErr error
	updateRes *model.Offer
	updateErr error
	deleteRes bool
	deleteErr error
	list      []model.Offer
	listErr   error

	createCalled bool
	lastArgs     []interface{}
}

func (f *fakeOfferRepo) Create(ctx context.Context, customerID uint, title, description string, price float64) (*model.Offer, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.createCalled = true
	f.lastArgs = []interface{}{customerID, title, description, price}
	if f.createErr != nil {
		return nil, f.createErr
	}
	return &model.Offer{ID: 1, CustomerID: customerID, Title: title, Description: description, Price: price}, nil
}

func (f *fakeOfferRepo) Update(ctx context.Context, offerID, customerID uint, title, description string, price float64) (*model.Offer, error) {
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	return &model.Offer{ID: offerID, CustomerID: customerID, Title: title, Description: description, Price: price}, nil
}

func (f *fakeOfferRepo) Delete(ctx context.Context, offerID, customerID uint) (bool, error) {
	return f.deleteRes, f.deleteErr
}

func (f *fakeOfferRepo) List(ctx context.Context) ([]model.Offer, error) {
	return f.list, f.listErr
}

type fakeServiceRepo struct {
	deleteRes bool
	deleteErr error
	list      []model.Service
	listErr   error
}

func (f *fakeServiceRepo) Create(ctx context.Context, performerID uint, title, description string, price float64) (*model.Service, error) {
	return &model.Service{ID: 1, PerformerID: performerID, Title: title, Description: description, Price: price}, nil
}

func (f *fakeServiceRepo) Update(ctx context.Context, serviceID, performerID uint, title, description string, price float64) (*model.Service, error) {
	return &model.Service{ID: serviceID, PerformerID: performerID, Title: title, Description: description, Price: price}, nil
}

func (f *fakeServiceRepo) Delete(ctx context.Context, serviceID, performerID uint) (bool, error) {
	return f.deleteRes, f.deleteErr
}

func (f *fakeServiceRepo) List(ctx context.Context) ([]model.Service, error) {
	return f.list, f.listErr
}

type fakeFavoriteRepo struct {
	addRes    *model.FavoriteReq
	addErr    error
	deleteRes bool
	deleteErr error
	list      []model.FavoriteInfoReq
	listErr   error
}

func (f *fakeFavoriteRepo) Add(ctx context.Context, customerID, serviceID uint) (*model.FavoriteReq, error) {
	if f.addErr != nil {
		return nil, f.addErr
	}
	return &model.FavoriteReq{ID: 1, CustomerID: customerID, ServiceID: serviceID}, nil
}

func (f *fakeFavoriteRepo) Delete(ctx context.Context, customerID, serviceID uint) (bool, error) {
	return f.deleteRes, f.deleteErr
}

func (f *fakeFavoriteRepo) List(ctx context.Context, customerID uint) ([]model.FavoriteInfoReq, error) {
	return f.list, f.listErr
}

type fakeAuthRepo struct {
	registerID    uint
	registerErr   error
	user          *model.User
	getUserErr    error
	verifyErr     error
	deleteCodeErr error
	saveCodeErr   error

	savedEmail string
	savedCode  string
}

func (f *fakeAuthRepo) RegisterUser(ctx context.Context, email, password, role, name string) (uint, error) {
	return f.registerID, f.registerErr
}

func (f *fakeAuthRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return f.user, f.getUserErr
}

func (f *fakeAuthRepo) VerifyCode(ctx context.Context, email, code string) error {
	return f.verifyErr
}

func (f *fakeAuthRepo) DeleteUsedCode(ctx context.Context, email string) error {
	return f.deleteCodeErr
}

func (f *fakeAuthRepo) SaveVerificationCode(ctx context.Context, email, code string) error {
	f.savedEmail = email
	f.savedCode = code
	return f.saveCodeErr
}

type fakeNotifier struct {
	verifyCh chan string
	loginCh  chan string
	err      error
}

func newFakeNotifier() *fakeNotifier {
	return &fakeNotifier{
		verifyCh: make(chan string, 1),
		loginCh:  make(chan string, 1),
	}
}

func (f *fakeNotifier) SendVerificationCode(ctx context.Context, to, code string) error {
	if f.err != nil {
		return f.err
	}
	select {
	case f.verifyCh <- code:
	default:
	}
	return nil
}

func (f *fakeNotifier) SendLoginNotification(ctx context.Context, to string) error {
	if f.err != nil {
		return f.err
	}
	select {
	case f.loginCh <- to:
	default:
	}
	return nil
}
