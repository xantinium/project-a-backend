package preferences_handler_get

import (
	"errors"

	flatbuffers "github.com/google/flatbuffers/go"
	api_preferences_get "github.com/xantinium/project-a-backend/api/preferences/get"
	"github.com/xantinium/project-a-backend/src/core"
	core_database "github.com/xantinium/project-a-backend/src/core/database"
)

type preferencesGetHandler struct {
	dbClient core.DatabaseClient
}

type requestType = *interface{}

func (h *preferencesGetHandler) Prepare(rawData []byte, ctx core.HttpCtx) requestType {
	return nil
}

func (h *preferencesGetHandler) Validate(data requestType, ctx core.HttpCtx) error {
	return nil
}

func (h *preferencesGetHandler) Handle(data requestType, ctx core.HttpCtx) (core.HttpStatus, []byte, error) {
	userId := core.ExtractUserId(ctx)

	users, err := h.dbClient.GetUsers(&core_database.GetUsersOptions{
		Id: core_database.CreateField(userId),
	})
	if err != nil {
		return core.HttpStatusInternalServerError, nil, err
	}

	if len(*users) == 0 {
		return core.HttpStatusInternalServerError, nil, errors.New("user not found")
	}

	user := (*users)[0]

	b := &flatbuffers.Builder{}

	api_preferences_get.PreferencesGetResStart(b)
	api_preferences_get.PreferencesGetResAddTheme(b, user.Theme)
	api_preferences_get.PreferencesGetResAddLang(b, user.Lang)
	offset := api_preferences_get.PreferencesGetResEnd(b)
	b.Finish(offset)

	return core.HttpStatusOK, b.FinishedBytes(), nil
}

func NewPreferencesGetHandler(dbClient core.DatabaseClient) core.Handler {
	h := &preferencesGetHandler{dbClient}

	return core.NewPrivateHttpHandler(h)
}
