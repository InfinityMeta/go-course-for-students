package adrepo

import (
	"context"
	"fmt"

	"homework6/internal/ads"
	"homework6/internal/app"
)

type RepositoryMap struct {
	storage map[int64]*ads.Ad
}

func New() app.Repository {
	storage := make(map[int64]*ads.Ad)
	return &RepositoryMap{storage: storage}
}

func (rs *RepositoryMap) GetAdByID(ctx context.Context, id int64) (*ads.Ad, error) {

	var err error

	ad, ok := rs.storage[id]

	if !ok {
		err = fmt.Errorf("ad with id %d not found", id)
		return &ads.Ad{}, err
	}

	return ad, err

}

func (rs *RepositoryMap) StoreAd(ctx context.Context, ad *ads.Ad) error {

	var err error

	_, ok := rs.storage[ad.ID]

	if ok {
		err = fmt.Errorf("ad with id %d already exists", ad.ID)
		return err
	}

	rs.storage[ad.ID] = ad

	return nil

}
