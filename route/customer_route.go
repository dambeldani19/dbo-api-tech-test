package route

import (
	"dbo-api/config"
	"dbo-api/handler"
	"dbo-api/repository"
	"dbo-api/service"

	"github.com/gin-gonic/gin"
)

func CustomerRoute(api *gin.RouterGroup) {
	custRepo := repository.NewCustomerRepository(config.DB)
	authRepo := repository.NewAuthRepository(config.DB)
	custService := service.NewCustomerService(custRepo, authRepo)
	handler := handler.NewCustomerHandler(custService)

	r := api.Group("/customer")

	//customer api
	r.GET("/", handler.GetList)
	r.GET("/:id", handler.GetDetail)
	r.PUT("/:id", handler.Update)
	r.DELETE("/:id", handler.Delete)
}
