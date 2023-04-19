package adrepo

import (
	"context"
	"fmt"

	"homework6/internal/ads"
	"homework6/internal/app"
)

type RepositoryApp struct {
	storage map[int64]*ads.Ad
}

func New() app.Repository {
	storage := make(map[int64]*ads.Ad)
	return &RepositoryApp{storage: storage}
}

func (rs *RepositoryApp) GetAdByID(ctx context.Context, id int64) (*ads.Ad, error) {

	var err error

	ad, ok := rs.storage[id]

	if !ok {
		err = fmt.Errorf("ad with id %d not found", id)
		return &ads.Ad{}, err
	}

	return ad, err

}

func (rs *RepositoryApp) StoreAd(ctx context.Context, ad *ads.Ad) error {

	var err error

	_, ok := rs.storage[ad.ID]

	if ok {
		err = fmt.Errorf("ad with id %d already exists", ad.ID)
		return err
	}

	rs.storage[ad.ID] = ad

	return nil

}

func (rs *RepositoryApp) Len(ctx context.Context) int64 {
	return int64(len(rs.storage))
}

func (rs *RepositoryApp) UpdateADByID(ctx context.Context, adID int64, title string, text string) {

	rs.storage[adID].Title = title
	rs.storage[adID].Text = text

}
