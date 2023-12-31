package service

import (
	"context"

	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/response"
)

type HistoryService interface {
	GetAllHistory(ctx context.Context, req request.GetAllHistoryRequest) response.BaseResponsePagination[response.History]
	GetAllHistoryStatistic(ctx context.Context, req request.GetAllHistoryStatisticRequest) []response.HistoryStatistic
	GetAllHistoryTopPerformance(ctx context.Context, req request.GetAllHistoryTopPerformanceRequest) []response.HistoryTopPerformance
}
