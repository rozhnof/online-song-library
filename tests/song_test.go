package tests

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	date "github.com/hardfinhq/go-date"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateNonExistentSong(t *testing.T) {
	if err := SetUpEmpty(); err != nil {
		t.Fatal(err)
	}

	req := CreateSongRequest{
		Group: nonExistentSong.Group,
		Song:  nonExistentSong.Name,
	}

	var (
		expectedStatusCode = http.StatusNotFound
	)

	_, code, err := songServiceClient.CreateSong(req, nil)

	require.NotNil(t, err)
	assert.Equal(t, expectedStatusCode, code)
}

func TestCreate(t *testing.T) {
	if err := SetUpCreateTest(); err != nil {
		t.Fatal(err)
	}

	var (
		expectedSong       = defaultSong
		expectedStatusCode = http.StatusOK
	)

	req := CreateSongRequest{
		Group: expectedSong.Group,
		Song:  expectedSong.Name,
	}

	t.Run("add song", func(t *testing.T) {
		resp, code, err := songServiceClient.CreateSong(req, nil)

		require.Nil(t, err)
		require.NotNil(t, resp)

		expectedSong.ID = resp.Song.ID

		assert.Equal(t, expectedStatusCode, code)
		assert.Equal(t, expectedSong, resp.Song)
	})

	t.Run("second add song", func(t *testing.T) {
		resp, code, err := songServiceClient.CreateSong(req, nil)

		require.Nil(t, err)
		require.NotNil(t, resp)

		assert.Equal(t, expectedStatusCode, code)
		assert.Equal(t, expectedSong, resp.Song)
	})
}

func TestUpdateNonExistentSong(t *testing.T) {
	if err := SetUpEmpty(); err != nil {
		t.Fatal(err)
	}

	req := UpdateSongRequest{
		Name:        nonExistentSong.Name,
		Group:       nonExistentSong.Group,
		ReleaseDate: &nonExistentSong.ReleaseDate,
		Text:        nonExistentSong.Text,
		Link:        nonExistentSong.Link,
	}

	var (
		expectedStatusCode = http.StatusNotFound
	)

	_, code, err := songServiceClient.UpdateSong(nonExistentSong.ID, req, nil)

	require.NotNil(t, err)
	assert.Equal(t, expectedStatusCode, code)
}

func TestPartialUpdateNonExistentSong(t *testing.T) {
	if err := SetUpEmpty(); err != nil {
		t.Fatal(err)
	}

	req := UpdateSongRequest{
		Name:  nonExistentSong.Name,
		Group: nonExistentSong.Group,
	}

	var (
		expectedStatusCode = http.StatusNotFound
	)

	_, code, err := songServiceClient.PartialUpdateSong(nonExistentSong.ID, req, nil)

	require.NotNil(t, err)
	assert.Equal(t, expectedStatusCode, code)
}

func TestUpdate(t *testing.T) {
	if err := SetUpDefault(); err != nil {
		t.Fatal(err)
	}

	expectedStatusCode := http.StatusOK

	expectedSong := Song{
		ID:          defaultSong.ID,
		Name:        "new-song-name",
		Group:       "new-song-group",
		Text:        "new-song-text",
		Link:        "new-song-link",
		ReleaseDate: date.NewDate(2026, 1, 1),
	}

	req := UpdateSongRequest{
		Name:        expectedSong.Name,
		Group:       expectedSong.Group,
		ReleaseDate: &expectedSong.ReleaseDate,
		Text:        expectedSong.Text,
		Link:        expectedSong.Link,
	}

	resp, code, err := songServiceClient.UpdateSong(expectedSong.ID, req, nil)

	require.Nil(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, expectedStatusCode, code)
	assert.Equal(t, expectedSong, resp.Song)
}

func TestPartialUpdate(t *testing.T) {
	if err := SetUpDefault(); err != nil {
		t.Fatal(err)
	}

	expectedStatusCode := http.StatusOK

	expectedSong := Song{
		ID:          defaultSong.ID,
		Name:        defaultSong.Name,
		Group:       "new-song-group",
		Text:        "new-song-text",
		Link:        defaultSong.Link,
		ReleaseDate: defaultSong.ReleaseDate,
	}

	req := UpdateSongRequest{
		Group: "new-song-group",
		Text:  "new-song-text",
	}

	resp, code, err := songServiceClient.PartialUpdateSong(expectedSong.ID, req, nil)

	require.Nil(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, expectedStatusCode, code)
	assert.Equal(t, expectedSong, resp.Song)
}

