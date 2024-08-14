package route

import (
	"dbo-api/config"
	"dbo-api/handler"
	"dbo-api/middleware"
	"dbo-api/repository"
	"dbo-api/service"

	"github.com/gin-gonic/gin"
)

func OrderRoute(api *gin.RouterGroup) {
	repo := repository.NewOrderRepository(config.DB)
	prodRepo := repository.NewProductRepository(config.DB)
	srv := service.NewOrderService(repo, prodRepo)
	handler := handler.NewOrderHandler(srv)

	r := api.Group("/order")
	r.Use(middleware.JWTMiddleware(config.DB))

	//order api
	r.POST("/", handler.CreateOrder)
	r.GET("/", handler.GetListOrder)
	r.GET("/:id", handler.GetDetailOrder)
	r.PUT("/:id", handler.UpdateOrder)
	r.DELETE("/:id/product", handler.DeleteOrderProduct)
}
