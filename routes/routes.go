package routes

import (
	"aegis/assessment-test/config"
	"aegis/assessment-test/core/controllers"
	midd "aegis/assessment-test/utils/middleware"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func RegisterPath(
	cfg *config.Config,
	e *echo.Echo,
	userCon *controllers.UserController,
	authCon *controllers.AuthController,
	productCon *controllers.ProductController,
	jwtHandler *midd.Jwt,
) {

	// Initial JWT
	jwtMiddleware := midd.NewJwt(jwtHandler)

	// Cors
	e.Use(middleware.CORS())

	// Logger
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}",
	}))

	// ROUTE REGISTER & LOGIN
	e.POST("/register", userCon.Register())
	e.POST("/login", authCon.Login())

	// Admin Routes (Protected)
	admin := e.Group("/admin", jwtMiddleware)
	admin.POST("/products", productCon.AddNewProduct())
	admin.PUT("/products", productCon.UpdateProductById())
	admin.DELETE("/products", productCon.DeleteProductById())

	// Cashier & Admin Routes
	e.GET("/products", productCon.GetProduct(), jwtMiddleware)
}
