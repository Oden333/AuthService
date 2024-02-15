package handler

import (
	"Auth/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	//Конфигурация страницы ответа
	router := gin.Default()
	//Конфигурация рутов

	//Группа маршрутов авторизации
	auth := router.Group("/auth")
	{
		//Маршрут добавляет пользвателя в БД
		auth.POST("/sign-up", h.singUp)
		//Маршрут выдает пару Access, Refresh токенов для пользователя с данными email,password
		auth.POST("/sign-in", h.singIn)
		//Маршрут выдает пару Access, Refresh токенов для пользователя сидентификатором (GUID) указанным в параметре запроса
		auth.POST("/sign-in_GUID", h.singInByGiud)
		//Маршрут выполняет Refresh операцию на пару Access, Refreshтокенов
		auth.POST("/auth/refresh", h.userRefresh)
	}
	//Группа маршрутов для авторизированных пользвателей (запросы с токеном)
	authenticated := router.Group("/", h.getUserIdentity)
	{
		authenticated.GET("/accountId", h.getUserAccount)
	}

	return router
}
