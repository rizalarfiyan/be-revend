package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/repository"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/response"
	baseModels "github.com/rizalarfiyan/be-revend/models"
	"github.com/rizalarfiyan/be-revend/utils"
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

func (s *userService) GetAllUser(ctx context.Context, req request.BasePagination) response.BaseResponsePagination[response.User] {
	data, err := s.repo.AllUser(ctx, req)
	utils.PanicIfError(err, true)
	utils.IsNotFound(data, true)

	content := models.ContentPagination[response.User]{
		Count:   data.Count,
		Content: []response.User{},
	}

	for _, val := range data.Content {
		user := response.User{
			Id:          utils.ToUUID(val.ID),
			FirstName:   val.FirstName,
			PhoneNumber: val.PhoneNumber,
			Identity:    val.Identity,
			Role:        val.Role,
		}

		if val.LastName.Valid {
			user.LastName = val.LastName.String
		}

		content.Content = append(content.Content, user)
	}

	return response.WithPagination[response.User](content, req)
}

func (s *userService) GetUserById(ctx context.Context, userId uuid.UUID) response.User {
	data, err := s.repo.GetUserById(ctx, userId)
	utils.PanicIfError(err, true)
	utils.IsNotFound(data, true)

	user := response.User{
		Id:          utils.ToUUID(data.ID),
		FirstName:   data.FirstName,
		PhoneNumber: data.PhoneNumber,
		Identity:    data.Identity,
		Role:        data.Role,
	}

	if data.LastName.Valid {
		user.LastName = data.LastName.String
	}

	return user
}
