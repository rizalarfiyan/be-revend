package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-revend/internal/sql"
	"github.com/rizalarfiyan/be-revend/utils"
)

type History struct {
	Id      uuid.UUID     `json:"id"`
	Success int32         `json:"success"`
	Failed  int32         `json:"failed"`
	User    HistoryUser   `json:"user"`
	Device  HistoryDevice `json:"device"`
	Date    time.Time     `json:"date"`
}

type HistoryUser struct {
	Id        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

type HistoryDevice struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (h *History) FromDB(history sql.GetAllHistoryRow) {
	h.Id = utils.PGToUUID(history.ID)
	h.Success = history.Success
	h.Failed = history.Failed
	h.User.Id = utils.PGToUUID(history.UserID)
	h.User.FirstName = history.FirstName
	if history.LastName.Valid {
		h.User.LastName = history.LastName.String
	}
	h.Device.Id = utils.PGToUUID(history.DeviceID)
	h.Device.Name = history.DeviceName
	if history.CreatedAt.Valid {
		h.Date = history.CreatedAt.Time
	}
}

type HistoryStatistic struct {
	Name    string `json:"name"`
	Success int64  `json:"success"`
	Failed  int64  `json:"failed"`
}

type HistoryTopPerformance struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Success     int64  `json:"success"`
	Failed      int64  `json:"failed"`
	IsMe        bool   `json:"is_me"`
}
