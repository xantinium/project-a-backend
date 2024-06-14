package core_database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type TaskStruct struct {
	Id          int       `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	IsPrivate   bool      `db:"is_private"`
	Elements    []byte    `db:"elements"`
	OwnerId     int       `db:"owner_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type GetTasksOptions struct {
	Id   *int    `db:"id"`
	Name *string `db:"name"`
}

func (dbClient *DatabaseClient) GetTasks(options *GetTasksOptions) (*[]TaskStruct, error) {
	queryStr := "SELECT * FROM tasks"
	q := CreateColumnsQuery(options)
	if q != "" {
		queryStr += fmt.Sprintf(" WHERE %s", q)
	}

	query, err := dbClient.p.Query(context.Background(), queryStr)
	if err != nil {
		return nil, err
	}

	tasks, err := pgx.CollectRows(query, pgx.RowToStructByName[TaskStruct])
	if err != nil {
		return nil, err
	}

	return &tasks, nil
}

type CreateTaskOptions struct {
	Name        string
	Description *string
	IsPrivate   bool
	Elements    []byte
	OwnerId     int
}

func (dbClient *DatabaseClient) CreateTask(options *CreateTaskOptions) error {
	currentTime := time.Now().UTC()

	_, err := dbClient.p.Exec(
		context.Background(),
		"INSERT INTO tasks VALUES (DEFAULT, $1, $2, $3, $4, $5, $6, $7)",
		options.Name,
		options.Description,
		options.IsPrivate,
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
	IsPrivate   *bool   `db:"is_private"`
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
