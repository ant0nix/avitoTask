package handler

import (
	"github.com/ant0nix/avitoTask/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	added := router.Group("/start")
	{
		added.POST("/user", h.CreateUser)
		added.POST("/services", h.CreateServices)
	}

	services := router.Group("/services")
	{
		services.POST("/change-balance", h.ChangeBalance)
		services.GET("/show-balance", h.ShowBalance)
		services.PUT("/p2p", h.P2p)
		services.GET("/", h.ListServices)
		services.POST("/new-order", h.MakeOrder)
		services.PATCH("/do-order/:id", h.DoOrder)
	}
	return router
}
