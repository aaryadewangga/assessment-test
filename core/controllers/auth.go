package controllers

import (
	"aegis/assessment-test/config"
	"aegis/assessment-test/core/constant"
	"aegis/assessment-test/core/entity"
	"aegis/assessment-test/core/repository"
	"aegis/assessment-test/utils/encrypt"
	"aegis/assessment-test/utils/middleware"
	"net/http"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	cfg      *config.Config
	userRepo repository.UserRepository
}

func NewAuthController(
	cfg *config.Config,
	userRepo repository.UserRepository,
) *AuthController {
	return &AuthController{
		cfg:      cfg,
		userRepo: userRepo,
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

		user, err := a.userRepo.GetUserByUsername(c.Request().Context(), loginReq.Username)
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

		token, err := middleware.GenerateToken(user.ID, user.Role, a.cfg.JWTPublicKey)
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
