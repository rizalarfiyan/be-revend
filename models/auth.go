package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-revend/internal/sql"
)

type AuthToken struct {
	Id          uuid.UUID `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Role        sql.Role  `json:"role"`
}

type AuthTokenClaims struct {
	AuthToken
	jwt.RegisteredClaims
}
