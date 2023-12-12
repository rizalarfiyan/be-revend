package response

import (
	"github.com/rizalarfiyan/be-revend/internal/sql"
)

type Device struct {
	Token    string `json:"token"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

func (u *Device) FromDB(device sql.Device) {
	u.Token = device.Token
	u.Name = device.Name
	u.Location = device.Location
}
