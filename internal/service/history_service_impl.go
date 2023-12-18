package service

import (
	"context"
	"time"

	"github.com/rizalarfiyan/be-revend/constants"
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
	var startDate, endDate time.Time
	var callbackName func(time.Time) string
	var callbackDate func(time.Time) time.Time

	now := time.Now()
	startDateToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	switch req.TimeFrequency {
	case constants.FilterTimeFrequencyWeek:
		first := now.AddDate(0, 0, -int(now.Weekday()-time.Sunday))
		startDate = time.Date(first.Year(), first.Month(), first.Day(), 0, 0, 0, 0, first.Location())
		end := now.AddDate(0, 0, int(time.Saturday-now.Weekday()))
		endDate = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())
		callbackName = func(val time.Time) string {
			return val.Format(time.DateOnly)
		}
		callbackDate = func(val time.Time) time.Time {
			return val.AddDate(0, 0, 1)
		}
	case constants.FilterTimeFrequencyMonth:
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()).Add(-time.Second)
		callbackName = func(val time.Time) string {
			return val.Format(time.DateOnly)
		}
		callbackDate = func(val time.Time) time.Time {
			return val.AddDate(0, 0, 1)
		}
	case constants.FilterTimeFrequencyQuarter:
		date := startDateToday.AddDate(0, -6, 0)
		startDate = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
		endDate = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()).Add(-time.Second)
		callbackName = func(val time.Time) string {
			return val.Format("Jan 2006")
		}
		callbackDate = func(val time.Time) time.Time {
			return val.AddDate(0, 1, 0)
		}
	case constants.FilterTimeFrequencyYear:
		date := startDateToday.AddDate(-1, 0, 0)
		startDate = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
		endDate = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()).Add(-time.Second)
		callbackName = func(val time.Time) string {
			return val.Format("Jan 2006")
		}
		callbackDate = func(val time.Time) time.Time {
			return val.AddDate(0, 1, 0)
		}
	default:
		startDate = startDateToday
		endDate = startDateToday.AddDate(0, 0, 1).Add(-time.Second)
		callbackName = func(val time.Time) string {
			return val.Format(time.TimeOnly)
		}
		callbackDate = func(val time.Time) time.Time {
			return val.Add(1 * time.Hour)
		}
	}

	payload := models.AllHistoryStatistic{
		StartDate:     startDate,
		EndDate:       endDate,
		UserId:        req.UserId,
		TimeFrequency: req.TimeFrequency,
	}

	idx := 0
	var tempArr = make(map[string]int)
	var res []response.HistoryStatistic
	for date := startDate; !date.After(endDate); date = callbackDate(date) {
		name := callbackName(date)
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

		name := callbackName(data.Date.Time)
		if idx, ok := tempArr[name]; ok {
			res[idx].Success += data.Success
			res[idx].Failed += data.Failed
		}
	}

	return res
}
