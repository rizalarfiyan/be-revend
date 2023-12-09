package handler

import (
	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/internal/service"
	baseModels "github.com/rizalarfiyan/be-revend/models"
)

type userHandler struct {
	service service.UserService
	conf    *baseModels.Config
}

func NewUserHandler(service service.UserService) UserHandler {
	return &userHandler{
		service: service,
		conf:    config.Get(),
	}
}
