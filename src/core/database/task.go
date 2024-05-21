package core_database

import (
	"context"
	"fmt"
	"time"
)

type TaskStruct struct {
	Id          int       `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Elements    []byte    `db:"elements"`
	OwnerId     int       `db:"owner_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type CreateTaskOptions struct {
	Name        string
	Description *string
	Elements    []byte
	OwnerId     int
}

func (dbClient *DatabaseClient) CreateTask(options *CreateTaskOptions) error {
	currentTime := time.Now().UTC()

	_, err := dbClient.p.Exec(
		context.Background(),
		"INSERT INTO tasks VALUES (DEFAULT, $1, $2, $3, $4, $5, $6)",
		options.Name,
		options.Description,
		options.Elements,
		options.OwnerId,
		currentTime,
		currentTime,
	)

	return err
}

type UpdateTaskOptionsFields struct {
	Name        *string `db:"name"`
	Description *string `db:"description"`
	Elements    *[]byte `db:"elements"`
}

type UpdateTaskOptions struct {
	Id     int
	Fields *UpdateTaskOptionsFields
}

func (dbClient *DatabaseClient) UpdateTask(options *UpdateTaskOptions) error {
	currentTime := time.Now().UTC()

	query := fmt.Sprintf("UPDATE tasks SET updated_at = $1, %s WHERE id = $2", CreateColumnsQuery(options.Fields))

	_, err := dbClient.p.Exec(context.Background(), query, currentTime, options.Id)

	return err
}
