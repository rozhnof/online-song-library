package repo

import "github.com/hardfinhq/go-date"

type Pagination struct {
	Limit  int32 `form:"limit"`
	Offset int32 `form:"offset"`
}

type SongFilter struct {
	Name            []string   `form:"song"`
	Group           []string   `form:"group"`
	ReleaseDateFrom *date.Date `form:"release_date_from"`
	ReleaseDateTo   *date.Date `form:"release_date_to"`
	Text            []string   `form:"text"`
	Link            []string   `form:"link"`
}
