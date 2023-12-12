package repository

import (
	"context"

	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/sql"
)

type HistoryRepository interface {
	AllHistory(ctx context.Context, req request.BasePagination) (*models.ContentPagination[sql.GetAllHistoryRow], error)
}
