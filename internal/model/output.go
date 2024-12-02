package model

type ErrRes struct {
	Error any `json:"error"`
}

type SongOut struct {
	ID          uint64 `json:"id"`
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Link        string `json:"link"`
	TotalVerses uint   `json:"totalVerses"`
}

type Songs struct {
	Songs []SongOut `json:"songs"`
}

type SongText struct {
	Text string `json:"text"`
}
