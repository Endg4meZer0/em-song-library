package delivery

import (
	"effective-mobile-song-library/internal/model"
	"effective-mobile-song-library/pkg/validator"
)

func ValidateSongFilters(v *validator.Validator, f model.SongFilters) {
	if f.ReleaseDate != "" {
		v.Check(
			v.Matches(f.ReleaseDate, validator.ReleaseDateRX) ||
				v.Matches(f.ReleaseDate, validator.ReleaseYearMonthRX) ||
				v.Matches(f.ReleaseDate, validator.ReleaseYearRX),
			"release_date", "invalid format of release date filter",
		)
	}

	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
}

func ValidateSongTextFilters(v *validator.Validator, f model.SongTextFilters) {
	v.Check(f.Verse > 0, "verse", "must be greater than zero")
	v.Check(f.Verse <= 10_000_000, "verse", "must be a maximum of 10 million")
}

func ValidateSongInput(v *validator.Validator, group string, song string) {
	v.Check(group != "", "group", "must be provided")
	v.Check(song != "", "song", "must be provided")
}

func ValidateSongInfo(v *validator.Validator, song *model.SongInfo) {
	v.Check(song.Group != "", "group", "must be provided")
	v.Check(song.Song != "", "song", "must be provided")

	v.Check(song.ReleaseDate != "", "release_date", "must be provided")
	v.Check(v.Matches(song.ReleaseDate, validator.ReleaseDateRX), "release_date", "invalid format of release date")

	v.Check(len(song.Text) <= 1_048_576, "text", "must not be more than 1MB long")

	v.Check(len(song.Link) <= 500, "link", "must not be more than 500 bytes long")
}
