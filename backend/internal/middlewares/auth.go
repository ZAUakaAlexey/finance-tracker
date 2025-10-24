package middlewares

import (
	"net/http"
	"strings"

	"github.com/ZAUakaAlexey/backend_go/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func Authenticate(context *gin.Context) {
	authHeader := context.Request.Header.Get("Authorization")

	if authHeader == "" {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		context.Abort()
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid authorization header format. Expected: Bearer <token>",
		})
		context.Abort()
		return
	}

	tokenString := parts[1]

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		cfg, _ := config.LoadConfig()
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid or expired token",
		})
		context.Abort()
		return
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		context.Abort()
		return
	}

	context.Set("user_id", claims.UserID)
	context.Next()
}
