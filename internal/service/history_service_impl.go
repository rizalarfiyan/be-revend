package service

import (
	"context"

	"github.com/rizalarfiyan/be-revend/exception"
	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/repository"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/response"
	"github.com/rizalarfiyan/be-revend/internal/sql"
	"github.com/rizalarfiyan/be-revend/utils"
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

	payload := models.AllHistoryStatistic{
		StartDate:     timeFrequency.StartDate,
		EndDate:       timeFrequency.EndDate,
		UserId:        req.UserId,
		TimeFrequency: req.TimeFrequency,
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

func (s *historyService) GetAllHistoryTopPerformance(ctx context.Context, req request.GetAllHistoryTopPerformanceRequest) []response.HistoryTopPerformance {
	timeFrequency := req.BuildTimeFrequency()

	payload := sql.GetAllHistoryTopPerformanceParams{
		StartDate: utils.PGTimeStamp(timeFrequency.StartDate),
		EndDate:   utils.PGTimeStamp(timeFrequency.EndDate),
		Limit:     int32(req.Limit),
	}

	data, err := s.repo.AllHistoryTopPerformance(ctx, payload)
	exception.PanicIfError(err, true)

	res := []response.HistoryTopPerformance{}
	for _, v := range data {
		resPayload := response.HistoryTopPerformance{
			FirstName:   v.FirstName,
			PhoneNumber: utils.CensorPhoneNumber(v.PhoneNumber),
			Success:     v.Success,
			Failed:      v.Failed,
		}

		if v.UserID.Valid {
			resPayload.IsMe = utils.PGToUUID(v.UserID) == req.UserId
		}

		if v.LastName.Valid {
			resPayload.LastName = v.LastName.String
		}

		res = append(res, resPayload)
	}

	return res
}
