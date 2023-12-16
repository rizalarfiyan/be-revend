package request

import (
	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-revend/internal/sql"
)

type CreateUserRequest struct {
	FirstName   string   `json:"first_name" field:"FirstName" validate:"required,min=3,max=100" example:"Rizal"`
	LastName    string   `json:"last_name" field:"LastName" validate:"omitempty,min=3,max=100" example:"Arfiyan"`
	PhoneNumber string   `json:"phone_number" field:"PhoneNumber" validate:"phone_number" example:"62895377233002"`
	GoogleId    string   `json:"google_id" field:"Google Id" validate:"omitempty,min=8,max=30" example:"1234567890"`
	Identity    string   `json:"identity" field:"Identity" validate:"required,min=8,max=50" example:"1234567890"`
	RawRole     string   `json:"role" field:"Role" validate:"required" example:"guest"`
	Role        sql.Role `json:"-"`
}

type UpdateUserRequest struct {
	Id uuid.UUID `json:"-"`
	CreateUserRequest
}

type GetAllUserRequest struct {
	BasePagination
	Role sql.Role `json:"role"`
}
