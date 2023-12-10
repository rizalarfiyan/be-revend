package utils

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func UUID(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: u,
		Valid: true,
	}
}

func ToUUID(t pgtype.UUID) uuid.UUID {
	if !t.Valid {
		return uuid.Nil
	}

	return t.Bytes
}
