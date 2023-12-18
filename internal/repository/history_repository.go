package repository

import (
	"context"

	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/sql"
)

type HistoryRepository interface {
	AllHistory(ctx context.Context, req request.GetAllHistoryRequest) (*models.ContentPagination[sql.GetAllHistoryRow], error)
	CreateHistory(ctx context.Context, payload sql.CreateHistoryParams) error
	AllHistoryStatistic(ctx context.Context, req models.AllHistoryStatistic) ([]sql.GetAllHistoryStatisticRow, error)
	AllHistoryTopPerformance(ctx context.Context, payload sql.GetAllHistoryTopPerformanceParams) ([]sql.GetAllHistoryTopPerformanceRow, error)
}
