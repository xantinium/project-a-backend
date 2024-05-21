package tasks_handler_patch

import (
	validation "github.com/go-ozzo/ozzo-validation"
	api_tasks_patch "github.com/xantinium/project-a-backend/api/tasks/patch"
	"github.com/xantinium/project-a-backend/src/core"
	core_database "github.com/xantinium/project-a-backend/src/core/database"
	core_resources "github.com/xantinium/project-a-backend/src/core/resources"
)

type tasksPatchHandlerSchema struct {
	Name        string
	Description string
}

func (s *tasksPatchHandlerSchema) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Name, validation.Required, validation.Length(0, 50)),
		validation.Field(&s.Description, validation.Length(0, 1024)),
	)
}

type tasksPatchHandler struct {
	dbClient core.DatabaseClient
}

func (h *tasksPatchHandler) Prepare(rawData []byte, ctx core.HttpCtx) *api_tasks_patch.TasksPatchReq {
	return api_tasks_patch.GetRootAsTasksPatchReq(rawData, 0)
}

func (h *tasksPatchHandler) Validate(data *api_tasks_patch.TasksPatchReq, ctx core.HttpCtx) error {
	r := tasksPatchHandlerSchema{
		Name:        string(data.Name()),
		Description: string(data.Description()),
	}

	return r.Validate()
}

func (h *tasksPatchHandler) Handle(data *api_tasks_patch.TasksPatchReq, ctx core.HttpCtx) (core.HttpStatus, []byte, error) {
	userId := core.ExtractUserId(ctx)

	elements := core_resources.DeserializeElements(data.Elements, data.ElementsLength())

	err := h.dbClient.UpdateTask(&core_database.UpdateTaskOptions{
		Id: userId,
		Fields: &core_database.UpdateTaskOptionsFields{
			Name:        core_database.CreateField(string(data.Name())),
			Description: core_database.CreateField(string(data.Description())),
			Elements:    core_database.CreateField(core_resources.SerializeElements(elements)),
		},
	})
	if err != nil {
		return core.HttpStatusInternalServerError, nil, err
	}

	return core.HttpStatusOK, nil, nil
}

func NewTasksPatchHandler(dbClient core.DatabaseClient) core.Handler {
	h := &tasksPatchHandler{dbClient}

	return core.NewPrivateHttpHandler(h)
}
