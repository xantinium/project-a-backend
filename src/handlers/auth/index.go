package auth_handler

import (
	"github.com/xantinium/project-a-backend/src/core"
	auth_handler_post "github.com/xantinium/project-a-backend/src/handlers/auth/post"
)

type authResource struct {
	dbClient core.DatabaseClient
}

func (res *authResource) RegisterHandlers(r core.Router) {
	r.POST("/auth", auth_handler_post.NewAutPosthHandler(res.dbClient))
}

func NewAuthResource(dbClient core.DatabaseClient) *authResource {
	return &authResource{dbClient}
}
