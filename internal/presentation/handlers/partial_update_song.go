package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	repo "song-service/internal/application/repository"
	"song-service/internal/domain/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hardfinhq/go-date"
)

type PartialUpdateSongRequest struct {
	Name        string    `json:"song"`
	Group       string    `json:"group"`
	ReleaseDate date.Date `json:"release_date" swaggertype:"primitive,string"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

type PartialUpdateSongResponse struct {
	Song models.Song `json:"song"`
}

// PartialUpdateSong godoc
// @Summary      Partially update song by ID
// @Description  Частичное обновление информации о песне в библиотеке
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        id       path     string                     true  "Song ID"
// @Param        request  body     PartialUpdateSongRequest   true  "Song details to be updated"
// @Success      200      {object} PartialUpdateSongResponse
// @Failure      400      {string} string                    "Invalid input data"
// @Failure      404      {string} string                    "Song not found"
// @Failure      409      {string} string                    "Name conflict"
// @Failure      500      {string} string                    "Internal Server Error"
// @Router       /songs/{id} [patch]
func (h *SongHandler) PartialUpdateSong(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "SongHandler.PartialUpdateSong")
	defer span.End()

	id, err := uuid.Parse(c.Param(pathParamID))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID format")
		return
	}

	var request PartialUpdateSongRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Debug("failed to parse request body", slog.String("error", err.Error()))

		c.String(http.StatusBadRequest, err.Error())
		return
	}

	song := models.Song{
		ID:          id,
		Name:        request.Name,
		Group:       request.Group,
		ReleaseDate: request.ReleaseDate,
		Text:        request.Text,
		Link:        request.Link,
	}

	updatedSong, err := h.songService.UpdateSong(ctx, song)
	if err != nil {
		if errors.Is(err, repo.ErrObjectNotFound) {
			c.String(http.StatusNotFound, err.Error())
			return
		}

		if errors.Is(err, repo.ErrDuplicate) {
			c.String(http.StatusConflict, err.Error())
			return
		}

		h.logger.Warn("failed to create song", slog.String("error", err.Error()))

		c.Status(http.StatusInternalServerError)
		return
	}

	response := SongResponse{
		Song: updatedSong,
	}

	c.JSON(http.StatusOK, response)
}
