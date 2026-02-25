package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   int    `json:"userId"`
	UserName string `json:"userName"`
	NickName string `json:"nickName"`
	Phone    string `json:"phone"`
	jwt.RegisteredClaims
}

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未授权"})
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(authHeader)
		parts := strings.SplitN(tokenString, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			tokenString = strings.TrimSpace(parts[1])
		}
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未授权"})
			c.Abort()
			return
		}
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未授权"})
			c.Abort()
			return
		}

		c.Set("userId", claims.UserID)
		c.Set("userName", claims.UserName)
		c.Set("nickName", claims.NickName)
		c.Set("phone", claims.Phone)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"},
	})
}
