package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-playground/form/v4"
	"github.com/google/uuid"
)

type SongServiceClient struct {
	client  *http.Client
	baseURL string
}

func NewSongClient(client *http.Client, baseURL string) *SongServiceClient {
	return &SongServiceClient{
		client:  client,
		baseURL: baseURL,
	}
}

func (c *SongServiceClient) CreateSong(request CreateSongRequest, queryParams any) (*CreateSongResponse, int, error) {
	return makeRequest[CreateSongRequest, CreateSongResponse](c.client, c.baseURL, "/songs", http.MethodPost, &request, queryParams)
}

func (c *SongServiceClient) UpdateSong(id uuid.UUID, request UpdateSongRequest, queryParams any) (*UpdateSongResponse, int, error) {
	return makeRequest[UpdateSongRequest, UpdateSongResponse](c.client, c.baseURL, fmt.Sprintf("/songs/%s", id.String()), http.MethodPut, &request, queryParams)
}

func (c *SongServiceClient) PartialUpdateSong(id uuid.UUID, request UpdateSongRequest, queryParams any) (*UpdateSongResponse, int, error) {
	return makeRequest[UpdateSongRequest, UpdateSongResponse](c.client, c.baseURL, fmt.Sprintf("/songs/%s", id.String()), http.MethodPatch, &request, queryParams)
}

func (c *SongServiceClient) DeleteSong(id uuid.UUID, queryParams any) (*DeleteSongResponse, int, error) {
	return makeRequest[struct{}, DeleteSongResponse](c.client, c.baseURL, fmt.Sprintf("/songs/%s", id.String()), http.MethodDelete, nil, queryParams)
}

func (c *SongServiceClient) GetSong(id uuid.UUID, queryParams any) (*SongResponse, int, error) {
	return makeRequest[struct{}, SongResponse](c.client, c.baseURL, fmt.Sprintf("/songs/%s", id.String()), http.MethodGet, nil, queryParams)
}

func (c *SongServiceClient) ListSong(queryParams any) (*ListSongResponse, int, error) {
	return makeRequest[struct{}, ListSongResponse](c.client, c.baseURL, "/songs", http.MethodGet, nil, queryParams)
}

func makeRequest[Req any, Resp any](client *http.Client, baseURL string, endpoint string, method string, request *Req, queryParams any) (*Resp, int, error) {
	url, err := buildURL(baseURL, endpoint, queryParams)
	if err != nil {
		return nil, 0, err
	}

	var body io.Reader

	if request != nil {
		requestBytes, err := json.Marshal(request)
		if err != nil {
			return nil, 0, err
		}

		body = bytes.NewBuffer(requestBytes)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, resp.StatusCode, err
		}

		return nil, resp.StatusCode, errors.New(string(body))
	}

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	var response Resp

	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return nil, 0, err
	}

	return &response, resp.StatusCode, nil
}

func buildURL(baseURL string, endpoint string, queryParams any) (string, error) {
	fullPath := fmt.Sprintf("%s%s", baseURL, endpoint)

	if queryParams == nil {
		return fullPath, nil
	}

	parsedURL, err := url.Parse(fullPath)
	if err != nil {
		return "", err
	}

	encoder := form.NewEncoder()
	values, err := encoder.Encode(queryParams)
	if err != nil {
		return "", err
	}

	parsedURL.RawQuery = values.Encode()

	return parsedURL.String(), nil
}
