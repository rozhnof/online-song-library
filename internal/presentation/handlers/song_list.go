package handlers

import (
	"log/slog"
	"net/http"
	repo "song-service/internal/application/repository"

	"song-service/internal/domain/models"

	"github.com/gin-gonic/gin"
)

type SongListQueryParams struct {
	repo.SongFilter
	repo.Pagination
}

type SongListResponse struct {
	SongList []models.Song `json:"song_list"`
}

// SongList godoc
// @Summary      Get songs list
// @Description  Получение списка песен с фильтрацией по всем полям и пагинацией
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        song                query    string  false  "Name of song"
// @Param        group               query    string  false  "Group name of song"
// @Param        release_date_from   query    string  false  "Start date for release date filter" example("2020-01-01")
// @Param        release_date_to     query    string  false  "End date for release date filter"   example("2023-01-01")
// @Param        text                query    string  false  "Text content of the song"
// @Param        link                query    string  false  "URL link for the song"
// @Param        limit               query    int     false  "Limit of songs"        default(10)
// @Param        offset              query    int     false  "Offset for pagination" default(0)
// @Success      200                 {object} SongListResponse
// @Failure      400                 {string} string  "Invalid query parameters"
// @Failure      500                 {string} string  "Internal Server Error"
// @Router       /songs [get]
func (h *SongHandler) SongList(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "SongHandler.SongList")
	defer span.End()

	var queryParams SongListQueryParams
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		h.logger.Debug("failed to parse query parameters", slog.String("error", err.Error()))

		c.String(http.StatusBadRequest, err.Error())
		return
	}

	songList, err := h.songService.SongList(ctx, &queryParams.SongFilter, &queryParams.Pagination)
	if err != nil {
		h.logger.Warn("failed to create song", slog.String("error", err.Error()))

		c.Status(http.StatusInternalServerError)
		return
	}

	response := SongListResponse{
		SongList: songList,
	}

	c.JSON(http.StatusOK, response)
}
