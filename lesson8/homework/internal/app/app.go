package app

import (
	"context"
	"errors"
	"time"

	"github.com/InfinityMeta/validator"

	"homework8/internal/ads"
	"homework8/internal/users"
)

var (
	ErrStatusForbidden = errors.New("forbidden")
	ErrNotFound        = errors.New("not found")
	ErrNotValid        = errors.New("not valid")
)

type App interface {
	// TODO: реализовать
	CreateAd(context.Context, string, string, int64) (*ads.Ad, error)
	CreateUser(context.Context, string, string) *users.User
	ChangeAdStatus(context.Context, int64, int64, bool) (*ads.Ad, error)
	UpdateAd(context.Context, int64, int64, string, string) (*ads.Ad, error)
	UpdateUser(context.Context, int64, string, string) (*users.User, error)
	CheckUserExists(context.Context, int64) bool
	GetAdByID(context.Context, int64) (*ads.Ad, error)
	SearchAdByName(context.Context, string) (*ads.Ad, error)
	FilterAds(context.Context, ...FilterOption) ([]*ads.Ad, error)
}

type Repository interface {
	// TODO: реализовать
	StoreAd(context.Context, *ads.Ad)
	StoreUser(context.Context, *users.User)
	GetAdByID(context.Context, int64) (*ads.Ad, error)
	GetUserByID(context.Context, int64) (*users.User, error)
	UpdateAdStatus(context.Context, int64, bool)
	UpdateADByID(context.Context, int64, string, string)
	UpdateUserByID(context.Context, int64, string, string)
	LenAd(context.Context) int64
	LenUser(context.Context) int64
	SearchAdByName(context.Context, string) (*ads.Ad, error)
	FilterAds(context.Context, *Filter) ([]*ads.Ad, error)
}

type AdApp struct {
	repository Repository
}

func NewApp(repo Repository) App {
	return &AdApp{repository: repo}
}

func (a *AdApp) CreateAd(ctx context.Context, title string, text string, authorId int64) (*ads.Ad, error) {

	if !a.CheckUserExists(ctx, authorId) {
		return &ads.Ad{}, ErrNotFound
	}

	ad := &ads.Ad{ID: a.repository.LenAd(ctx), Title: title, Text: text, AuthorID: authorId, Published: false, CreationDate: time.Now().UTC(), UpdateDate: time.Time{}}

	err := validator.Validate(ad)

	if err != nil {
		return &ads.Ad{}, ErrNotValid
	}

	a.repository.StoreAd(ctx, ad)

	return ad, nil

}

func (a *AdApp) ChangeAdStatus(ctx context.Context, adID int64, authorID int64, published bool) (*ads.Ad, error) {

	if !a.CheckUserExists(ctx, authorID) {
		return &ads.Ad{}, ErrNotFound
	}

	ad, err := a.repository.GetAdByID(ctx, adID)

	if err != nil {
		return &ads.Ad{}, ErrNotFound
	}

	if ad.AuthorID != authorID {
		return &ads.Ad{}, ErrStatusForbidden
	}

	a.repository.UpdateAdStatus(ctx, adID, published)

	return ad, nil

}

func (a *AdApp) UpdateAd(ctx context.Context, adID int64, authorID int64, title string, text string) (*ads.Ad, error) {

	if !a.CheckUserExists(ctx, authorID) {
		return &ads.Ad{}, ErrNotFound
	}

	ad, err := a.repository.GetAdByID(ctx, adID)

	if err != nil {
		return &ads.Ad{}, ErrNotFound
	}

	if ad.AuthorID != authorID {
		return &ads.Ad{}, ErrStatusForbidden
	}

	a.repository.UpdateADByID(ctx, adID, title, text)

	err = validator.Validate(ad)

	if err != nil {
		return &ads.Ad{}, ErrNotValid
	}

	return ad, nil

}

func (a *AdApp) CheckUserExists(ctx context.Context, userID int64) bool {

	_, err := a.repository.GetUserByID(ctx, userID)

	return err == nil

}

func (a *AdApp) CreateUser(ctx context.Context, nickname string, email string) *users.User {

	user := &users.User{ID: a.repository.LenUser(ctx), Nickname: nickname, Email: email}

	a.repository.StoreUser(ctx, user)

	return user

}

func (a *AdApp) UpdateUser(ctx context.Context, userID int64, nickname string, email string) (*users.User, error) {

	if !a.CheckUserExists(ctx, userID) {
		return &users.User{}, ErrNotFound
	}

	user, err := a.repository.GetUserByID(ctx, userID)

	if err != nil {
		return &users.User{}, ErrNotFound
	}

	a.repository.UpdateUserByID(ctx, userID, nickname, email)

	return user, nil

}

func (a *AdApp) GetAdByID(ctx context.Context, adID int64) (*ads.Ad, error) {

	ad, err := a.repository.GetAdByID(ctx, adID)

	if err != nil {
		return &ads.Ad{}, ErrNotFound
	}

	return ad, nil

}

func (a *AdApp) SearchAdByName(ctx context.Context, adName string) (*ads.Ad, error) {

	ad, err := a.repository.SearchAdByName(ctx, adName)

	if err != nil {
		return &ads.Ad{}, ErrNotFound
	}

	return ad, nil

}

func (a *AdApp) FilterAds(ctx context.Context, options ...FilterOption) ([]*ads.Ad, error) {

	filteredAds, err := a.repository.FilterAds(ctx, NewFilter(options...))

	if err != nil {
		return []*ads.Ad{}, ErrNotFound
	}

	return filteredAds, nil

}
