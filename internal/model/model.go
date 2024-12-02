package model

type SongInfo struct {
	ID          uint64   `json:"id"`
	Group       string   `json:"group"`
	Song        string   `json:"song"`
	ReleaseDate string   `json:"releaseDate"`
	Text        []string `json:"text"`
	Link        string   `json:"link"`
}
