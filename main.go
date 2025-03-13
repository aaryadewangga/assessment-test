package main

import (
	"aegis/assessment-test/config"
	"aegis/assessment-test/core/controllers"
	"aegis/assessment-test/core/repository"
	"aegis/assessment-test/routes"
	"aegis/assessment-test/utils/middleware"
	"time"

	"fmt"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.Get()
	db := config.NewConnPG()

	edDSASecret, err := middleware.NewEdDSASecret(cfg.JWTPrivateKey)
	if err != nil {
		logrus.Errorf("failed to create EdDSA secret: %v", err)
	}

	jwtDeps := &middleware.Jwt{
		Issuer:        "my-local",
		Secret:        edDSASecret,
		Expiration:    time.Minute * time.Duration(15),
		SigningMethod: jwt.SigningMethodEdDSA,
	}

	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)

	userCon := controllers.NewUserController(userRepo)
	authCon := controllers.NewAuthController(cfg, userRepo, jwtDeps)
	productCon := controllers.NewProductController(productRepo)

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	routes.RegisterPath(
		cfg,
		e,
		userCon,
		authCon,
		productCon,
		jwtDeps,
	)

	logrus.Fatal(e.Start(fmt.Sprintf(":%s", cfg.AppPort)))
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
