package controllers

import (
	"aegis/assessment-test/config"
	"aegis/assessment-test/core/constant"
	"aegis/assessment-test/core/entity"
	"aegis/assessment-test/core/repository"
	"aegis/assessment-test/utils/converter"
	"aegis/assessment-test/utils/encrypt"
	"aegis/assessment-test/utils/middleware"
	"context"
	"net/http"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	cfg        *config.Config
	userRepo   repository.UserRepository
	jwtHandler *middleware.Jwt
}

func NewAuthController(
	cfg *config.Config,
	userRepo repository.UserRepository,
	jwtHandler *middleware.Jwt,
) *AuthController {
	return &AuthController{
		cfg:        cfg,
		userRepo:   userRepo,
		jwtHandler: jwtHandler,
	}
}

func (a *AuthController) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		loginReq := entity.LoginRequest{}
		c.Bind(&loginReq)
		err := c.Validate(&loginReq)
		if err != nil {
			logrus.Errorf("err validate request=%s", err.Error())
			return c.JSON(
				http.StatusBadRequest,
				constant.BadRequest(constant.CodeErrBadRequest, "there is some problem from input", err))
		}

		user, err := a.userRepo.GetUserByUsername(context.Background(), loginReq.Username)
		if err != nil || user == nil {
			logrus.Errorf("err get user by username err=%s username=%s", err.Error(), loginReq.Username)
			return c.JSON(
				http.StatusBadRequest,
				constant.InternalServerError(constant.CodeErrBadRequest, "failed to get user", err))
		}

		if !encrypt.CheckPasswordHash(loginReq.Password, user.Password) {
			logrus.Errorf("invalid password for username=%s", loginReq.Username)
			return c.JSON(
				http.StatusInternalServerError,
				constant.InternalServerError(constant.CodeErrBadRequest, "invalid password", err))
		}

		claims := middleware.Claims{
			UserId: converter.IntToString(user.ID),
			Role:   user.Role,
		}

		token, err := a.jwtHandler.Generate(context.Background(), &claims)
		if err != nil {
			logrus.Errorf("err generate token=%s", err.Error())
			return c.JSON(
				http.StatusBadRequest,
				constant.BadRequest(constant.InternalServerError, "failed generate token login", err))
		}

		resp := &entity.LoginResponse{
			Token: token,
		}

		return c.JSON(
			http.StatusOK,
			constant.Success(constant.CodeSuccess, "success login", resp))
	}
}
