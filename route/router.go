package route

import (
	"dbo-api/config"
	"dbo-api/handler"
	"dbo-api/repository"
	"dbo-api/service"

	"github.com/gin-gonic/gin"
)

func AuthRouter(api *gin.RouterGroup) {
	authRepository := repository.NewAuthRepository(config.DB)
	custRepository := repository.NewCustomerRepository(config.DB)
	authService := service.NewAuthService(authRepository, custRepository)
	authHandler := handler.NewAuthHandler(authService)

	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login)
}
