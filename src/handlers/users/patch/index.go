package users_handler_patch

import (
	"errors"

	api_users_patch "github.com/xantinium/project-a-backend/api/users/patch"
	"github.com/xantinium/project-a-backend/src/core"
	core_database "github.com/xantinium/project-a-backend/src/core/database"
)

type usersPatchHandler struct {
	dbClient core.DatabaseClient
}

func (h *usersPatchHandler) Prepare(rawData []byte, ctx core.HttpCtx) *api_users_patch.UsersPatchReq {
	return api_users_patch.GetRootAsUsersPatchReq(rawData, 0)
}

func (h *usersPatchHandler) Validate(data *api_users_patch.UsersPatchReq, ctx core.HttpCtx) error {
	userId := core.ExtractIntParam(ctx, "id")

	if userId == 0 {
		return errors.New("invalid user id")
	}

	return nil
}

func (h *usersPatchHandler) Handle(data *api_users_patch.UsersPatchReq, ctx core.HttpCtx) (core.HttpStatus, []byte, error) {
	userId := core.ExtractIntParam(ctx, "id")

	options := core_database.UpdateUserOptions{
		Id: userId,
		Fields: &core_database.UpdateUserOptionsFields{
			FirstName: core_database.CreateField(string(data.FirstName())),
			LastName:  core_database.CreateField(string(data.LastName())),
		},
	}

	if data.AvatarBytesLength() != 0 {
		avatarBytes := make([]byte, 0, data.AvatarBytesLength())

		for i := 0; i < data.AvatarBytesLength(); i++ {
			avatarBytes = append(avatarBytes, byte(data.AvatarBytes(i)))
		}

		imgId, err := h.dbClient.CreateImage(&core_database.CreateImageOptions{
			Data:    avatarBytes,
			OwnerId: userId,
		})
		if err != nil {
			return core.HttpStatusInternalServerError, nil, err
		}

		options.Fields.AvatarId = core_database.CreateField(&imgId)
	}

	err := h.dbClient.UpdateUser(&options)
	if err != nil {
		return core.HttpStatusInternalServerError, nil, err
	}

	return core.HttpStatusOK, nil, nil
}

func NewUsersPatchHandler(dbClient core.DatabaseClient) core.Handler {
	h := &usersPatchHandler{dbClient}

	return core.NewPrivateHttpHandler(h)
}
