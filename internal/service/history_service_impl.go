package service

import (
	"context"

	"github.com/rizalarfiyan/be-revend/exception"
	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/repository"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/response"
)

type historyService struct {
	repo repository.HistoryRepository
}

func NewHistoryService(repo repository.HistoryRepository) HistoryService {
	return &historyService{
		repo: repo,
	}
}

func (s *historyService) GetAllHistory(ctx context.Context, req request.BasePagination) response.BaseResponsePagination[response.History] {
	data, err := s.repo.AllHistory(ctx, req)
	exception.PanicIfError(err, true)
	exception.IsNotFound(data, true)

	content := models.ContentPagination[response.History]{
		Count:   data.Count,
		Content: []response.History{},
	}

	for _, val := range data.Content {
		user := response.History{}
		user.FromDB(val)
		content.Content = append(content.Content, user)
	}

	return response.WithPagination[response.History](content, req)
}
