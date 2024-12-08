package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/hardfinhq/go-date"
)

type SongInfo struct {
	ReleaseDate date.Date `json:"releaseDate"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

type MusicServiceClient struct {
	client  *http.Client
	baseURL string
}

func NewMusicServiceClient(client *http.Client, baseURL string) *MusicServiceClient {
	return &MusicServiceClient{
		client:  client,
		baseURL: baseURL,
	}
}

func (c MusicServiceClient) Info(ctx context.Context, group string, song string) (*SongInfo, int, error) {
	url := fmt.Sprintf("%s/info?group=%s&song=%s", c.baseURL, group, song)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, resp.StatusCode, err
		}

		return nil, resp.StatusCode, errors.New(string(body))
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	var songInfo SongInfo
	if err := json.Unmarshal(bytes, &songInfo); err != nil {
		return nil, 0, err
	}

	return &songInfo, resp.StatusCode, nil
}
