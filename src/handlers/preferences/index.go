package preferences_handler

import (
	"github.com/xantinium/project-a-backend/src/core"
	preferences_handler_get "github.com/xantinium/project-a-backend/src/handlers/preferences/get"
	preferences_handler_patch "github.com/xantinium/project-a-backend/src/handlers/preferences/patch"
)

type preferencesResource struct {
	dbClient core.DatabaseClient
}

func (res *preferencesResource) RegisterHandlers(r core.Router) {
	r.GET("/preferences", preferences_handler_get.NewPreferencesGetHandler(res.dbClient))
	r.PATCH("/preferences", preferences_handler_patch.NewPreferencesPatchHandler(res.dbClient))
}

func NewPreferencesResource(dbClient core.DatabaseClient) *preferencesResource {
	return &preferencesResource{dbClient}
}
