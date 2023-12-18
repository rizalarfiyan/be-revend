package request

import (
	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-revend/constants"
)

type GetAllHistoryRequest struct {
	BasePagination
	DeviceId uuid.UUID `json:"device_id"`
	UserId   uuid.UUID `json:"user_id"`
}

type GetAllHistoryStatisticRequest struct {
	WithTimeFrequency
	UserId uuid.UUID `json:"-"`
}

func (req *GetAllHistoryStatisticRequest) Normalize() {
	req.WithTimeFrequency.Normalize()
}

type GetAllHistoryTopPerformanceRequest struct {
	WithTimeFrequency
	Limit  int       `json:"limit"`
	UserId uuid.UUID `json:"-"`
}

func (req *GetAllHistoryTopPerformanceRequest) Normalize() {
	req.WithTimeFrequency.Normalize()

	if req.Limit > constants.DefaultPageLimitStatistic {
		req.Limit = constants.DefaultPageLimitStatistic
	}
}
