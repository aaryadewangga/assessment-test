package repository

import (
	"aegis/assessment-test/core/repository/models"
	"context"

	"github.com/go-pg/pg/v10"
)

type UserRepository interface {
	InsertNewUser(ctx context.Context, user *models.UserSchema) error
	GetUserByUsername(ctx context.Context, username string) (*models.UserSchema, error)
}

type userRepository struct {
	db *pg.DB
}

func NewUserRepository(db *pg.DB) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) InsertNewUser(ctx context.Context, user *models.UserSchema) error {
	_, err := u.db.Model(user).Context(ctx).Insert()
	return err
}

func (u *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.UserSchema, error) {
	user := new(models.UserSchema)

	if err := u.db.Model(user).
		Context(ctx).
		Where("? = ?", pg.Ident("USERNAME"), username).
		Limit(1).
		Select(); err != nil {
		return nil, err
	}

	return user, nil
}
