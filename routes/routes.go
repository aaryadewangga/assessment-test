package routes

import (
	"aegis/assessment-test/config"
	"aegis/assessment-test/core/controllers"
	midd "aegis/assessment-test/utils/middleware"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func RegisterPath(
	e *echo.Echo,
	cfg *config.Config,
	jwtHandler *midd.Jwt,
	userCon *controllers.UserController,
	authCon *controllers.AuthController,
	productCon *controllers.ProductController,
	trxCon *controllers.TransactionController,
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

	trx := e.Group("/transactions", jwtMiddleware)
	trx.POST("/create", trxCon.CreateTransaction(), jwtMiddleware)
	trx.GET("/list", trxCon.GetAllTransactions(), jwtMiddleware)
	trx.GET("/details", trxCon.GetTransactionDetailsById(), jwtMiddleware)
	trx.GET("/generate/pdf", trxCon.GetTransactionPDF(), jwtMiddleware)
	trx.GET("/generate/excel", trxCon.GetTransactionExcel(), jwtMiddleware)
}
