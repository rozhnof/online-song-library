package tests

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hardfinhq/go-date"
)

type MockMusicService struct {
	storage []Song
}

func NewMockMusicService() *MockMusicService {
	return &MockMusicService{}
}

func (s *MockMusicService) AddSong(song Song) {
	s.storage = append(s.storage, song)
}

func (s *MockMusicService) ClearStorage() {
	s.storage = nil
}

func (s *MockMusicService) Run() {
	r := gin.Default()

	type songInfo struct {
		ReleaseDate date.Date `json:"releaseDate"`
		Text        string    `json:"text"`
		Link        string    `json:"link"`
	}

	r.GET("/info", func(c *gin.Context) {
		group := c.DefaultQuery("group", "")
		song := c.DefaultQuery("song", "")

		if song == "" || group == "" {
			c.String(http.StatusBadRequest, "Missing required parameters")
			return
		}

		for _, s := range s.storage {
			if s.Group == group && s.Name == song {
				c.JSON(http.StatusOK, songInfo{
					ReleaseDate: s.ReleaseDate,
					Text:        s.Text,
					Link:        s.Link,
				})

				return
			}
		}

		c.String(http.StatusNotFound, "Song not found")
	})

	r.Run(":9091")
}
