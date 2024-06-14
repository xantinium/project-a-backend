package preferences_handler_patch

import (
	api_preferences_patch "github.com/xantinium/project-a-backend/api/preferences/patch"
	"github.com/xantinium/project-a-backend/src/core"
	core_database "github.com/xantinium/project-a-backend/src/core/database"
)

type preferencesPatchHandler struct {
	dbClient core.DatabaseClient
}

func (h *preferencesPatchHandler) Prepare(rawData []byte, ctx core.HttpCtx) *api_preferences_patch.PreferencesPatchRes {
	return api_preferences_patch.GetRootAsPreferencesPatchRes(rawData, 0)
}

func (h *preferencesPatchHandler) Validate(data *api_preferences_patch.PreferencesPatchRes, ctx core.HttpCtx) error {
	return nil
}

func (h *preferencesPatchHandler) Handle(data *api_preferences_patch.PreferencesPatchRes, ctx core.HttpCtx) (core.HttpStatus, []byte, error) {
	userId := core.ExtractUserId(ctx)

	err := h.dbClient.UpdateUser(&core_database.UpdateUserOptions{
		Id: userId,
		Fields: &core_database.UpdateUserOptionsFields{
			Theme: core_database.CreateField(data.Theme()),
			Lang:  core_database.CreateField(data.Lang()),
		},
	})
	if err != nil {
		return core.HttpStatusInternalServerError, nil, err
	}

	return core.HttpStatusOK, nil, nil
}

func NewPreferencesPatchHandler(dbClient core.DatabaseClient) core.Handler {
	h := &preferencesPatchHandler{dbClient}

	return core.NewPrivateHttpHandler(h)
}
