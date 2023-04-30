package httpgin

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"homework8/internal/app"
)

// Метод для создания объявления (ad)
func createAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		err := c.ShouldBindJSON(&reqBody)

		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		//TODO: вызов логики, например, CreateAd(c.Context(), reqBody.Title, reqBody.Text, reqBody.UserID)
		// TODO: метод должен возвращать AdSuccessResponse или ошибку.
		ad, err := a.CreateAd(c, reqBody.Title, reqBody.Text, reqBody.UserID)

		if err != nil {
			if errors.Is(err, app.ErrNotFound) {
				c.JSON(http.StatusNotFound, AdErrorResponse(err))
				return
			}
			if errors.Is(err, app.ErrNotValid) {
				c.JSON(http.StatusBadRequest, AdErrorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
func changeAdStatus(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody changeAdStatusRequest
		err := c.ShouldBindJSON(&reqBody)

		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		adID, err := strconv.ParseInt(c.Param("ad_id"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		// TODO: вызов логики ChangeAdStatus(c.Context(), int64(adID), reqBody.UserID, reqBody.Published)
		// TODO: метод должен возвращать AdSuccessResponse или ошибку.

		ad, err := a.ChangeAdStatus(c, adID, reqBody.UserID, reqBody.Published)

		if err != nil {
			if errors.Is(err, app.ErrNotFound) {
				c.JSON(http.StatusNotFound, AdErrorResponse(err))
				return
			}

			if errors.Is(err, app.ErrStatusForbidden) {
				c.JSON(http.StatusForbidden, AdErrorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateAdRequest
		err := c.ShouldBindJSON(&reqBody)

		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		adID, err := strconv.ParseInt(c.Param("ad_id"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		// TODO: вызов логики, например, UpdateAd(c.Context(), int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)
		// TODO: метод должен возвращать AdSuccessResponse или ошибку.
		ad, err := a.UpdateAd(c, adID, reqBody.UserID, reqBody.Title, reqBody.Text)

		if err != nil {
			if errors.Is(err, app.ErrNotFound) {
				c.JSON(http.StatusNotFound, AdErrorResponse(err))
				return
			}

			if errors.Is(err, app.ErrStatusForbidden) {
				c.JSON(http.StatusForbidden, AdErrorResponse(err))
				return
			}

			if errors.Is(err, app.ErrNotValid) {
				c.JSON(http.StatusBadRequest, AdErrorResponse(err))
				return
			}

			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для получения объявления по id

func getAdByID(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		adID, err := strconv.ParseInt(c.Param("ad_id"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, err := a.GetAdByID(c, adID)

		if err != nil {
			if errors.Is(err, app.ErrNotFound) {
				c.JSON(http.StatusNotFound, AdErrorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))

	}
}

// Метод для получения опубликованных объявлений

func filterAds(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {

		authorID, err := strconv.ParseInt(c.Query("author_id"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, AdsErrorResponse(err))
			return
		}

		pubBefore, err := time.Parse(time.RFC3339Nano, c.Query("pub_before"))

		if err != nil {
			c.JSON(http.StatusBadRequest, AdsErrorResponse(err))
			return
		}

		pubAfter, err := time.Parse(time.RFC3339Nano, c.Query("pub_after"))

		if err != nil {
			c.JSON(http.StatusBadRequest, AdsErrorResponse(err))
			return
		}

		adsList, err := a.FilterAds(c, app.WithAuthorID(authorID), app.WithPublishedBefore(pubBefore.UTC()), app.WithPublishedAfter(pubAfter.UTC()))

		if err != nil {
			if errors.Is(err, app.ErrNotFound) {
				c.JSON(http.StatusNotFound, AdErrorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdsSuccessResponse(adsList))

	}
}

// Метод для создания пользователя (user)
func createUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createUserRequest
		err := c.ShouldBindJSON(&reqBody)

		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}

		user := a.CreateUser(c, reqBody.Nickname, reqBody.Email)

		c.JSON(http.StatusOK, UserSuccessResponse(user))
	}
}

// Метод для обновления никнейма(Nickname) или емейла(Email) пользователя
func updateUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateUserRequest
		err := c.ShouldBindJSON(&reqBody)

		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}

		userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}

		user, err := a.UpdateUser(c, userID, reqBody.Nickname, reqBody.Email)

		if err != nil {
			if errors.Is(err, app.ErrNotFound) {
				c.JSON(http.StatusNotFound, AdErrorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, UserErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(user))
	}
}

// Метод для получения объявления по имени

func searchAdByName(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {

		title := c.Param("title")

		ad, err := a.SearchAdByName(c, title)

		if err != nil {
			if errors.Is(err, app.ErrNotFound) {
				c.JSON(http.StatusNotFound, AdErrorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))

	}
}
