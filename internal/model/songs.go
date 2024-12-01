package model

type SongInfo struct {
	ID          uint64   `json:"id"`
	Group       string   `json:"group"`
	Song        string   `json:"song"`
	ReleaseDate string   `json:"releaseDate"`
	Text        []string `json:"text"`
	Link        string   `json:"link"`
}

type SongFilters struct {
	Group       string
	Song        string
	ReleaseDate string
	Text        string
	Link        string
	PageSize    uint
	Page        uint
}

type SongTextFilters struct {
	ID    uint64
	Verse uint
}

type ErrRes struct {
	Error any `json:"error"`
}

type Songs struct {
	Songs []SongInfo `json:"songs"`
}

type SongText struct {
	Verse       string
	TotalVerses uint
}

type SongInput struct {
	Group       *string   `json:"group"`
	Song        *string   `json:"song"`
	ReleaseDate *string   `json:"releaseDate"`
	Text        *[]string `json:"text"`
	Link        *string   `json:"link"`
}

type SongsInput struct {
	Groups []string `json:"groups"`
	Songs  []string `json:"songs"`
}
