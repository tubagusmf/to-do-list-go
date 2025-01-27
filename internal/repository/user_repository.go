package repository

import (
	"context"
	"errors"
	"log"
	"time"
	"to-do-list/internal/model"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) model.IUserRepository {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) FindByEmail(ctx context.Context, email string) *model.User {
	var user model.User

	err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("No user found with email: %s", email)
			return nil
		}
		log.Printf("Error finding user by email: %s, err: %v", email, err)
		return nil
	}

	return &user
}

func (u *UserRepo) Create(ctx context.Context, user model.User) (newUser *model.User, err error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	err = u.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) FindById(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := u.db.WithContext(ctx).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindAll(ctx context.Context, user model.User) ([]*model.User, error) {
	var users []*model.User
	query := r.db.WithContext(ctx).Model(&model.User{}).Where("deleted_at IS NULL")

	if user.Username != "" {
		query = query.Where("username LIKE ?", "%"+user.Username+"%")
	}
	if user.Email != "" {
		query = query.Where("email LIKE ?", "%"+user.Email+"%")
	}

	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserRepo) Update(ctx context.Context, user model.User) error {
	user.UpdatedAt = time.Now()

	err := u.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ? AND deleted_at IS NULL", user.Id).
		Updates(user).Error

	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) Delete(ctx context.Context, id int64) error {
	err := u.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
	if err != nil {
		return err
	}
	return nil
}
