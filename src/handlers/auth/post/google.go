package auth_handler_post

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/xantinium/project-a-backend/api/auth"
	"github.com/xantinium/project-a-backend/src/core"
	core_database "github.com/xantinium/project-a-backend/src/core/database"
)

type googleProfileInfo struct {
	Id        string `json:"sub"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	AvatarUri string `json:"picture"`
}

func createNewGoogleUser(dbClient core.DatabaseClient, info *googleProfileInfo) (int, error) {
	createdUserId, err := dbClient.CreateUser(&core_database.CreateUserOptions{
		FirstName:       info.FirstName,
		LastName:        info.LastName,
		AvatarId:        nil,
		OAuthSerive:     auth.OAuthServicesGOOGLE,
		YandexProfileId: nil,
		GoogleProfileId: &info.Id,
	})
	if err != nil {
		return -1, err
	}

	go updateUserAvatar(dbClient, createdUserId, info.AvatarUri)

	return createdUserId, nil
}

func handleGoogle(dbClient core.DatabaseClient, accessToken string, ctx core.HttpCtx) (core.HttpStatus, []byte, error) {
	req, err := http.NewRequest(http.MethodGet, "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return core.HttpStatusBadRequest, nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", string(accessToken)))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return core.HttpStatusBadRequest, nil, err
	}

	defer res.Body.Close()

	var info googleProfileInfo
	err = json.NewDecoder(res.Body).Decode(&info)
	if err != nil {
		return core.HttpStatusBadRequest, nil, err
	}

	user, err := dbClient.GetUserByServiceProfileId(&core_database.GetUserByServiceProfileIdOptions{
		Service:          auth.OAuthServicesYANDEX,
		ServiceProfileId: info.Id,
	})

	if err != nil {
		userId, err := createNewGoogleUser(dbClient, &info)
		if err != nil {
			return core.HttpStatusBadRequest, nil, err
		}

		setAuthCookie(ctx, userId)

		return core.HttpStatusOK, nil, nil
	}

	setAuthCookie(ctx, user.Id)

	return core.HttpStatusOK, nil, nil
}
