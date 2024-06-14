package tasks_handler_get

import (
	"strconv"
	"strings"

	flatbuffers "github.com/google/flatbuffers/go"
	api_tasks_get "github.com/xantinium/project-a-backend/api/tasks/get"
	"github.com/xantinium/project-a-backend/src/core"
	core_database "github.com/xantinium/project-a-backend/src/core/database"
	core_resources "github.com/xantinium/project-a-backend/src/core/resources"
)

type tasksGetHandler struct {
	dbClient core.DatabaseClient
}

func (h *tasksGetHandler) Prepare(rawData []byte, ctx core.HttpCtx) *api_tasks_get.TasksGetReq {
	v := ctx.GetHeader("X-Request-Payload")
	chars := strings.Split(v, ",")
	bytes := make([]byte, 0, len(chars))
	for _, char := range chars {
		n, err := strconv.Atoi(char)
		if err != nil {
			return nil
		}
		bytes = append(bytes, byte(n))
	}
	return api_tasks_get.GetRootAsTasksGetReq(bytes, 0)
}

func (h *tasksGetHandler) Validate(data *api_tasks_get.TasksGetReq, ctx core.HttpCtx) error {
	return nil
}

func (h *tasksGetHandler) Handle(data *api_tasks_get.TasksGetReq, ctx core.HttpCtx) (core.HttpStatus, []byte, error) {
	options := &core_database.GetTasksOptions{}

	optionId := int(data.Id())
	if optionId != 0 {
		options.Id = core_database.CreateField(optionId)
	}

	optionName := string(data.Name())
	if optionName != "" {
		options.Name = core_database.CreateField(optionName)
	}

	tasks, err := h.dbClient.GetTasks(options)
	if err != nil {
		return core.HttpStatusInternalServerError, nil, err
	}

	b := &flatbuffers.Builder{}
	var offsets = make([]flatbuffers.UOffsetT, 0, len(*tasks))

	for _, task := range *tasks {
		offset := core_resources.SerializeTask(b, core_resources.TaskType{
			Id:          task.Id,
			Name:        task.Name,
			Description: task.Description,
			IsPrivate:   task.IsPrivate,
			Elements:    core_resources.DeserializeElementsFromBytes(task.Elements),
		})

		offsets = append(offsets, offset)
	}

	offset := b.CreateVectorOfTables(offsets)

	api_tasks_get.TasksGetResStart(b)
	api_tasks_get.TasksGetResAddTasks(b, offset)
	offset = api_tasks_get.TasksGetResEnd(b)
	b.Finish(offset)

	return core.HttpStatusOK, b.FinishedBytes(), nil
}

func NewTasksGetHandler(dbClient core.DatabaseClient) core.Handler {
	h := &tasksGetHandler{dbClient}

	return core.NewPrivateHttpHandler(h)
}
