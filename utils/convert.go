package utils

import (
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func Str(str string) string {
	return strings.ToLower(strings.TrimSpace(str))
}

func PGUUID(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: u,
		Valid: true,
	}
}

func PGText(str string) pgtype.Text {
	return pgtype.Text{
		String: str,
		Valid:  true,
	}
}

func PGToUUID(t pgtype.UUID) uuid.UUID {
	if !t.Valid {
		return uuid.Nil
	}

	return t.Bytes
}
