package main

import (
	"aegis/assessment-test/config"
	"aegis/assessment-test/core/controllers"
	"aegis/assessment-test/core/repository"
	"aegis/assessment-test/routes"

	"fmt"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.Get()
	db := config.NewConnPG()

	userRepo := repository.NewUserRepository(db)

	userCon := controllers.NewUserController(userRepo)
	authCon := controllers.NewAuthController(cfg, userRepo)

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	routes.RegisterPath(
		e,
		userCon,
		authCon,
	)

	logrus.Fatal(e.Start(fmt.Sprintf(":%s", cfg.AppPort)))
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
