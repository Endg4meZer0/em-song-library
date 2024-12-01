package model

type SongInput struct {
	Group       *string   `json:"group"`
	Song        *string   `json:"song"`
	ReleaseDate *string   `json:"releaseDate"`
	Link        *string   `json:"link"`
	Text        *[]string `json:"text"`
}

type SongsInput struct {
	Groups []string `json:"groups"`
	Songs  []string `json:"songs"`
}