func TestGetNonExistentSong(t *testing.T) {
	if err := SetUpEmpty(); err != nil {
		t.Fatal(err)
	}

	queryParams := SongQueryParams{
		Group: nonExistentSong.Group,
		Song:  nonExistentSong.Name,
	}

	var (
		expectedStatusCode = http.StatusNotFound
	)

	_, code, err := songServiceClient.GetSong(nonExistentSong.ID, queryParams)

	require.NotNil(t, err)
	assert.Equal(t, expectedStatusCode, code)
}

func TestGet(t *testing.T) {
	if err := SetUpDefault(); err != nil {
		t.Fatal(err)
	}

	var (
		expectedStatusCode = http.StatusOK
		expectedSong       = defaultSong
	)

	resp, code, err := songServiceClient.GetSong(defaultSong.ID, nil)

	require.Nil(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, expectedStatusCode, code)
	assert.Equal(t, expectedSong, resp.Song)
}

func TestGetTextPagination(t *testing.T) {
	song := Song{
		ID:          uuid.New(),
		Group:       "song-group",
		Name:        "song-name",
		ReleaseDate: date.NewDate(2025, 1, 1),
		Text:        "verse-1\n\nverse-2\n\nverse-3\n\nverse-4\n\nverse-5\n\nverse-6",
		Link:        "song-link",
	}

	var (
		musicServiceSongs = []Song{song}
		songServiceSongs  = []Song{song}
	)

	if err := SetUp(musicServiceSongs, songServiceSongs); err != nil {
		t.Fatal(err)
	}

	expectedVersesByLimitAndOffset := func(limit int32, offset int32) string {
		const sep = "\n\n"

		verseList := strings.Split(song.Text, sep)

		if offset < 0 {
			offset = 0
		} else if offset >= int32(len(verseList)) {
			return ""
		}

		if limit <= 0 || offset+limit >= int32(len(verseList)) {
			return strings.Join(verseList[offset:], sep)
		}

		return strings.Join(verseList[offset:offset+limit], sep)
	}

	limits := []int32{-1, 0, 1, 3, 6, 7}
	offsets := []int32{-1, 0, 1, 3, 6, 7}

	for _, limit := range limits {
		for _, offset := range offsets {
			t.Run(fmt.Sprintf("limit %d offset %d", limit, offset), func(t *testing.T) {
				queryParams := SongQueryParams{
					Group:  song.Group,
					Song:   song.Name,
					Limit:  limit,
					Offset: offset,
				}

				var (
					expectedStatusCode = http.StatusOK
					expectedText       = expectedVersesByLimitAndOffset(limit, offset)
				)

				resp, code, err := songServiceClient.GetSong(song.ID, queryParams)

				require.NoError(t, err)
				require.NotNil(t, resp)

				assert.Equal(t, expectedStatusCode, code)
				assert.Equal(t, expectedText, resp.Song.Text)
			})
		}
	}
}

func TestDeleteNonExistentSong(t *testing.T) {
	if err := SetUpEmpty(); err != nil {
		t.Fatal(err)
	}

	var (
		expectedStatusCode = http.StatusNotFound
	)

	_, code, err := songServiceClient.DeleteSong(nonExistentSong.ID, nil)

	require.NotNil(t, err)
	assert.Equal(t, expectedStatusCode, code)
}

func TestDelete(t *testing.T) {
	if err := SetUpDefault(); err != nil {
		t.Fatal(err)
	}

	var (
		expectedStatusCode = http.StatusOK
	)

	var deletedTime time.Time

	t.Run("delete", func(t *testing.T) {
		resp, code, err := songServiceClient.DeleteSong(defaultSong.ID, nil)

		require.Nil(t, err)
		require.NotNil(t, resp)

		assert.Equal(t, expectedStatusCode, code)

		deletedTime = resp.DeletedTime
	})

	t.Run("second delete", func(t *testing.T) {
		resp, code, err := songServiceClient.DeleteSong(defaultSong.ID, nil)

		require.Nil(t, err)
		require.NotNil(t, resp)

		assert.Equal(t, expectedStatusCode, code)
		assert.Equal(t, deletedTime, resp.DeletedTime)
	})
}

