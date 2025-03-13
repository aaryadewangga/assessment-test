package middleware

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Middleware to check JWT Token
func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(secret),
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
	})
}

// Extract user info from token
func ExtractUser(c echo.Context) (int, string, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int(claims["user_id"].(float64))
	role := claims["role"].(string)

	return userID, role, nil
}

// Admin-Only Middleware
func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, role, err := ExtractUser(c)
		if err != nil || role != "admin" {
			return c.JSON(http.StatusForbidden, map[string]string{"message": "Access denied"})
		}
		return next(c)
	}
}

// generate token
func GenerateToken(userID int, role, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Minute * 15).Unix(), // Token expires in 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
