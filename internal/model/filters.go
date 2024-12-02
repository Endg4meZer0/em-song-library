package model

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
