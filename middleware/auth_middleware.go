package middleware

import (
	"bunaken-boat-backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.APIError(c, http.StatusUnauthorized, "Token tidak ditemukan")
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, " ")
		if len(tokenString) != 2 {
			utils.APIError(c, http.StatusUnauthorized, "Format token salah")
			c.Abort()
			return
		}

		token, err := utils.ValidateToken(tokenString[1])
		if err != nil || !token.Valid {
			utils.APIError(c, http.StatusUnauthorized, "Token tidak valid")
			c.Abort()
			return
		}
		
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			userId := uint(claims["user_id"].(float64))
			c.Set("user_id", userId)
		}

		c.Next()
	}
}
