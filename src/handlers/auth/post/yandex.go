package auth_handler_post

import (
	"encoding/json"
	"fmt"
	"net/http"

	api_auth "github.com/xantinium/project-a-backend/api/auth"
	"github.com/xantinium/project-a-backend/src/core"
	core_database "github.com/xantinium/project-a-backend/src/core/database"
)

type yandexProfileInfo struct {
	Id            string `json:"id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	IsAvatarEmpty bool   `json:"is_avatar_empty"`
	AvatarId      string `json:"default_avatar_id"`
}

func createNewYandexUser(dbClient core.DatabaseClient, info *yandexProfileInfo) (int, error) {
	createdUserId, err := dbClient.CreateUser(&core_database.CreateUserOptions{
		FirstName:       info.FirstName,
		LastName:        info.LastName,
		AvatarId:        nil,
		OAuthSerive:     api_auth.OAuthServicesYANDEX,
		YandexProfileId: &info.Id,
		GoogleProfileId: nil,
	})
	if err != nil {
		return -1, err
	}

	if !info.IsAvatarEmpty {
		avatarUri := fmt.Sprintf("https://avatars.yandex.net/get-yapic/%s/islands-200", info.AvatarId)
		go updateUserAvatar(dbClient, createdUserId, avatarUri)
	}

	return createdUserId, nil
}

func handleYandex(dbClient core.DatabaseClient, accessToken string, ctx core.HttpCtx) (core.HttpStatus, []byte, error) {
	req, err := http.NewRequest(http.MethodGet, "https://login.yandex.ru/info", nil)
	if err != nil {
		return core.HttpStatusBadRequest, nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", string(accessToken)))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return core.HttpStatusBadRequest, nil, err
	}

	defer res.Body.Close()

	var info yandexProfileInfo
	err = json.NewDecoder(res.Body).Decode(&info)
	if err != nil {
		return core.HttpStatusBadRequest, nil, err
	}

	user, err := dbClient.GetUserByServiceProfileId(&core_database.GetUserByServiceProfileIdOptions{
		Service:          api_auth.OAuthServicesYANDEX,
		ServiceProfileId: info.Id,
	})

	if err != nil {
		userId, err := createNewYandexUser(dbClient, &info)
		if err != nil {
			return core.HttpStatusBadRequest, nil, err
		}

		setAuthCookie(ctx, userId)

		return core.HttpStatusOK, nil, nil
	}

	setAuthCookie(ctx, user.Id)

	return core.HttpStatusOK, nil, nil
}
