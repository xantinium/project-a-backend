package core_database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	api_auth "github.com/xantinium/project-a-backend/api/auth"
)

type UserStruct struct {
	Id              int                    `db:"id"`
	FirstName       string                 `db:"first_name"`
	LastName        string                 `db:"last_name"`
	AvatarId        *string                `db:"avatar_id"`
	OAuthSerive     api_auth.OAuthServices `db:"oauth_service"`
	YandexProfileId *string                `db:"yandex_profile_id"`
	GoogleProfileId *string                `db:"google_profile_id"`
	CreatedAt       time.Time              `db:"created_at"`
	UpdatedAt       time.Time              `db:"updated_at"`
}

type GetUsersOptions struct {
	Id        *int    `db:"id"`
	FirstName *string `db:"first_name"`
	LastName  *string `db:"last_name"`
}

func (dbClient *DatabaseClient) GetUsers(options *GetUsersOptions) (*[]UserStruct, error) {
	queryStr := "SELECT * FROM users"
	q := CreateColumnsQuery(options)
	if q != "" {
		queryStr += fmt.Sprintf(" WHERE %s", q)
	}

	query, err := dbClient.p.Query(context.Background(), queryStr)
	if err != nil {
		return nil, err
	}

	users, err := pgx.CollectRows(query, pgx.RowToStructByName[UserStruct])
	if err != nil {
		return nil, err
	}

	return &users, nil
}

type GetUserByServiceProfileIdOptions struct {
	ServiceProfileId string
	Service          api_auth.OAuthServices
}

func (dbClient *DatabaseClient) GetUserByServiceProfileId(options *GetUserByServiceProfileIdOptions) (*UserStruct, error) {
	var serviceName string
	switch options.Service {
	case api_auth.OAuthServicesYANDEX:
		serviceName = "yandex_profile_id"
	case api_auth.OAuthServicesGOOGLE:
		serviceName = "google_profile_id"
	}

	query, err := dbClient.p.Query(
		context.Background(),
		fmt.Sprintf("SELECT * FROM users WHERE oauth_service = $1 AND %s = $2", serviceName),
		options.Service,
		options.ServiceProfileId,
	)
	if err != nil {
		return nil, err
	}

	user, err := pgx.CollectOneRow(query, pgx.RowToStructByName[UserStruct])
	if err != nil {
		return nil, err
	}

	return &user, nil
}

type CreateUserOptions struct {
	FirstName       string
	LastName        string
	AvatarId        *string
	OAuthSerive     api_auth.OAuthServices
	YandexProfileId *string
	GoogleProfileId *string
}

func (dbClient *DatabaseClient) CreateUser(options *CreateUserOptions) (int, error) {
	currentTime := time.Now().UTC()
	var createdUserId int

	err := dbClient.p.QueryRow(
		context.Background(),
		"INSERT INTO users VALUES (DEFAULT, $1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		options.FirstName,
		options.LastName,
		options.AvatarId,
		options.OAuthSerive,
		options.YandexProfileId,
		options.GoogleProfileId,
		currentTime,
		currentTime,
	).Scan(&createdUserId)

	return createdUserId, err
}

type UpdateUserOptionsFields struct {
	FirstName *string  `db:"first_name"`
	LastName  *string  `db:"last_name"`
	AvatarId  **string `db:"avatar_id"`
}

type UpdateUserOptions struct {
	Id     int
	Fields *UpdateUserOptionsFields
}

func (dbClient *DatabaseClient) UpdateUser(options *UpdateUserOptions) error {
	currentTime := time.Now().UTC()

	query := fmt.Sprintf("UPDATE users SET updated_at = $1, %s WHERE id = $2", CreateColumnsQuery(options.Fields))

	_, err := dbClient.p.Exec(context.Background(), query, currentTime, options.Id)

	return err
}
