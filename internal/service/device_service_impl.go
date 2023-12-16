package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-revend/exception"
	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/repository"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/response"
	"github.com/rizalarfiyan/be-revend/internal/sql"
	"github.com/rizalarfiyan/be-revend/utils"
	"github.com/segmentio/ksuid"
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

func (s *deviceService) GetAllDropdownDevice(ctx context.Context, req request.AllDropdownDeviceRequest) response.BaseResponsePagination[response.BaseDropdown] {
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

	return response.WithPagination[response.BaseDropdown](content, req.BasePagination)
}

func (s *deviceService) CreateDevice(ctx context.Context, req request.CreateDeviceRequest) {
	err := s.repo.CreateDevice(ctx, sql.CreateDeviceParams{
		Name:     req.Name,
		Location: req.Location,
		Token:    ksuid.New().String(),
	})
	exception.PanicIfError(err, true)
}

func (s *deviceService) UpdateDevice(ctx context.Context, req request.UpdateDeviceRequest) {
	err := s.repo.UpdateDevice(ctx, sql.UpdateDeviceParams{
		ID:       utils.PGUUID(req.ID),
		Name:     req.Name,
		Location: req.Location,
	})
	exception.PanicIfError(err, true)
}

func (s *deviceService) ToggleDeleteDevice(ctx context.Context, deviceId, userId uuid.UUID) {
	err := s.repo.ToggleDeleteDevice(ctx, sql.ToggleDeleteDeviceParams{
		ID:        utils.PGUUID(deviceId),
		DeletedBy: utils.PGUUID(userId),
	})
	exception.PanicIfError(err, true)
}
