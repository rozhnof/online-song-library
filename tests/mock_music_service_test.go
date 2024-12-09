package tests

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hardfinhq/go-date"
)

const (
	mockMusicServerAddress = "localhost:9091"
)

type songInfo struct {
	ReleaseDate date.Date `json:"releaseDate"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

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

	r.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.Run(mockMusicServerAddress)
}

func (s *MockMusicService) waitForServer(timeout time.Duration, retryInterval time.Duration) error {
	url := fmt.Sprintf("http://%s/ping", mockMusicServerAddress)

	for start := time.Now(); time.Since(start) < timeout; time.Sleep(retryInterval) {
		if resp, err := http.Get(url); err == nil {
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}
	}

	return fmt.Errorf("mock music server not available")
}
