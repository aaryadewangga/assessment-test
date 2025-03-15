package repository

import (
	"aegis/assessment-test/core/repository/models"
	"context"

	"github.com/go-pg/pg/v10"
)

type UserRepository interface {
	InsertNewUser(ctx context.Context, user *models.UserSchema) error
	GetUserByUsername(ctx context.Context, username string) (*models.UserSchema, error)
	GetUserById(ctx context.Context, id string) (*models.UserSchema, error)
	GetAllUsers(ctx context.Context) (*[]models.UserSchema, error)
	DeleteUser(ctx context.Context, id string) error
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

func (u *userRepository) GetUserById(ctx context.Context, id string) (*models.UserSchema, error) {
	user := new(models.UserSchema)

	if err := u.db.Model(user).
		Context(ctx).
		Where("? = ?", pg.Ident("ID"), id).
		Limit(1).
		Select(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userRepository) GetAllUsers(ctx context.Context) (*[]models.UserSchema, error) {
	user := new([]models.UserSchema)

	if err := u.db.Model(user).
		Context(ctx).
		Select(); err != nil {
		return nil, err
	}

	return user, nil
}

func (p *userRepository) DeleteUser(ctx context.Context, id string) error {
	_, err := p.db.Model((*models.UserSchema)(nil)).Context(ctx).
		Where("? = ?", pg.Ident("ID"), id).
		Delete()
	return err
}
