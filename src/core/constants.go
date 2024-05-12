package core

import (
	"github.com/gin-gonic/gin"
	core_database "github.com/xantinium/project-a-backend/src/core/database"
)

const (
	BINARY_MIME_TYPE = "application/octet-stream"
	HTML_MIME_TYPE   = "text/html"
)

const (
	HttpStatusOK                  HttpStatus = 200
	HttpStatusBadRequest          HttpStatus = 400
	HttpStatusNotFound            HttpStatus = 404
	HttpStatusUnauthorized        HttpStatus = 401
	HttpStatusInternalServerError HttpStatus = 500
)

type (
	Router         = *gin.RouterGroup
	HttpCtx        = *gin.Context
	Handler        = func(HttpCtx)
	HttpStatus     = int
	DatabaseClient = *core_database.DatabaseClient
)
