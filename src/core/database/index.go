package core_database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseClient struct {
	p *pgxpool.Pool
}

func NewDatabaseClient() (*DatabaseClient, error) {
	pool, err := pgxpool.New(context.Background(), os.Getenv("POSTGRESQL_CONNECTION_STRING"))
	if err != nil {
		return nil, err
	}

	// Порядок имеет значение
	tables := []string{
		"users",
		"images",
	}

	for _, tableName := range tables {
		fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s.sql", os.Getenv("POSTGRESQL_TABLES_PATH"), tableName))
		if err != nil {
			pool.Close()
			return nil, err
		}

		_, err = pool.Exec(context.Background(), string(fileContent))
		if err != nil {
			pool.Close()
			return nil, err
		}
	}

	return &DatabaseClient{p: pool}, nil
}

func (dbClient *DatabaseClient) Dispose() {
	dbClient.p.Close()
}
