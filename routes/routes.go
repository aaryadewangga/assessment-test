package routes

import (
	"aegis/assessment-test/core/controllers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func RegisterPath(e *echo.Echo,
	userCon *controllers.UserController,
	authCon *controllers.AuthController,
) {

	//CORS
	e.Use(middleware.CORS())

	//LOGGER
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}",
	}))

	//ROUTE REGISTER & LOGIN
	e.POST("/register", userCon.Register())
	e.POST("/login", authCon.Login())
}
