package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/xantinium/project-a-backend/src/core"
	core_database "github.com/xantinium/project-a-backend/src/core/database"
	"github.com/xantinium/project-a-backend/src/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("no .env file")
	}

	if os.Getenv("MODE") == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	dbClient, err := core_database.NewDatabaseClient()

	defer dbClient.Dispose()

	if err != nil {
		panic(err)
	}

	core.RegisterStaticResolver(router)

	router.Use(gin.RecoveryWithWriter(nil, func(ctx *gin.Context, err any) {
		fmt.Println(err)
		ctx.AbortWithStatus(500)
	}))

	handlers.RegisterHandlers(router.Group("/api"), dbClient)

	router.Run(fmt.Sprintf("%s:80", os.Getenv("PLATFORM_IP")))
}
