package response

import (
	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-revend/internal/sql"
	"github.com/rizalarfiyan/be-revend/utils"
)

type User struct {
	Id          uuid.UUID `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Identity    string    `json:"identity"`
	Role        sql.Role  `json:"role"`
	GoogleId    string    `json:"google_id"`
	IsDeleted   bool      `json:"is_deleted"`
}

func (u *User) FromDB(user sql.User) {
	u.Id = utils.PGToUUID(user.ID)
	u.FirstName = user.FirstName
	u.PhoneNumber = user.PhoneNumber
	u.Identity = user.Identity
	u.IsDeleted = user.DeletedAt.Valid

	u.Role = user.Role
	if user.LastName.Valid {
		u.LastName = user.LastName.String
	}

	if user.GoogleID.Valid {
		u.GoogleId = user.GoogleID.String
	}
}

type BindGoogleUserProfile struct {
	Url string `json:"url"`
}
