package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-revend/constants"
)

type AllHistoryStatistic struct {
	StartDate     time.Time
	EndDate       time.Time
	UserId        uuid.UUID
	TimeFrequency constants.FilterTimeFrequency
}
