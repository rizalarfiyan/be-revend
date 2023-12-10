package service

import (
	"context"

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

func (s *userService) AllUser(ctx context.Context, req request.BasePagination) response.BaseResponsePagination[response.User] {
	data, err := s.repo.AllUser(ctx, req)
	utils.PanicIfError(err, true)
	utils.IsNotFound(data, true)

	content := models.ContentPagination[response.User]{
		Count:   data.Count,
		Content: []response.User{},
	}

	for _, val := range data.Content {
		user := response.User{
			FirstName:   val.FirstName,
			PhoneNumber: val.PhoneNumber,
		}

		if val.LastName.Valid {
			user.LastName = val.LastName.String
		}

		content.Content = append(content.Content, user)
	}

	return response.WithPagination[response.User](content, req)
}
