package httpgin

import (
	"time"

	"github.com/gin-gonic/gin"

	"homework8/internal/ads"
	"homework8/internal/users"
)

type createAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type createUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type adResponse struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Text         string    `json:"text"`
	AuthorID     int64     `json:"author_id"`
	Published    bool      `json:"published"`
	CreationDate time.Time `json:"creation_date"`
	UpdateDate   time.Time `json:"update_date"`
}

type userResponse struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type changeAdStatusRequest struct {
	Published bool  `json:"published"`
	UserID    int64 `json:"user_id"`
}

type updateAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type updateUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

func AdSuccessResponse(ad *ads.Ad) *gin.H {
	return &gin.H{
		"data": adResponse{
			ID:           ad.ID,
			Title:        ad.Title,
			Text:         ad.Text,
			AuthorID:     ad.AuthorID,
			Published:    ad.Published,
			CreationDate: ad.CreationDate,
			UpdateDate:   ad.UpdateDate,
		},
		"error": nil,
	}
}

func AdErrorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}

func AdsSuccessResponse(ads []*ads.Ad) *gin.H {

	resps := []adResponse{}

	for _, ad := range ads {
		adResp := adResponse{
			ID:           ad.ID,
			Title:        ad.Title,
			Text:         ad.Text,
			AuthorID:     ad.AuthorID,
			Published:    ad.Published,
			CreationDate: ad.CreationDate,
			UpdateDate:   ad.UpdateDate,
		}
		resps = append(resps, adResp)
	}

	return &gin.H{
		"data":  resps,
		"error": nil,
	}

}

func AdsErrorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}

func UserSuccessResponse(user *users.User) *gin.H {
	return &gin.H{
		"data": userResponse{
			ID:       user.ID,
			Nickname: user.Nickname,
			Email:    user.Email,
		},
		"error": nil,
	}
}

func UserErrorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}
