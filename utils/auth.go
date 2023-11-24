package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/constants"
	"github.com/rizalarfiyan/be-revend/internal/sql"
	baseModels "github.com/rizalarfiyan/be-revend/models"
)

func ValidateUser(ctx *fiber.Ctx) (*baseModels.AuthToken, error) {
	user, isValidToken := ctx.Locals(constants.KeyLocalsUser).(*jwt.Token)
	if !isValidToken {
		return nil, fmt.Errorf("could not extract user from context")
	}

	claims, isValidMapClaims := user.Claims.(jwt.MapClaims)
	if !isValidMapClaims {
		return nil, fmt.Errorf("could not extract claims from JWT token")
	}

	firstName, firstNameOk := claims["first_name"].(string)
	lastName, lastNameOk := claims["last_name"].(string)
	phoneNumber, phoneNumberOk := claims["phone_number"].(string)
	role, roleOk := claims["role"].(sql.Role)
	if !firstNameOk || !lastNameOk || !phoneNumberOk || !roleOk {
		return nil, fmt.Errorf("one or more claims are missing or have incorrect types")
	}

	res := baseModels.AuthToken{
		FirstName:   firstName,
		LastName:    lastName,
		PhoneNumber: phoneNumber,
		Role:        role,
	}

	return &res, nil
}

func GetUser(ctx *fiber.Ctx) baseModels.AuthToken {
	user, err := ValidateUser(ctx)
	if err != nil {
		return baseModels.AuthToken{}
	}

	return *user
}

func GenerateJwtToken(claims jwt.Claims) (string, error) {
	conf := config.Get()
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte(conf.JWT.Secret))
	if err != nil {
		return "", err
	}

	return token, nil
}
