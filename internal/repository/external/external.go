package external

import (
	"effective-mobile-song-library/config"
	"effective-mobile-song-library/internal/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

var (
	ErrBadRequest = errors.New("incorrect request")
	ErrNoResponse = errors.New("no response from API")
)

type ApiClient struct {
	config *config.Config
}

func NewApiClient(config *config.Config) *ApiClient {
	return &ApiClient{config: config}
}

func (ac *ApiClient) GetSongInfoWithDetails(group string, song string) (*model.SongInfo, error) {
	resp, err := http.Get(fmt.Sprintf("%s/info?group=%s&song=%s", ac.config.ExternalAPIURL, url.PathEscape(group), url.PathEscape(song)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return nil, ErrBadRequest
		} else {
			return nil, ErrNoResponse
		}
	}

	var songInfoWithDetailsDTO SongInfoWithDetailsDTO = SongInfoWithDetailsDTO{Group: group, Song: song}
	if err := json.NewDecoder(resp.Body).Decode(&songInfoWithDetailsDTO); err != nil {
		return nil, err
	}

	return songInfoWithDetailsDTO.ToModel(), nil
}
