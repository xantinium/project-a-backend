package core

import (
	"net/http"
	"strconv"

	flatbuffers "github.com/google/flatbuffers/go"
	httpError "github.com/xantinium/project-a-backend/api/error"
)

const AuthCookieName = "X-Auth-Token"

type httpHandler[T any] interface {
	Prepare([]byte, HttpCtx) *T
	Validate(*T, HttpCtx) error
	Handle(*T, HttpCtx) (statusCode HttpStatus, resultData []byte, err error)
}

func newHttpError(e error) []byte {
	b := &flatbuffers.Builder{}
	errMsg := b.CreateString(e.Error())
	httpError.ErrorResStart(b)
	httpError.ErrorResAddMessage(b, errMsg)
	b.Finish(httpError.ErrorResEnd(b))
	return b.FinishedBytes()
}

func NewPublicHttpHandler[T any](h httpHandler[T]) Handler {
	handler := func(ctx HttpCtx) {
		// Получаем бинарные данные
		rawData, err := ctx.GetRawData()
		if err != nil {
			ctx.Data(HttpStatusInternalServerError, BINARY_MIME_TYPE, nil)
			return
		}
		// Преобразовываем бинарные данные в объект конкретного типа
		data := h.Prepare(rawData, ctx)
		// Валидируем полученный объект
		err = h.Validate(data, ctx)
		if err != nil {
			ctx.Data(HttpStatusBadRequest, BINARY_MIME_TYPE, newHttpError(err))
			return
		}
		// Обрабатываем запрос, генерируем ответ
		statusCode, resultData, err := h.Handle(data, ctx)
		if err != nil {
			ctx.Data(statusCode, BINARY_MIME_TYPE, newHttpError(err))
		}
		ctx.Data(statusCode, BINARY_MIME_TYPE, resultData)
	}

	return handler
}

func NewPrivateHttpHandler[T any](h httpHandler[T]) Handler {
	publicHandler := NewPublicHttpHandler[T](h)

	handler := func(ctx HttpCtx) {
		cookie, err := ctx.Request.Cookie(AuthCookieName)
		if err != nil {
			ctx.Data(HttpStatusUnauthorized, BINARY_MIME_TYPE, newHttpError(ErrInvalidToken))
			return
		}

		tokenPayload, err := decode(cookie.Value)
		if err != nil {
			ctx.Data(HttpStatusUnauthorized, BINARY_MIME_TYPE, newHttpError(err))
			return
		}

		http.SetCookie(ctx.Writer, &http.Cookie{
			Name:  AuthCookieName,
			Value: CreateToken(tokenPayload.UserId),
		})
		ctx.Set(AuthCookieName, tokenPayload)

		publicHandler(ctx)
	}

	return handler
}

func ExtractIntParam(ctx HttpCtx, param string) int {
	value := ctx.Param(param)

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return parsedValue
}
