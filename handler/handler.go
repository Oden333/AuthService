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
		auth.POST("/sign-up", h.singUp)
		auth.POST("/sign-in", h.singIn)
		auth.POST("/auth/refresh", h.userRefresh)
	}
	//Группа маршрутов для авторизированных пользвателей (запросы с токеном)
	authenticated := router.Group("/", h.getUserIdentity)
	{
		authenticated.GET("/accountId", h.getUserAccount)
	}

	return router
}
