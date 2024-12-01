package http

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "effective-mobile-song-library/docs"
	responses "effective-mobile-song-library/pkg/errors"

	"github.com/julienschmidt/httprouter"
)

func (h *Handler) Routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(responses.NotFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(responses.MethodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/songs", h.listSongsHandler)
	router.HandlerFunc(http.MethodGet, "/songs/:id/text", h.listSongTextHandler)
	router.HandlerFunc(http.MethodPost, "/songs", h.addSongInfoHandler)
	router.HandlerFunc(http.MethodPatch, "/songs/:id", h.updateSongInfoHandler)
	router.HandlerFunc(http.MethodDelete, "/songs/:id", h.deleteSongInfoHandler)

	router.HandlerFunc(http.MethodGet, "/swagger/:any", httpSwagger.WrapHandler)

	return router
}
