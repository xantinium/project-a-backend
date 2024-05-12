package tasks_handler_get

import (
	validation "github.com/go-ozzo/ozzo-validation"
	api_tasks_post "github.com/xantinium/project-a-backend/api/tasks/post"
	"github.com/xantinium/project-a-backend/src/core"
)

type tasksPostHandlerSchema struct {
	Id          int
	Name        string
	Description string
}

func (s *tasksPostHandlerSchema) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Id, validation.Required, validation.Min(0)),
		validation.Field(&s.Name, validation.Required, validation.Max(50)),
		validation.Field(&s.Description, validation.Max(1024)),
	)
}

type tasksPostHandler struct {
	dbClient core.DatabaseClient
}

func (h *tasksPostHandler) Prepare(rawData []byte, ctx core.HttpCtx) *api_tasks_post.TasksPostReq {
	return api_tasks_post.GetRootAsTasksPostReq(rawData, 0)
}

func (h *tasksPostHandler) Validate(data *api_tasks_post.TasksPostReq, ctx core.HttpCtx) error {
	taskId := core.ExtractIntParam(ctx, "id")

	r := tasksPostHandlerSchema{
		Id:          taskId,
		Name:        string(data.Name()),
		Description: string(data.Description()),
	}

	return r.Validate()
}

func (h *tasksPostHandler) Handle(data *api_tasks_post.TasksPostReq, ctx core.HttpCtx) (core.HttpStatus, []byte, error) {
	// taskId := core.ExtractIntParam(ctx, "id")
	return core.HttpStatusOK, nil, nil
}

func NewTasksPostHandler(dbClient core.DatabaseClient) core.Handler {
	h := &tasksPostHandler{dbClient}

	return core.NewPrivateHttpHandler(h)
}
