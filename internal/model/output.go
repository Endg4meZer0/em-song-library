package model

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