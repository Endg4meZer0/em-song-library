package service

import (
	"errors"
	"reflect"

	"effective-mobile-song-library/internal/model"
	"effective-mobile-song-library/internal/repository/api"
	"effective-mobile-song-library/pkg/logger"
)

type (
	SongStorage interface {
		Get(id uint64) (*model.SongInfo, error)
		GetAll(filters model.SongFilters) ([]*model.SongInfo, error)
		GetText(filters model.SongTextFilters) (*string, error)
		Insert(*model.SongInfo) error
		Update(cars *model.SongInfo) error
		Delete(id uint64) error
	}

	ApiClient interface {
		GetSongInfoWithDetails(group string, song string) (*model.SongInfo, error)
	}
)

type SongLibraryService struct {
	songRepo  SongStorage
	apiClient ApiClient
}

func NewSongLibraryService(songRepo SongStorage, apiClient ApiClient) *SongLibraryService {
	return &SongLibraryService{
		songRepo:  songRepo,
		apiClient: apiClient,
	}
}

func (sl *SongLibraryService) Get(id uint64) (*model.SongInfo, error) {
	return sl.songRepo.Get(id)
}

func (sl *SongLibraryService) GetAll(filters model.SongFilters) ([]*model.SongInfo, error) {
	return sl.songRepo.GetAll(filters)
}

func (sl *SongLibraryService) GetText(filters model.SongTextFilters) (*string, error) {
	return sl.songRepo.GetText(filters)
}

func (sl *SongLibraryService) Insert(group string, song string) (*model.SongInfo, error) {
	songInfo, err := sl.apiClient.GetSongInfoWithDetails(group, song)

	logger.PrintDebug("info from external API", map[string]any{
		"songInfo": songInfo,
	})

	if err != nil {
		if errors.Is(err, api.ErrBadRequest) {
			logger.PrintDebug("did not add song", map[string]any{
				"group": group,
				"song":  song,
				"error": err,
			})
		} else {
			return nil, err
		}
	}

	if songInfo != nil {
		err = sl.songRepo.Insert(songInfo)
		if err != nil {
			return nil, err
		}
	}

	return songInfo, nil
}

func (sl *SongLibraryService) Update(song *model.SongInfo) error {
	if !reflect.DeepEqual(*song, model.SongInfo{}) {
		err := sl.songRepo.Update(song)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sl *SongLibraryService) Delete(id uint64) error {
	return sl.songRepo.Delete(id)
}
