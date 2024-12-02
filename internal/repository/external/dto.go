package external

import (
	"effective-mobile-song-library/internal/model"
	"strings"
)

type SongInfoWithDetailsDTO struct {
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Link        string `json:"link"`
	Text        string `json:"text"`
}

func (dto *SongInfoWithDetailsDTO) ToModel() *model.SongInfo {
	var songInfo model.SongInfo
	songInfo.Group = dto.Group
	songInfo.Song = dto.Song
	songInfo.ReleaseDate = dto.ReleaseDate
	songInfo.Link = dto.Link
	songInfo.Text = strings.Split(dto.Text, "\n\n")
	return &songInfo
}
