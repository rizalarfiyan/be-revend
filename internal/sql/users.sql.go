// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: users.sql

package sql

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countAllUsers = `-- name: CountAllUsers :one
SELECT count(*) FROM users
`

func (q *Queries) CountAllUsers(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countAllUsers)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createUser = `-- name: CreateUser :exec
INSERT INTO users (first_name, last_name, phone_number, google_id, identity)
VALUES ($1, $2, $3, $4, $5)
`

type CreateUserParams struct {
	FirstName   string
	LastName    pgtype.Text
	PhoneNumber string
	GoogleID    pgtype.Text
	Identity    string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.Exec(ctx, createUser,
		arg.FirstName,
		arg.LastName,
		arg.PhoneNumber,
		arg.GoogleID,
		arg.Identity,
	)
	return err
}

const getAllNameUsers = `-- name: GetAllNameUsers :many
SELECT id, first_name, last_name FROM users
`

type GetAllNameUsersRow struct {
	ID        pgtype.UUID
	FirstName string
	LastName  pgtype.Text
}

func (q *Queries) GetAllNameUsers(ctx context.Context) ([]GetAllNameUsersRow, error) {
	rows, err := q.db.Query(ctx, getAllNameUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllNameUsersRow
	for rows.Next() {
		var i GetAllNameUsersRow
		if err := rows.Scan(&i.ID, &i.FirstName, &i.LastName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllUsers = `-- name: GetAllUsers :many
SELECT id, first_name, last_name, phone_number, google_id, identity, role, created_at, updated_at, deleted_by, deleted_at FROM users
`

func (q *Queries) GetAllUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, getAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.PhoneNumber,
			&i.GoogleID,
			&i.Identity,
			&i.Role,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedBy,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByGoogleId = `-- name: GetUserByGoogleId :one
SELECT id, first_name, last_name, phone_number, google_id, identity, role, created_at, updated_at, deleted_by, deleted_at FROM users
WHERE google_id = $1 LIMIT 1
`

func (q *Queries) GetUserByGoogleId(ctx context.Context, googleID pgtype.Text) (User, error) {
	row := q.db.QueryRow(ctx, getUserByGoogleId, googleID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.PhoneNumber,
		&i.GoogleID,
		&i.Identity,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedBy,
		&i.DeletedAt,
	)
	return i, err
}

const getUserByGoogleIdOrPhoneNumber = `-- name: GetUserByGoogleIdOrPhoneNumber :one
SELECT id, first_name, last_name, phone_number, google_id, identity, role, created_at, updated_at, deleted_by, deleted_at FROM users
WHERE google_id = $1 OR phone_number = $2 LIMIT 1
`

type GetUserByGoogleIdOrPhoneNumberParams struct {
	GoogleID    pgtype.Text
	PhoneNumber string
}

func (q *Queries) GetUserByGoogleIdOrPhoneNumber(ctx context.Context, arg GetUserByGoogleIdOrPhoneNumberParams) (User, error) {
	row := q.db.QueryRow(ctx, getUserByGoogleIdOrPhoneNumber, arg.GoogleID, arg.PhoneNumber)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.PhoneNumber,
		&i.GoogleID,
		&i.Identity,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedBy,
		&i.DeletedAt,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, first_name, last_name, phone_number, google_id, identity, role, created_at, updated_at, deleted_by, deleted_at FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserById(ctx context.Context, id pgtype.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.PhoneNumber,
		&i.GoogleID,
		&i.Identity,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedBy,
		&i.DeletedAt,
	)
	return i, err
}

const getUserByIdentity = `-- name: GetUserByIdentity :one
SELECT id, first_name, last_name, phone_number, google_id, identity, role, created_at, updated_at, deleted_by, deleted_at FROM users
WHERE identity = $1 LIMIT 1
`

func (q *Queries) GetUserByIdentity(ctx context.Context, identity string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByIdentity, identity)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.PhoneNumber,
		&i.GoogleID,
		&i.Identity,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedBy,
		&i.DeletedAt,
	)
	return i, err
}

const getUserByPhoneNumber = `-- name: GetUserByPhoneNumber :one
SELECT id, first_name, last_name, phone_number, google_id, identity, role, created_at, updated_at, deleted_by, deleted_at FROM users
WHERE phone_number = $1 LIMIT 1
`

func (q *Queries) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByPhoneNumber, phoneNumber)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.PhoneNumber,
		&i.GoogleID,
		&i.Identity,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedBy,
		&i.DeletedAt,
	)
	return i, err
}

const toggleDeleteUser = `-- name: ToggleDeleteUser :exec
UPDATE users SET
deleted_by = CASE WHEN deleted_by IS NULL THEN $1::UUID ELSE NULL END,
deleted_at = CASE WHEN deleted_at IS NULL THEN CURRENT_TIMESTAMP ELSE NULL
END WHERE id = $2
`

type ToggleDeleteUserParams struct {
	DeletedBy pgtype.UUID
	ID        pgtype.UUID
}

func (q *Queries) ToggleDeleteUser(ctx context.Context, arg ToggleDeleteUserParams) error {
	_, err := q.db.Exec(ctx, toggleDeleteUser, arg.DeletedBy, arg.ID)
	return err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
SET first_name = $1, last_name = $2, phone_number = $3, google_id = $4, identity = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $6
`

type UpdateUserParams struct {
	FirstName   string
	LastName    pgtype.Text
	PhoneNumber string
	GoogleID    pgtype.Text
	Identity    string
	ID          pgtype.UUID
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser,
		arg.FirstName,
		arg.LastName,
		arg.PhoneNumber,
		arg.GoogleID,
		arg.Identity,
		arg.ID,
	)
	return err
}

const updateUserProfile = `-- name: UpdateUserProfile :exec
UPDATE users
SET first_name = $1, last_name = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $3
`

type UpdateUserProfileParams struct {
	FirstName string
	LastName  pgtype.Text
	ID        pgtype.UUID
}

func (q *Queries) UpdateUserProfile(ctx context.Context, arg UpdateUserProfileParams) error {
	_, err := q.db.Exec(ctx, updateUserProfile, arg.FirstName, arg.LastName, arg.ID)
	return err
}
