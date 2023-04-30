package app

import (
	"context"
	"errors"

	"homework6/internal/ads"
)

var ErrStatusForbidden = errors.New("status forbidden")

type App interface {
	// TODO: реализовать
	CreateAd(context.Context, string, string, int64) (*ads.Ad, error)
	ChangeAdStatus(context.Context, int64, int64, bool) (*ads.Ad, error)
	UpdateAd(context.Context, int64, int64, string, string) (*ads.Ad, error)
}

type Repository interface {
	// TODO: реализовать
	StoreAd(context.Context, *ads.Ad) error
	GetAdByID(context.Context, int64) (*ads.Ad, error)
	UpdateADByID(context.Context, int64, string, string)
	Len(context.Context) int64
}

type AdApp struct {
	repository Repository
}

func NewApp(repo Repository) App {
	return &AdApp{repository: repo}
}

func (a *AdApp) CreateAd(ctx context.Context, title string, text string, authorId int64) (*ads.Ad, error) {

	ad := &ads.Ad{ID: a.repository.Len(ctx), Title: title, Text: text, AuthorID: authorId, Published: false}

	err := a.repository.StoreAd(ctx, ad)

	if err != nil {
		return &ads.Ad{}, err
	}

	return ad, nil

}

func (a *AdApp) ChangeAdStatus(ctx context.Context, adId int64, authorId int64, published bool) (*ads.Ad, error) {

	ad, err := a.repository.GetAdByID(ctx, adId)

	if err != nil {
		return &ads.Ad{}, err
	}

	if ad.AuthorID != authorId {
		err = ErrStatusForbidden
		return &ads.Ad{}, err
	}

	ad.Published = published

	return ad, nil

}

func (a *AdApp) UpdateAd(ctx context.Context, adID int64, userID int64, title string, text string) (*ads.Ad, error) {

	ad, err := a.repository.GetAdByID(ctx, adID)

	if err != nil {
		return &ads.Ad{}, err
	}

	if ad.AuthorID != userID {
		err = ErrStatusForbidden
		return &ads.Ad{}, err
	}

	a.repository.UpdateADByID(ctx, adID, title, text)

	return ad, nil

}
