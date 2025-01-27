package usecase

import (
	"context"
	"errors"
	"time"
	"to-do-list/internal/helper"
	"to-do-list/internal/model"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var v = validator.New()

type UserUsecase struct {
	userRepo model.IUserRepository
}

func NewUserUsecase(
	userRepo model.IUserRepository,
) model.IUserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (u *UserUsecase) Login(ctx context.Context, in model.LoginInput) (token string, err error) {
	log := logrus.WithFields(logrus.Fields{
		"email": in.Email,
	})

	if err := v.Struct(in); err != nil {
		log.Error("Validation error: ", err)
		return "", err
	}

	user := u.userRepo.FindByEmail(ctx, in.Email)
	if user == nil {
		err = errors.New("wrong email or password")
		return
	}

	if !helper.CheckPasswordHash(in.Password, user.Password) {
		err = errors.New("missmatch password")
		return
	}

	token, err = helper.GenerateToken(user.Id)

	if err != nil {
		log.Error(err)
	}

	return
}

func (u *UserUsecase) FindAll(ctx context.Context, user model.User) ([]*model.User, error) {
	log := logrus.WithFields(logrus.Fields{
		"filter": user,
	})

	users, err := u.userRepo.FindAll(ctx, user)
	if err != nil {
		log.Error("Failed to fetch users: ", err)
		return nil, err
	}

	return users, nil
}

func (u *UserUsecase) FindById(ctx context.Context, id int64) (*model.User, error) {
	log := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	user, err := u.userRepo.FindById(ctx, int64(id))
	if err != nil {
		log.Error("Failed to fetch user by ID: ", err)
		return nil, err
	}

	if user == nil {
		log.Error("User not found")
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (u *UserUsecase) Create(ctx context.Context, in model.CreateUserInput) (token string, err error) {
	logger := logrus.WithFields(logrus.Fields{
		"in": in,
	})

	passwordHashed, err := helper.HashRequestPassword(in.Password)
	if err != nil {
		logger.Error(err)
		return
	}

	newUser, err := u.userRepo.Create(ctx, model.User{
		Username: in.Username,
		Email:    in.Email,
		Password: passwordHashed,
	})

	if err != nil {
		logger.Error(err)
		return
	}

	acceesToken, err := helper.GenerateToken(newUser.Id)
	if err != nil {
		logger.Error(err)
		return
	}

	return acceesToken, nil
}

func (u *UserUsecase) Update(ctx context.Context, id int64, in model.UpdateUserInput) error {
	log := logrus.WithFields(logrus.Fields{
		"id":       id,
		"username": in.Username,
		"email":    in.Email,
	})

	err := v.StructCtx(ctx, in)
	if err != nil {
		log.Error("Validation error:", err)
		return err
	}

	existingUser, err := u.userRepo.FindById(ctx, id)
	if err != nil {
		log.Error("Failed to fetch user: ", err)
		return err
	}
	if existingUser == nil || (existingUser.DeletedAt != nil && !existingUser.DeletedAt.IsZero()) {
		log.Error("User is deleted or does not exist")
		return errors.New("user is deleted or does not exist")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to hash password: ", err)
		return err
	}

	user := model.User{
		Id:        id,
		Username:  in.Username,
		Password:  string(hashedPassword),
		Email:     in.Email,
		UpdatedAt: time.Now(),
	}

	err = u.userRepo.Update(ctx, user)
	if err != nil {
		log.Error("Failed to update user: ", err)
		return err
	}

	return nil
}

func (u *UserUsecase) Delete(ctx context.Context, id int64) error {
	log := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	user, err := u.userRepo.FindById(ctx, id)
	if err != nil {
		log.Error("Failed to find user for deletion: ", err)
		return err
	}

	if user == nil {
		log.Error("User not found")
		return err
	}

	now := time.Now()
	user.DeletedAt = &now

	err = u.userRepo.Delete(ctx, id)
	if err != nil {
		log.Error("Failed to delete user: ", err)
		return err
	}

	log.Info("Successfully deleted user with ID: ", id)
	return nil
}
