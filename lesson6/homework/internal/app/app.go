package app

import (
	"context"
	"errors"

	"homework6/internal/ads"
)

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
}

type AdApp struct {
	repository Repository
	adCount    int64
}

func NewApp(repo Repository) App {
	return &AdApp{repository: repo}
}

func (a *AdApp) CreateAd(ctx context.Context, title string, text string, authorId int64) (*ads.Ad, error) {

	ad := &ads.Ad{ID: a.adCount, Title: title, Text: text, AuthorID: authorId, Published: false}
	a.adCount++

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
		err = errors.New("status forbidden")
		return &ads.Ad{}, err
	}

	ad.Published = published

	return ad, nil

}

func (a *AdApp) UpdateAd(ctx context.Context, adId int64, userId int64, title string, text string) (*ads.Ad, error) {

	ad, err := a.repository.GetAdByID(ctx, adId)

	if err != nil {
		return &ads.Ad{}, err
	}

	if ad.AuthorID != userId {
		err = errors.New("status forbidden")
		return &ads.Ad{}, err
	}

	ad.Title = title
	ad.Text = text

	return ad, nil

}
