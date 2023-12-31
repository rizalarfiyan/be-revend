package request

import "github.com/google/uuid"

type CreateDeviceRequest struct {
	Name     string `json:"name" field:"Name" validate:"required,min=3,max=50" example:"Revend AM"`
	Location string `json:"location" field:"Location" validate:"required,min=5,max=150" example:"Revend Universitas Amikom Yogyakarta"`
}

type UpdateDeviceRequest struct {
	Id uuid.UUID `json:"-"`
	CreateDeviceRequest
}
