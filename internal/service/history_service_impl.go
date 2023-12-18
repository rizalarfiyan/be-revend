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

func (s *historyService) GetAllHistory(ctx context.Context, req request.GetAllHistoryRequest) response.BaseResponsePagination[response.History] {
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

	return response.WithPagination[response.History](content, req.BasePagination)
}

func (s *historyService) GetAllHistoryStatistic(ctx context.Context, req request.GetAllHistoryStatisticRequest) []response.HistoryStatistic {
	timeFrequency := req.BuildTimeFrequency()

	payload := models.AllHistoryStatistic{
		StartDate:     timeFrequency.StartDate,
		EndDate:       timeFrequency.EndDate,
		UserId:        req.UserId,
		TimeFrequency: req.TimeFrequency,
	}

	idx := 0
	var tempArr = make(map[string]int)
	var res []response.HistoryStatistic
	for date := timeFrequency.StartDate; !date.After(timeFrequency.EndDate); date = timeFrequency.CallbackDate(date) {
		name := timeFrequency.CallbackName(date)
		if _, ok := tempArr[name]; !ok {
			tempArr[name] = idx
		}
		res = append(res, response.HistoryStatistic{
			Name:    name,
			Success: 0,
			Failed:  0,
		})
		idx++
	}

	data, err := s.repo.AllHistoryStatistic(ctx, payload)
	exception.PanicIfError(err, true)

	for _, data := range data {
		if !data.Date.Valid {
			continue
		}

		name := timeFrequency.CallbackName(data.Date.Time)
		if idx, ok := tempArr[name]; ok {
			res[idx].Success += data.Success
			res[idx].Failed += data.Failed
		}
	}

	return res
}
