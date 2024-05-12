package core_database

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	core_files "github.com/xantinium/project-a-backend/src/core/files"
)

type CreateImageOptions struct {
	Data     []byte
	FileName *string
	OwnerId  int
}

func (dbClient *DatabaseClient) CreateImage(options *CreateImageOptions) (string, error) {
	currentTime := time.Now().UTC()
	imgId := uuid.New().String()
	imgName := imgId
	if options.FileName != nil {
		imgName = *options.FileName
	}

	_, err := dbClient.p.Exec(
		context.Background(),
		"INSERT INTO images VALUES ($1, $2, $3, $4, $5)",
		imgId,
		imgName,
		options.OwnerId,
		currentTime,
		currentTime,
	)
	if err != nil {
		return "", err
	}

	return imgId, core_files.SaveImage(options.Data, imgName)
}

type CreateImageFromURLOptions struct {
	Url     string
	OwnerId int
}

func (dbClient *DatabaseClient) CreateImageFromURL(options *CreateImageFromURLOptions) (string, error) {
	res, err := http.Get(options.Url)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	imgData, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return dbClient.CreateImage(&CreateImageOptions{
		Data:     imgData,
		FileName: nil,
		OwnerId:  options.OwnerId,
	})
}

func (dbClient *DatabaseClient) DeleteImage(id string) error {
	_, err := dbClient.p.Exec(
		context.Background(),
		"DELETE FROM images WHERE id = $1",
		id,
	)

	return err
}
