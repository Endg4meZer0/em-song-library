package service

import (
	"errors"
	"reflect"

	"effective-mobile-song-library/internal/model"
	"effective-mobile-song-library/internal/repository/external"
	"effective-mobile-song-library/pkg/logger"
)

type (
	SongStorage interface {
		Get(id uint64) (*model.SongInfo, error)
		GetAll(filters model.SongFilters) ([]*model.SongInfo, error)
		GetFullText(id uint64) (*string, error)
		GetText(filters model.SongTextFilters) (*string, error)
		Insert(*model.SongInfo) error
		Update(songs *model.SongInfo) error
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

func (sl *SongLibraryService) GetAll(filters model.SongFilters) ([]*model.SongOut, error) {
	songs, err := sl.songRepo.GetAll(filters)
	if err != nil {
		return nil, err
	}

	songOuts := make([]*model.SongOut, 0, len(songs))
	for _, song := range songs {
		songOuts = append(songOuts, &model.SongOut{
			ID:          song.ID,
			Group:       song.Group,
			Song:        song.Song,
			ReleaseDate: song.ReleaseDate,
			Link:        song.Link,
			TotalVerses: uint(len(song.Text)),
		})
	}
	return songOuts, nil
}

func (sl *SongLibraryService) GetText(filters model.SongTextFilters) (*string, error) {
	if filters.Verse == 0 {
		return sl.songRepo.GetFullText(filters.ID)
	} else {
		return sl.songRepo.GetText(filters)
	}
}

func (sl *SongLibraryService) Insert(group string, song string) (*model.SongInfo, error) {
	songInfo, err := sl.apiClient.GetSongInfoWithDetails(group, song)

	logger.PrintDebug("info from external API", map[string]any{
		"songInfo": songInfo,
	})

	if err != nil {
		if errors.Is(err, external.ErrBadRequest) {
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
