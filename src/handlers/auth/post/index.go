package auth_handler_post

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
	api_auth "github.com/xantinium/project-a-backend/api/auth"
	api_auth_post "github.com/xantinium/project-a-backend/api/auth/post"
	"github.com/xantinium/project-a-backend/src/core"
)

type authPostHandlerSchema struct {
	Service     api_auth.OAuthServices
	AccessToken string
}

func isOAuthService(value interface{}) error {
	service, ok := value.(api_auth.OAuthServices)

	if !ok {
		return errors.New("cannot be blank")
	}

	if _, ok := api_auth.EnumNamesOAuthServices[service]; !ok {
		return errors.New("unknown service")
	}

	return nil
}

func (s *authPostHandlerSchema) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Service, validation.By(isOAuthService)),
		validation.Field(&s.AccessToken, validation.Required),
	)
}

type authPostHandler struct {
	dbClient core.DatabaseClient
}

func (h *authPostHandler) Prepare(rawData []byte, ctx core.HttpCtx) *api_auth_post.AuthPostReq {
	return api_auth_post.GetRootAsAuthPostReq(rawData, 0)
}

func (h *authPostHandler) Validate(data *api_auth_post.AuthPostReq, ctx core.HttpCtx) error {
	r := authPostHandlerSchema{
		Service:     data.Service(),
		AccessToken: string(data.AccessToken()),
	}

	return r.Validate()
}

func (h *authPostHandler) Handle(data *api_auth_post.AuthPostReq, ctx core.HttpCtx) (core.HttpStatus, []byte, error) {
	accessToken := string(data.AccessToken())

	switch data.Service() {
	case api_auth.OAuthServicesYANDEX:
		return handleYandex(h.dbClient, accessToken, ctx)
	case api_auth.OAuthServicesGOOGLE:
		return handleGoogle(h.dbClient, accessToken, ctx)
	}

	return core.HttpStatusInternalServerError, nil, nil
}

func NewAutPosthHandler(dbClient core.DatabaseClient) core.Handler {
	h := &authPostHandler{dbClient}

	return core.NewPublicHttpHandler(h)
}
