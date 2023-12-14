package request

import "github.com/rizalarfiyan/be-revend/internal/sql"

type GetAllUserRequest struct {
	BasePagination
	Role sql.Role `json:"role"`
}
