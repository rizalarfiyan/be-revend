package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/exception"
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

func (s *userService) GetAllUser(ctx context.Context, req request.GetAllUserRequest) response.BaseResponsePagination[response.User] {
	data, err := s.repo.AllUser(ctx, req)
	exception.PanicIfError(err, true)
	exception.IsNotFound(data, true)

	content := models.ContentPagination[response.User]{
		Count:   data.Count,
		Content: []response.User{},
	}

	for _, val := range data.Content {
		user := response.User{}
		user.FromDB(val)
		content.Content = append(content.Content, user)
	}

	return response.WithPagination[response.User](content, req.BasePagination)
}

func (s *userService) GetUserById(ctx context.Context, userId uuid.UUID) response.User {
	data, err := s.repo.GetUserById(ctx, userId)
	exception.PanicIfError(err, true)
	exception.IsNotFound(data, true)

	user := response.User{}
	user.FromDB(data)
	return user
}

func (s *userService) GetAllDropdownUser(ctx context.Context, req request.BasePagination) response.BaseResponsePagination[response.BaseDropdown] {
	data, err := s.repo.AllDropdownUsers(ctx, req)
	exception.PanicIfError(err, true)
	exception.IsNotFound(data, true)

	content := models.ContentPagination[response.BaseDropdown]{
		Count:   data.Count,
		Content: []response.BaseDropdown{},
	}

	for _, val := range data.Content {
		content.Content = append(content.Content, response.BaseDropdown{
			Key:   utils.FullName(val.FirstName, val.LastName),
			Value: utils.PGToUUID(val.ID).String(),
		})
	}

	return response.WithPagination[response.BaseDropdown](content, req)
}
