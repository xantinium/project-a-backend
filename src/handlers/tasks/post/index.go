package tasks_handler_get

import (
	validation "github.com/go-ozzo/ozzo-validation"
	api_tasks_post "github.com/xantinium/project-a-backend/api/tasks/post"
	"github.com/xantinium/project-a-backend/src/core"
	core_database "github.com/xantinium/project-a-backend/src/core/database"
	core_resources "github.com/xantinium/project-a-backend/src/core/resources"
)

type tasksPostHandlerSchema struct {
	Name        string
	Description string
}

func (s *tasksPostHandlerSchema) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Name, validation.Required, validation.Length(0, 50)),
		validation.Field(&s.Description, validation.Length(0, 1024)),
	)
}

type tasksPostHandler struct {
	dbClient core.DatabaseClient
}

func (h *tasksPostHandler) Prepare(rawData []byte, ctx core.HttpCtx) *api_tasks_post.TasksPostReq {
	return api_tasks_post.GetRootAsTasksPostReq(rawData, 0)
}

func (h *tasksPostHandler) Validate(data *api_tasks_post.TasksPostReq, ctx core.HttpCtx) error {
	r := tasksPostHandlerSchema{
		Name:        string(data.Name()),
		Description: string(data.Description()),
	}

	return r.Validate()
}

func (h *tasksPostHandler) Handle(data *api_tasks_post.TasksPostReq, ctx core.HttpCtx) (core.HttpStatus, []byte, error) {
	userId := core.ExtractUserId(ctx)

	elements := core_resources.DeserializeElements(data.Elements, data.ElementsLength())

	err := h.dbClient.CreateTask(&core_database.CreateTaskOptions{
		Name:        string(data.Name()),
		Description: core_database.CreateField(string(data.Description())),
		IsPrivate:   data.IsPrivate(),
		Elements:    core_resources.SerializeElements(elements),
		OwnerId:     userId,
	})
	if err != nil {
		return core.HttpStatusInternalServerError, nil, err
	}

	return core.HttpStatusOK, nil, nil
}

func NewTasksPostHandler(dbClient core.DatabaseClient) core.Handler {
	h := &tasksPostHandler{dbClient}

	return core.NewPrivateHttpHandler(h)
}
