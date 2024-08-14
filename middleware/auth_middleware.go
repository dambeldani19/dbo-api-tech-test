package middleware

import (
	"dbo-api/entity"
	"dbo-api/errorhandler"
	"dbo-api/helper"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func JWTMiddleware(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			errorhandler.HandlerError(c, &errorhandler.UnathorizedError{Message: "Unauthorize"})
			c.Abort()
			return
		}

		userId, err := helper.ValidateToken(tokenString)
		if err != nil {
			errorhandler.HandlerError(c, &errorhandler.UnathorizedError{Message: err.Error()})
			c.Abort()
			return
		}

		var customer entity.Customer
		if err := db.Where("user_id = ?", *userId).First(&customer).Error; err != nil {
			errorhandler.HandlerError(c, &errorhandler.UnathorizedError{Message: "Customer not found"})
			c.Abort()
			return
		}

		c.Set("userID", *userId)
		c.Set("customerID", customer.ID)
		c.Next()

	}
}
