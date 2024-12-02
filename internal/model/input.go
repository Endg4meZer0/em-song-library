package model

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
