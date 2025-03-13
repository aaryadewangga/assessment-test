package middleware

import (
	"aegis/assessment-test/core/constant"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type contextKey string

const claimsContextKey contextKey = "jwtClaims"

type Jwt struct {
	Issuer        string
	Secret        JwtSecret
	Expiration    time.Duration
	SigningMethod jwt.SigningMethod
}

type JwtOpt func(*Jwt)

type Claims struct {
	UserId string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func WithSigningMethod(method jwt.SigningMethod) JwtOpt {
	return func(j *Jwt) {
		j.SigningMethod = method
	}
}

func WithSecret(secret JwtSecret) JwtOpt {
	return func(j *Jwt) {
		j.Secret = secret
	}
}

func NewJwt(jwtInstance *Jwt) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, constant.UnauthorizeError(constant.CodeErrInternalServer, "unauthorize", nil))
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, constant.UnauthorizeError(constant.CodeErrInternalServer, "unauthorize", nil))
			}

			claims := &Claims{}
			ctx, err := jwtInstance.Verify(c.Request().Context(), tokenParts[1], claims)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, constant.UnauthorizeError(constant.CodeErrInternalServer, "invalid token", nil))
			}

			c.Set("jwtClaims", claims)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func (j *Jwt) Generate(ctx context.Context, claims *Claims) (string, error) {
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(j.Expiration))
	claims.Issuer = j.Issuer
	token := jwt.NewWithClaims(j.SigningMethod, claims)
	signedString := j.Secret.GetSign()
	return token.SignedString(signedString)
}

func (j *Jwt) Verify(ctx context.Context, tokenString string, claims *Claims) (context.Context, error) {
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != j.SigningMethod.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.Secret.GetVerify(), nil
	})

	if err != nil {
		return ctx, fmt.Errorf("failed to parse token: %w", err)
	}
	if !parsedToken.Valid {
		return ctx, errors.New("invalid token")
	}

	ctx = context.WithValue(ctx, claimsContextKey, claims)
	return ctx, nil
}

func GetTokenClaims(c echo.Context) (*Claims, error) {
	rawClaims := c.Get("jwtClaims")
	claims, ok := rawClaims.(*Claims)
	if !ok {
		return nil, errors.New("claims wrong type")
	}
	return claims, nil
}

func IsAdmin(c echo.Context) bool {
	claims, err := GetTokenClaims(c)
	if err != nil {
		logrus.Errorf("failed claim token by context=%s", err.Error())
		return false
	}

	if claims.Role != "admin" {
		logrus.Errorf("access denied userId=%s", claims.UserId)
		return false
	}

	return true
}
