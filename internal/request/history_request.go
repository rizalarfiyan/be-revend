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
	UserId        uuid.UUID                     `json:"user_id"`
	TimeFrequency constants.FilterTimeFrequency `json:"time_range"`
}
