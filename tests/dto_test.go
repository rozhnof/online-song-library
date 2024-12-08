package tests

import (
	"time"

	"github.com/google/uuid"
	"github.com/hardfinhq/go-date"
)

type Song struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"song"`
	Group       string    `json:"group"`
	ReleaseDate date.Date `json:"release_date"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

type CreateSongRequest struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

type CreateSongResponse struct {
	Song Song `json:"song"`
}

type DeleteSongResponse struct {
	DeletedTime time.Time `json:"deleted_time"`
}

type UpdateSongRequest struct {
	Name        string     `json:"song,omitempty"`
	Group       string     `json:"group,omitempty"`
	ReleaseDate *date.Date `json:"release_date,omitempty"`
	Text        string     `json:"text,omitempty"`
	Link        string     `json:"link,omitempty"`
}

type UpdateSongResponse struct {
	Song Song `json:"song"`
}

type SongQueryParams struct {
	Group  string `form:"group"`
	Song   string `form:"song"`
	Limit  int32  `form:"limit"`
	Offset int32  `form:"offset"`
}

type SongResponse struct {
	Song Song `json:"song"`
}

type SongListQueryParams struct {
	Name            []string   `form:"song"`
	Group           []string   `form:"group"`
	ReleaseDateFrom *date.Date `form:"release_date_from"`
	ReleaseDateTo   *date.Date `form:"release_date_to"`
	Text            []string   `form:"text"`
	Link            []string   `form:"link"`
	Limit           int32      `form:"limit"`
	Offset          int32      `form:"offset"`
}

type ListSongResponse struct {
	SongList []Song `json:"song_list"`
}
