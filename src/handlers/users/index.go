package users_handler

import (
	"github.com/xantinium/project-a-backend/src/core"
	users_handler_get "github.com/xantinium/project-a-backend/src/handlers/users/get"
	users_handler_patch "github.com/xantinium/project-a-backend/src/handlers/users/patch"
)

type usersResource struct {
	dbClient core.DatabaseClient
}

func (res *usersResource) RegisterHandlers(r core.Router) {
	r.GET("/users", users_handler_get.NewUsersGetHandler(res.dbClient))
	r.GET("/users/:id", users_handler_get.NewUsersGetHandler(res.dbClient))
	r.PATCH("/users/:id", users_handler_patch.NewUsersPatchHandler(res.dbClient))
}

func NewAuthResource(dbClient core.DatabaseClient) *usersResource {
	return &usersResource{dbClient}
}
