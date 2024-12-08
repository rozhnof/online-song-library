package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	repo "song-service/internal/application/repository"
	"song-service/internal/domain/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SongQueryParams struct {
	repo.Pagination
}

type SongResponse struct {
	Song models.Song `json:"song"`
}

// Song godoc
// @Summary      Get song
// @Description  Получение песни с пагинацией по куплетам
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        id       path     string  true   "Song ID"
// @Param        limit    query    int     false  "Limit number of verses"
// @Param        offset   query    int     false  "Offset for pagination"
// @Success      200      {object} SongResponse
// @Failure      400      {string} string  "Invalid ID format"
// @Failure      404      {string} string  "Song not found"
// @Failure      500      {string} string  "Internal Server Error"
// @Router       /songs/{id} [get]
func (h *SongHandler) Song(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "SongHandler.Song")
	defer span.End()

	id, err := uuid.Parse(c.Param(pathParamID))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID format")
		return
	}

	var queryParams SongQueryParams
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		h.logger.Debug("failed to parse query parameters", slog.String("error", err.Error()))

		c.String(http.StatusBadRequest, err.Error())
		return
	}

	pagination := repo.Pagination{
		Limit:  queryParams.Limit,
		Offset: queryParams.Offset,
	}

	song, err := h.songService.Song(ctx, id, &pagination)
	if err != nil {
		if errors.Is(err, repo.ErrObjectNotFound) {
			c.String(http.StatusNotFound, err.Error())
			return
		}

		h.logger.Warn("failed to create song", slog.String("error", err.Error()))

		c.Status(http.StatusInternalServerError)
		return
	}

	response := SongResponse{
		Song: song,
	}

	c.JSON(http.StatusOK, response)
}
