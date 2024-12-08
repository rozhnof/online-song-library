package app

import (
	handlers "song-service/internal/presentation/handlers"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router gin.IRoutes, songHandler *handlers.SongHandler) {
	router.POST("/songs", songHandler.CreateSong)
	router.GET("/songs", songHandler.SongList)
	router.GET("/songs/:id", songHandler.Song)
	router.PUT("/songs/:id", songHandler.UpdateSong)
	router.PATCH("/songs/:id", songHandler.PartialUpdateSong)
	router.DELETE("/songs/:id", songHandler.DeleteSong)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
