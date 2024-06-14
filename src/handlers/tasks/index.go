package tasks_handler

import (
	"github.com/xantinium/project-a-backend/src/core"
	tasks_handler_get "github.com/xantinium/project-a-backend/src/handlers/tasks/get"
	tasks_handler_patch "github.com/xantinium/project-a-backend/src/handlers/tasks/patch"
	tasks_handler_post "github.com/xantinium/project-a-backend/src/handlers/tasks/post"
)

type tasksResource struct {
	dbClient core.DatabaseClient
}

func (res *tasksResource) RegisterHandlers(r core.Router) {
	r.GET("/tasks", tasks_handler_get.NewTasksGetHandler(res.dbClient))
	r.GET("/tasks/:id", tasks_handler_get.NewTasksGetHandler(res.dbClient))
	r.POST("/tasks", tasks_handler_post.NewTasksPostHandler(res.dbClient))
	r.PATCH("/tasks/:id", tasks_handler_patch.NewTasksPatchHandler(res.dbClient))
}

func NewTasksResource(dbClient core.DatabaseClient) *tasksResource {
	return &tasksResource{dbClient}
}
