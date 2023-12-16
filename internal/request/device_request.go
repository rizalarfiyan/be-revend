package request

import "github.com/google/uuid"

type AllDropdownDeviceRequest struct {
	BasePagination
	HideDeleted bool `json:"-"`
}

type CreateDeviceRequest struct {
	Name     string `json:"name" field:"Name" validate:"required" example:"Revend AM"`
	Location string `json:"location" field:"Location" validate:"required" example:"Revend Universitas Amikom Yogyakarta"`
}

type UpdateDeviceRequest struct {
	ID uuid.UUID `json:"-"`
	CreateDeviceRequest
}
