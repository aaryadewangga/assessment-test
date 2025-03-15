package controllers

import (
	"aegis/assessment-test/core/constant"
	"aegis/assessment-test/core/entity"
	"aegis/assessment-test/core/repository"
	"aegis/assessment-test/core/repository/models"
	hash "aegis/assessment-test/utils/encrypt"
	"aegis/assessment-test/utils/middleware"
	"context"
	"net/http"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	userRepo repository.UserRepository
}

func NewUserController(
	userRepo repository.UserRepository,
) *UserController {
	return &UserController{
		userRepo: userRepo,
	}
}

func (u *UserController) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		registerReq := entity.RegisterRequest{}
		c.Bind(&registerReq)
		err := c.Validate(&registerReq)
		if err != nil {
			logrus.Errorf("err validate request=%s", err.Error())
			return c.JSON(
				http.StatusBadRequest,
				constant.BadRequest(constant.CodeErrBadRequest, "there is some problem from input", err))
		}

		isValidUsername, err := u.isExistingUsername(context.Background(), registerReq.Username)
		if err != nil && isValidUsername {
			return c.JSON(
				http.StatusBadRequest,
				constant.BadRequest(constant.CodeErrBadRequest, err.Error(), err))
		}

		hashPassword, err := hash.HashPassword(registerReq.Password)
		if err != nil {
			logrus.Errorf("err hash password=%s", err.Error())
			return c.JSON(
				http.StatusInternalServerError,
				constant.InternalServerError(constant.CodeErrInternalServer, "failed to hash password", err))
		}

		err = u.userRepo.InsertNewUser(context.Background(), &models.UserSchema{
			Name:     registerReq.Name,
			Username: registerReq.Username,
			Password: hashPassword,
			Role:     registerReq.Role,
		})
		if err != nil {
			logrus.Errorf("err insert user=%s", err.Error())
			return c.JSON(
				http.StatusInternalServerError,
				constant.InternalServerError(constant.CodeErrInternalServer, "failed to register new user", err))
		}

		resp := entity.RegisterResponse{
			Name:     registerReq.Name,
			Username: registerReq.Username,
			Role:     registerReq.Role,
		}

		return c.JSON(
			http.StatusOK,
			constant.Success(constant.CodeSuccess, "success register new user", resp))
	}
}

func (u *UserController) GetAllUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		if !middleware.IsAdmin(c) {
			logrus.Errorf("access denied")
			return c.JSON(
				http.StatusForbidden,
				constant.UnauthorizeError(constant.CodeErrForbidden, "access denied", nil))
		}

		users, err := u.userRepo.GetAllUsers(context.Background())
		if err != nil {
			logrus.Errorf("err get all users=%s", err.Error())
			return c.JSON(
				http.StatusInternalServerError,
				constant.InternalServerError(constant.CodeErrInternalServer, "failed to get all users", err))
		}

		usersList := []entity.RegisterResponse{}
		for _, v := range *users {
			tmp := entity.RegisterResponse{
				Name:     v.Name,
				Username: v.Username,
				Role:     v.Username,
				Id:       &v.ID,
			}
			usersList = append(usersList, tmp)
		}

		return c.JSON(
			http.StatusOK,
			constant.Success(constant.CodeSuccess, "success get all users", entity.GetAllUsersResponse{Users: usersList}))
	}
}

func (u *UserController) DeleteUserById() echo.HandlerFunc {
	return func(c echo.Context) error {
		if !middleware.IsAdmin(c) {
			logrus.Errorf("access denied")
			return c.JSON(
				http.StatusForbidden,
				constant.UnauthorizeError(constant.CodeErrForbidden, "access denied", nil))
		}

		id := c.QueryParam("id")
		err := u.userRepo.DeleteUser(context.Background(), id)
		if err != nil {
			logrus.Errorf("err delete user by id=%s err=%s", id, err.Error())
			return c.JSON(
				http.StatusInternalServerError,
				constant.InternalServerError(constant.CodeErrInternalServer, "err delete user by id", id))
		}

		return c.JSON(
			http.StatusOK,
			constant.Success(constant.CodeSuccess, "success delete user by id", entity.DeleteUserResponse{Id: id}))
	}
}