func TestListEmpty(t *testing.T) {
	if err := SetUpEmpty(); err != nil {
		t.Fatal(err)
	}

	resp, code, err := songServiceClient.ListSong(nil)

	var (
		expectedStatusCode          = http.StatusOK
		expectedResponseSongListLen = 0
	)

	require.Nil(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, expectedStatusCode, code)
	assert.Equal(t, expectedResponseSongListLen, len(resp.SongList))
}

func TestList(t *testing.T) {
	var songList []Song

	for i := 0; i < 3; i++ {
		for j := 0; j < 3-i; j++ {
			var (
				songID          = uuid.New()
				songGroup       = fmt.Sprintf("group-%d", i)
				songName        = fmt.Sprintf("name-%d", j)
				songReleaseDate = date.NewDate(2024, 1, 1)
				songText        = "song-text"
				songLink        = "song-link"
			)

			song := Song{
				ID:          songID,
				Name:        songName,
				Group:       songGroup,
				ReleaseDate: songReleaseDate,
				Text:        songText,
				Link:        songLink,
			}

			songList = append(songList, song)
		}
	}

	var (
		musicServiceSongs = songList
		songServiceSongs  = songList
	)

	if err := SetUp(musicServiceSongs, songServiceSongs); err != nil {
		t.Fatal(err)
	}

	sortAndAssertSongList := func(s1 []Song, s2 []Song) {
		assert.Equal(t, len(s1), len(s2))

		sort.Slice(s1, func(i, j int) bool {
			return s1[i].ID.String() < s1[j].ID.String()
		})

		sort.Slice(s2, func(i, j int) bool {
			return s2[i].ID.String() < s2[j].ID.String()
		})

		for i := 0; i < len(s1); i++ {
			assert.Equal(t, s1[i], s2[i])
		}
	}

	var (
		expectedStatusCode = http.StatusOK
	)

	t.Run("list all", func(t *testing.T) {
		var (
			expectedSongList = songList
		)

		resp, code, err := songServiceClient.ListSong(nil)

		require.Nil(t, err)
		require.NotNil(t, resp)

		assert.Equal(t, expectedStatusCode, code)
		sortAndAssertSongList(expectedSongList, resp.SongList)
	})

	t.Run("list with group filter", func(t *testing.T) {
		queryParams := SongListQueryParams{
			Group: []string{"group-1"},
		}

		resp, code, err := songServiceClient.ListSong(queryParams)

		require.Nil(t, err)
		require.NotNil(t, resp)

		assert.Equal(t, expectedStatusCode, code)
		for _, song := range resp.SongList {
			assert.Contains(t, queryParams.Group, song.Group)
		}
	})

	t.Run("list with name filter", func(t *testing.T) {
		queryParams := SongListQueryParams{
			Name: []string{"name-1"},
		}

		resp, code, err := songServiceClient.ListSong(queryParams)

		require.Nil(t, err)
		require.NotNil(t, resp)

		assert.Equal(t, expectedStatusCode, code)
		for _, song := range resp.SongList {
			assert.Contains(t, queryParams.Name, song.Name)
		}
	})

	t.Run("list with limit", func(t *testing.T) {
		queryParams := SongListQueryParams{
			Limit: 3,
		}

		resp, code, err := songServiceClient.ListSong(queryParams)

		require.Nil(t, err)
		require.NotNil(t, resp)

		assert.Equal(t, expectedStatusCode, code)
		assert.Equal(t, queryParams.Limit, int32(len(resp.SongList)))
	})

	t.Run("list with offset", func(t *testing.T) {
		queryParams := SongListQueryParams{
			Offset: 3,
		}

		resp, code, err := songServiceClient.ListSong(queryParams)

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, expectedStatusCode, code)

		assert.Equal(t, len(songList)-int(queryParams.Offset), len(resp.SongList))
	})
}
