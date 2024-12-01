package http

type Handler struct {
	service SongLibraryService
}

func NewHandler(service SongLibraryService) *Handler {
	return &Handler{service: service}
}
