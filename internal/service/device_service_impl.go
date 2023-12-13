package service

import (
	"context"

	"github.com/rizalarfiyan/be-revend/exception"
	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/repository"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/response"
	"github.com/rizalarfiyan/be-revend/utils"
)

type deviceService struct {
	repo repository.DeviceRepository
}

func NewDeviceService(repo repository.DeviceRepository) DeviceService {
	return &deviceService{
		repo: repo,
	}
}

func (s *deviceService) GetAllDevice(ctx context.Context, req request.BasePagination) response.BaseResponsePagination[response.Device] {
	data, err := s.repo.AllDevice(ctx, req)
	exception.PanicIfError(err, true)
	exception.IsNotFound(data, true)

	content := models.ContentPagination[response.Device]{
		Count:   data.Count,
		Content: []response.Device{},
	}

	for _, val := range data.Content {
		user := response.Device{}
		user.FromDB(val)
		content.Content = append(content.Content, user)
	}

	return response.WithPagination[response.Device](content, req)
}

func (s *deviceService) GetAllDropdownDevice(ctx context.Context, req request.BasePagination) response.BaseResponsePagination[response.BaseDropdown] {
	data, err := s.repo.AllDropdownDevice(ctx, req)
	exception.PanicIfError(err, true)
	exception.IsNotFound(data, true)

	content := models.ContentPagination[response.BaseDropdown]{
		Count:   data.Count,
		Content: []response.BaseDropdown{},
	}

	for _, val := range data.Content {
		content.Content = append(content.Content, response.BaseDropdown{
			Key:   val.Name,
			Value: utils.PGToUUID(val.ID).String(),
		})
	}

	return response.WithPagination[response.BaseDropdown](content, req)
}
