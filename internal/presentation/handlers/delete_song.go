package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	repo "song-service/internal/application/repository"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeleteSongResponse struct {
	DeletedTime time.Time `json:"deleted_time"`
}

// DeleteSong godoc
// @Summary      Delete song by ID
// @Description  Удаление песни из библиотеки по ID
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        id     path     string  true  "Song ID"
// @Success      200    {object} DeleteSongResponse
// @Failure      400    {string} string  "Invalid ID format"
// @Failure      404    {string} string  "Song not found"
// @Failure      500    {string} string  "Internal Server Error"
// @Router       /songs/{id} [delete]
func (h *SongHandler) DeleteSong(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "SongHandler.DeleteSong")
	defer span.End()

	id, err := uuid.Parse(c.Param(pathParamID))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID format")
		return
	}

	deletedTime, err := h.songService.DeleteSong(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrObjectNotFound) {
			c.String(http.StatusNotFound, err.Error())
			return
		}

		h.logger.Warn("failed to create song", slog.String("error", err.Error()))

		c.Status(http.StatusInternalServerError)
		return
	}

	response := DeleteSongResponse{
		DeletedTime: *deletedTime,
	}

	c.JSON(http.StatusOK, response)
}
