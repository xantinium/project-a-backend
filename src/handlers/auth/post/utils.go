package auth_handler_post

import (
	"net/http"

	"github.com/xantinium/project-a-backend/src/core"
	core_database "github.com/xantinium/project-a-backend/src/core/database"
)

func updateUserAvatar(dbClient core.DatabaseClient, userId int, avatarUri string) {
	createdImgId, err := dbClient.CreateImageFromURL(&core_database.CreateImageFromURLOptions{
		Url:     avatarUri,
		OwnerId: userId,
	})

	if err == nil {
		dbClient.UpdateUser(&core_database.UpdateUserOptions{
			Id: userId,
			Fields: &core_database.UpdateUserOptionsFields{
				AvatarId: core_database.CreateField(&createdImgId),
			},
		})
	}
}

func setAuthCookie(ctx core.HttpCtx, userId int) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:  core.AuthCookieName,
		Value: core.CreateToken(userId),
	})
}
