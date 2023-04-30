package adrepo

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"homework8/internal/ads"
	"homework8/internal/app"
	"homework8/internal/users"
)

var ErrNotFound = fmt.Errorf("not found")

type StorageAd struct {
	mx   *sync.RWMutex
	data map[int64]*ads.Ad
}

type StorageUser struct {
	mx   *sync.RWMutex
	data map[int64]*users.User
}

type RepositoryApp struct {
	storageAd   *StorageAd
	storageUser *StorageUser
}

func New() app.Repository {
	storageAd := &StorageAd{mx: &sync.RWMutex{}, data: make(map[int64]*ads.Ad)}
	storageUser := &StorageUser{mx: &sync.RWMutex{}, data: make(map[int64]*users.User)}
	return &RepositoryApp{storageAd: storageAd, storageUser: storageUser}
}

func (rs *RepositoryApp) GetAdByID(ctx context.Context, adID int64) (*ads.Ad, error) {
	rs.storageAd.mx.RLock()
	defer rs.storageAd.mx.RUnlock()

	ad, ok := rs.storageAd.data[adID]

	if !ok {
		return &ads.Ad{}, ErrNotFound
	}

	return ad, nil

}

func (rs *RepositoryApp) StoreAd(ctx context.Context, ad *ads.Ad) {
	rs.storageAd.mx.Lock()
	defer rs.storageAd.mx.Unlock()

	rs.storageAd.data[ad.ID] = ad

}

func (rs *RepositoryApp) LenAd(ctx context.Context) int64 {
	rs.storageAd.mx.RLock()
	defer rs.storageAd.mx.RUnlock()

	return int64(len(rs.storageAd.data))
}

func (rs *RepositoryApp) UpdateADByID(ctx context.Context, adID int64, title string, text string) {
	rs.storageAd.mx.Lock()
	defer rs.storageAd.mx.Unlock()

	rs.storageAd.data[adID].Title = title
	rs.storageAd.data[adID].Text = text
	rs.storageAd.data[adID].UpdateDate = time.Now().UTC()

}

func (rs *RepositoryApp) UpdateAdStatus(ctx context.Context, adID int64, status bool) {
	rs.storageAd.mx.Lock()
	defer rs.storageAd.mx.Unlock()

	rs.storageAd.data[adID].Published = status
	rs.storageAd.data[adID].UpdateDate = time.Now().UTC()

}

func (rs *RepositoryApp) LenUser(ctx context.Context) int64 {
	rs.storageUser.mx.RLock()
	defer rs.storageUser.mx.RUnlock()

	return int64(len(rs.storageUser.data))
}

func (rs *RepositoryApp) GetUserByID(ctx context.Context, userID int64) (*users.User, error) {
	rs.storageUser.mx.RLock()
	defer rs.storageUser.mx.RUnlock()

	user, ok := rs.storageUser.data[userID]

	if !ok {
		return &users.User{}, ErrNotFound
	}

	return user, nil

}

func (rs *RepositoryApp) StoreUser(ctx context.Context, user *users.User) {
	rs.storageUser.mx.Lock()
	defer rs.storageUser.mx.Unlock()

	rs.storageUser.data[user.ID] = user

}

func (rs *RepositoryApp) UpdateUserByID(ctx context.Context, userID int64, nickname string, email string) {
	rs.storageUser.mx.Lock()
	defer rs.storageUser.mx.Unlock()

	rs.storageUser.data[userID].Nickname = nickname
	rs.storageUser.data[userID].Email = email

}

func (rs *RepositoryApp) SearchAdByName(ctx context.Context, adName string) (*ads.Ad, error) {
	rs.storageAd.mx.RLock()
	defer rs.storageAd.mx.RUnlock()

	for _, v := range rs.storageAd.data {
		if strings.Contains(v.Title, adName) {
			return v, nil
		}
	}

	return &ads.Ad{}, ErrNotFound

}

func (rs *RepositoryApp) FilterAds(ctx context.Context, filter *app.Filter) ([]*ads.Ad, error) {
	rs.storageAd.mx.RLock()
	defer rs.storageAd.mx.RUnlock()

	res := []*ads.Ad{}

	for _, v := range rs.storageAd.data {
		if !v.Published {
			continue
		}
		if filter.AuthorID != -1 && v.AuthorID != filter.AuthorID {
			continue
		}
		if !filter.PublishedAfter.IsZero() && v.CreationDate.Before(filter.PublishedAfter) {
			continue
		}
		if !filter.PublishedBefore.IsZero() && v.CreationDate.After(filter.PublishedBefore) {
			continue
		}
		res = append(res, v)
	}

	if len(res) == 0 {
		return res, ErrNotFound
	}

	return res, nil
}
