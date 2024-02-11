package handlers

import (
	"github.com/dbubel/intake"
)

func Endpoints(visionAPI *App) intake.Endpoints {
	endpoints := intake.Endpoints{
		intake.POST("/api/v1/workspaces", visionAPI.CreateWorkspace),
		intake.POST("/api/v1/workspaces/repos", visionAPI.CreateRepo),
		intake.POST("/api/v1/workspaces/vectors", visionAPI.InsertVector),

		intake.GET("/api/v1/search", visionAPI.Search),
		intake.GET("/health", visionAPI.Health),
	}

	return endpoints
}
