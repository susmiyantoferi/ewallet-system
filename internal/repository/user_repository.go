package repository

import "ewallet/internal/entity"

type UserRepository interface {
	Repository[entity.User]
}

type userRepositoryImpl struct {
	repositoryImpl[entity.User]
}

func NewUserRepositoryImpl() UserRepository {
	return &userRepositoryImpl{}
}
