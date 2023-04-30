package httpgin

import (
	"github.com/gin-gonic/gin"

	"homework8/internal/app"
)

func AppRouter(r gin.IRouter, a app.App) {
	r.Use(Logger())                                //Логгер
	r.Use(gin.Recovery())                          //panic recovery
	r.POST("/ads", createAd(a))                    // Метод для создания объявления (ad)
	r.PUT("/ads/:ad_id/status", changeAdStatus(a)) // Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
	r.PUT("/ads/:ad_id", updateAd(a))              // Метод для обновления текста(Text) или заголовка(Title) объявления
	r.GET("/ads/:ad_id", getAdByID(a))             // Метод для получения объявления по id
	r.GET("/ads", filterAds(a))                    // Метод для получения опубликованных объявлений
	r.POST("/users", createUser(a))                // Метод для создания пользователя (user)
	r.PUT("/users/:user_id", updateUser(a))        // Метод для обновления никнейма(Nickname) или емейла(Email) пользователя
	r.GET("/ads/search/:title", searchAdByName(a)) // Метод для поиска объявления по названию
}
