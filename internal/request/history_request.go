package request

import (
	"github.com/google/uuid"
)

type GetAllHistoryRequest struct {
	BasePagination
	DeviceId uuid.UUID `json:"device_id"`
	UserId   uuid.UUID `json:"user_id"`
}

type GetAllHistoryStatisticRequest struct {
	WithTimeFrequency
	UserId uuid.UUID `json:"user_id"`
}

func (req *GetAllHistoryStatisticRequest) Normalize() {
	req.WithTimeFrequency.Normalize()
}
