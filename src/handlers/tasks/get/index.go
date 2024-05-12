package tasks_handler_get

import (
	api_users_get "github.com/xantinium/project-a-backend/api/users/get"
	"github.com/xantinium/project-a-backend/src/core"
)

type tasksGetHandler struct {
	dbClient core.DatabaseClient
}

func (h *tasksGetHandler) Prepare(rawData []byte, ctx core.HttpCtx) *api_users_get.UsersGetReq {
	return nil
}

func (h *tasksGetHandler) Validate(data *api_users_get.UsersGetReq, ctx core.HttpCtx) error {
	return nil
}

func (h *tasksGetHandler) Handle(data *api_users_get.UsersGetReq, ctx core.HttpCtx) (core.HttpStatus, []byte, error) {
	return core.HttpStatusOK, nil, nil
}

func NewTasksGetHandler(dbClient core.DatabaseClient) core.Handler {
	h := &tasksGetHandler{dbClient}

	return core.NewPrivateHttpHandler(h)
}
