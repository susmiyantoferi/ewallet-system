package service

import (
	"context"
	"errors"
	"ewallet/internal/dto"
	"ewallet/internal/entity"
	"ewallet/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserService interface {
	Create(c context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	Update(c context.Context, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(c context.Context, id uuid.UUID) error
	GetByID(c context.Context, id uuid.UUID) (*dto.UserResponse, error)
}

type userServiceImpl struct {
	UserRepo repository.UserRepository
	Db       *gorm.DB
	Validate *validator.Validate
	Log      *logrus.Logger
}

func NewUserServiceImpl(userRepo repository.UserRepository, db *gorm.DB, validate *validator.Validate, log *logrus.Logger) UserService {
	return &userServiceImpl{
		UserRepo: userRepo,
		Db:       db,
		Validate: validate,
		Log:      log,
	}
}

func (s *userServiceImpl) Create(c context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	if err := s.Validate.Struct(req); err != nil {
		s.Log.WithError(err).Error("Validation failed")
		return nil, err
	}

	user := entity.User{
		Name:    req.Name,
		Address: req.Address,
	}

	if err := s.UserRepo.Create(s.Db.WithContext(c), &user); err != nil {
		s.Log.WithError(err).Error("Failed to create user")
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID.String(),
		Name:      user.Name,
		Address:   user.Address,
		CreatedAt: user.CreatedAt,
		UpdatedAt: *user.UpdatedAt,
	}, nil
}

func (s *userServiceImpl) Update(c context.Context, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	db := s.Db.WithContext(c)

	if err := s.Validate.Struct(req); err != nil {
		s.Log.WithError(err).Error("Validation failed")
		return nil, err
	}

	user, err := s.UserRepo.GetByID(db, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.WithError(err).Error("get user by id not found")
			return nil, gorm.ErrRecordNotFound
		}
		s.Log.WithError(err).Error("failed get user by id")
		return nil, err
	}

	u := &entity.User{
		Name:    req.Name,
		Address: req.Address,
	}

	if err := s.UserRepo.Update(db, u, user.ID); err != nil {
		s.Log.WithError(err).Error("failed to update user")
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID.String(),
		Name:      u.Name,
		Address:   u.Address,
		CreatedAt: user.CreatedAt,
		UpdatedAt: *user.UpdatedAt,
	}, nil
}

func (s *userServiceImpl) Delete(c context.Context, id uuid.UUID) error {
	db := s.Db.WithContext(c)

	user, err := s.UserRepo.GetByID(db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.WithError(err).Error("get user by id not found")
			return gorm.ErrRecordNotFound
		}
		s.Log.WithError(err).Error("failed get user by id")
		return err
	}

	if err := s.UserRepo.Delete(db, user.ID); err != nil {
		s.Log.WithError(err).Error("failed delte user")
		return err
	}

	return nil
}

func (s *userServiceImpl) GetByID(c context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	db := s.Db.WithContext(c)

	user, err := s.UserRepo.GetByID(db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Log.WithError(err).Error("get user by id not found")
			return nil, gorm.ErrRecordNotFound
		}
		s.Log.WithError(err).Error("failed get user by id")
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID.String(),
		Name:      user.Name,
		Address:   user.Address,
		CreatedAt: user.CreatedAt,
		UpdatedAt: *user.UpdatedAt,
	}, nil

}
