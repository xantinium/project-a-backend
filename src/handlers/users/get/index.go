package users_handler_get

import (
	flatbuffers "github.com/google/flatbuffers/go"
	api_users_get "github.com/xantinium/project-a-backend/api/users/get"
	"github.com/xantinium/project-a-backend/src/core"
	core_database "github.com/xantinium/project-a-backend/src/core/database"
)

type usersGetHandler struct {
	dbClient core.DatabaseClient
}

func (h *usersGetHandler) Prepare(rawData []byte, ctx core.HttpCtx) *api_users_get.UsersGetReq {
	return nil
}

func (h *usersGetHandler) Validate(data *api_users_get.UsersGetReq, ctx core.HttpCtx) error {
	return nil
}

func (h *usersGetHandler) Handle(data *api_users_get.UsersGetReq, ctx core.HttpCtx) (core.HttpStatus, []byte, error) {
	options := core_database.GetUsersOptions{}

	userId := core.ExtractIntParam(ctx, "id")
	if userId != 0 {
		options.Id = core_database.CreateField(userId)
	}

	users, err := h.dbClient.GetUsers(&options)
	if err != nil {
		return core.HttpStatusInternalServerError, nil, err
	}

	b := &flatbuffers.Builder{}
	var offsets = make([]flatbuffers.UOffsetT, 0, len(*users))

	for _, user := range *users {
		firstName := b.CreateString(user.FirstName)
		lastName := b.CreateString(user.LastName)

		var avatarId flatbuffers.UOffsetT
		if user.AvatarId != nil {
			avatarId = b.CreateString(*user.AvatarId)
		}

		api_users_get.UserStart(b)
		api_users_get.UserAddId(b, uint32(user.Id))
		api_users_get.UserAddFirstName(b, firstName)
		api_users_get.UserAddLastName(b, lastName)
		api_users_get.UserAddAvatarId(b, avatarId)

		offset := api_users_get.UserEnd(b)
		offsets = append(offsets, offset)
	}

	offset := b.CreateVectorOfTables(offsets)

	api_users_get.UsersGetResStart(b)
	api_users_get.UsersGetResAddUsers(b, offset)
	offset = api_users_get.UsersGetResEnd(b)
	b.Finish(offset)

	return core.HttpStatusOK, b.FinishedBytes(), nil
}

func NewUsersGetHandler(dbClient core.DatabaseClient) core.Handler {
	h := &usersGetHandler{dbClient}

	return core.NewPrivateHttpHandler(h)
}
