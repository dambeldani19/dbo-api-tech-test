package route

import (
	"dbo-api/config"
	"dbo-api/handler"
	"dbo-api/middleware"
	"dbo-api/repository"
	"dbo-api/service"

	"github.com/gin-gonic/gin"
)

func ProductRoute(api *gin.RouterGroup) {
	productRepo := repository.NewProductRepository(config.DB)
	productSrv := service.NewProductService(productRepo)
	handler := handler.NewProductHandler(productSrv)

	r := api.Group("/product")
	r.Use(middleware.JWTMiddleware(config.DB))

	//product api
	r.POST("/", handler.AddProduct)
	r.GET("/", handler.GetListProduct)
	r.GET("/:id", handler.GetDetail)
	r.PUT("/:id", handler.Update)
	r.DELETE("/:id", handler.Delete)
}
