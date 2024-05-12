package handlers

import (
	"github.com/xantinium/project-a-backend/src/core"
	core_files "github.com/xantinium/project-a-backend/src/core/files"
	auth_handler "github.com/xantinium/project-a-backend/src/handlers/auth"
	users_handler "github.com/xantinium/project-a-backend/src/handlers/users"
)

type resource interface {
	RegisterHandlers(router core.Router)
}

func RegisterHandlers(router core.Router, dbClient core.DatabaseClient) {
	resources := []resource{
		auth_handler.NewAuthResource(dbClient),
		users_handler.NewAuthResource(dbClient),
	}

	for _, r := range resources {
		r.RegisterHandlers(router)
	}

	router.GET("/images/:id", func(ctx core.HttpCtx) {
		imageId := ctx.Param("id")

		image, err := core_files.GetImage(imageId)
		if err != nil {
			ctx.Status(core.HttpStatusNotFound)
			return
		}

		ctx.Data(core.HttpStatusOK, core.BINARY_MIME_TYPE, image)
	})
}
