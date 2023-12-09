package service

import (
	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/internal/repository"
	baseModels "github.com/rizalarfiyan/be-revend/models"
)

type userService struct {
	repo repository.UserRepository
	conf *baseModels.Config
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
		conf: config.Get(),
	}
}
