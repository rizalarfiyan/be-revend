package response

import (
	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-revend/internal/sql"
)

type User struct {
	Id          uuid.UUID `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Identity    string    `json:"identity"`
	Role        sql.Role  `json:"role"`
}
