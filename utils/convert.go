package utils

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ConvertUUID(t pgtype.UUID) uuid.UUID {
	if !t.Valid {
		return uuid.Nil
	}

	return t.Bytes
}
