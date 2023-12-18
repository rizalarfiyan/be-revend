package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/exception"
	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/repository"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/response"
	"github.com/rizalarfiyan/be-revend/internal/sql"
	baseModels "github.com/rizalarfiyan/be-revend/models"
	"github.com/rizalarfiyan/be-revend/utils"
	"github.com/segmentio/ksuid"
)

type userService struct {
	repo        repository.UserRepository
	conf        *baseModels.Config
	authService AuthService
}

func NewUserService(repo repository.UserRepository, authService AuthService) UserService {
	return &userService{
		repo:        repo,
		conf:        config.Get(),
		authService: authService,
	}
}

func (s *userService) GetAllUser(ctx context.Context, req request.GetAllUserRequest) response.BaseResponsePagination[response.User] {
	data, err := s.repo.AllUser(ctx, req)
	exception.PanicIfError(err, true)
	exception.IsNotFound(data, true)

	content := models.ContentPagination[response.User]{
		Count:   data.Count,
		Content: []response.User{},
	}

	for _, val := range data.Content {
		user := response.User{}
		user.FromDB(val)
		content.Content = append(content.Content, user)
	}

	return response.WithPagination[response.User](content, req.BasePagination)
}

func (s *userService) GetUserById(ctx context.Context, userId uuid.UUID) response.User {
	data, err := s.repo.GetUserById(ctx, userId)
	exception.PanicIfError(err, false)
	exception.IsNotFound(data, false)

	user := response.User{}
	user.FromDB(data)
	return user
}

func (s *userService) GetAllDropdownUser(ctx context.Context, req request.BasePagination) response.BaseResponsePagination[response.BaseDropdown] {
	data, err := s.repo.AllDropdownUsers(ctx, req)
	exception.PanicIfError(err, true)
	exception.IsNotFound(data, true)

	content := models.ContentPagination[response.BaseDropdown]{
		Count:   data.Count,
		Content: []response.BaseDropdown{},
	}

	for _, val := range data.Content {
		content.Content = append(content.Content, response.BaseDropdown{
			Key:   utils.FullName(val.FirstName, val.LastName),
			Value: utils.PGToUUID(val.ID).String(),
		})
	}

	return response.WithPagination[response.BaseDropdown](content, req)
}

func (s *userService) handleErrorUniqueUser(err error) {
	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return
	}

	if pgErr.Code != pgerrcode.UniqueViolation {
		return
	}

	switch pgErr.ConstraintName {
	case "users_phone_number_key":
		exception.ErrorManualValidation("phone_number", "Phone Number already exist.")
	case "users_google_id_key":
		exception.ErrorManualValidation("google_id", "Google Id already exist.")
	case "users_identity_key":
		exception.ErrorManualValidation("identity", "Identity already exist.")
	}
}

func (s *userService) CreateUser(ctx context.Context, req request.CreateUserRequest) {
	payload := sql.CreateUserParams{
		FirstName:   req.FirstName,
		PhoneNumber: req.PhoneNumber,
		Identity:    req.Identity,
	}

	if !utils.IsEmpty(req.LastName) {
		payload.LastName = utils.PGText(req.LastName)
	}

	if !utils.IsEmpty(req.GoogleId) {
		payload.GoogleID = utils.PGText(req.GoogleId)
	}

	err := s.repo.CreateUser(ctx, payload)
	s.handleErrorUniqueUser(err)
	exception.PanicIfError(err, false)
}

func (s *userService) UpdateUser(ctx context.Context, req request.UpdateUserRequest) {
	payload := sql.UpdateUserParams{
		ID:          utils.PGUUID(req.Id),
		FirstName:   req.FirstName,
		PhoneNumber: req.PhoneNumber,
		Identity:    req.Identity,
	}

	if !utils.IsEmpty(req.LastName) {
		payload.LastName = utils.PGText(req.LastName)
	}

	if !utils.IsEmpty(req.GoogleId) {
		payload.GoogleID = utils.PGText(req.GoogleId)
	}

	err := s.repo.UpdateUser(ctx, payload)
	s.handleErrorUniqueUser(err)
	exception.PanicIfError(err, true)
}

func (s *userService) ToggleDeleteUser(ctx context.Context, userId, currentUserId uuid.UUID) {
	err := s.repo.ToggleDeleteUser(ctx, sql.ToggleDeleteUserParams{
		ID:        utils.PGUUID(userId),
		DeletedBy: utils.PGUUID(currentUserId),
	})
	exception.PanicIfError(err, false)
}

func (s *userService) UpdateUserProfile(ctx context.Context, req request.UpdateUserProfileRequest) {
	payload := sql.UpdateUserProfileParams{
		ID:        utils.PGUUID(req.Id),
		FirstName: req.FirstName,
	}

	if !utils.IsEmpty(req.LastName) {
		payload.LastName = utils.PGText(req.LastName)
	}

	err := s.repo.UpdateUserProfile(ctx, payload)
	exception.PanicIfError(err, false)
}

func (s *userService) DeleteGoogleUserProfile(ctx context.Context, userId uuid.UUID) {
	err := s.repo.DeleteGoogleUserProfile(ctx, userId)
	exception.PanicIfError(err, false)
}

func (s *userService) BindGoogleUserProfile(ctx context.Context, userId uuid.UUID) response.BindGoogleUserProfile {
	user, err := s.repo.GetUserById(ctx, userId)
	exception.PanicIfErrorWithoutNoSqlResult(err, false)
	exception.IsNotProcess(user, false)

	if user.GoogleID.Valid {
		exception.IsNotProcessMessage(nil, "Account is already bind.", false)
	}

	token := ksuid.New().String()
	err = s.repo.CreateBindGoogleFresh(ctx, token, userId)
	exception.PanicIfError(err, false)

	return response.BindGoogleUserProfile{
		Url: s.authService.Google(token),
	}
}
