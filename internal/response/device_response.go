package response

import (
	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-revend/internal/sql"
	"github.com/rizalarfiyan/be-revend/utils"
)

type Device struct {
	Id       uuid.UUID `json:"id"`
	Token    string    `json:"token"`
	Name     string    `json:"name"`
	Location string    `json:"location"`
}

func (d *Device) FromDB(device sql.Device) {
	d.Id = utils.PGToUUID(device.ID)
	d.Token = device.Token
	d.Name = device.Name
	d.Location = device.Location
}
