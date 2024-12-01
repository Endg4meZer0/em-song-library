package http

import (
	"errors"
	"net/http"

	"effective-mobile-song-library/internal/delivery"
	"effective-mobile-song-library/internal/model"
	"effective-mobile-song-library/internal/repository/db"
	errResponses "effective-mobile-song-library/pkg/errors"
	"effective-mobile-song-library/pkg/jsonutil"
	"effective-mobile-song-library/pkg/logger"
	"effective-mobile-song-library/pkg/validator"
)

type SongLibraryService interface {
	Get(id uint64) (*model.SongInfo, error)
	GetAll(model.SongFilters) ([]*model.SongInfo, error)
	GetText(model.SongTextFilters) (*string, error)
	InsertSongs(groups []string, songs []string) ([]*model.SongInfo, error)
	Update(cars *model.SongInfo) error
	Delete(id uint64) error
}

// @Summary list
// @Tags songs
// @Description listing songs data
// @Accept json
// @Produce json
// @Param  group   query string  false  "name search by group"
// @Param  song   query string  false  "name search by song"
// @Param  page   query uint  false  "page number, default 1"
// @Param  pageSize   query uint  false  "page size, default 10"
// @Success 200 {object} model.Songs
// @Failure 422 {object} model.ErrRes
// @Failure 500 {object} model.ErrRes
// @Router       /songs [get]
func (h *Handler) listSongsHandler(w http.ResponseWriter, r *http.Request) {
	var filters model.SongFilters
	qs := r.URL.Query()
	v := validator.New()

	filters.Group = readString(qs, "group", "")
	filters.Song = readString(qs, "song", "")

	filters.Page = readUint(qs, "page", 1, v)
	filters.PageSize = readUint(qs, "page_size", 10, v)

	if delivery.ValidateSongFilters(v, filters); !v.Valid() {
		errResponses.FailedValidationResponse(w, r, v.Errors)
		return
	}

	logger.PrintDebug("", map[string]any{
		"method":  r.Method,
		"url":     r.URL.String(),
		"filters": filters,
	})

	songs, err := h.service.GetAll(filters)
	if err != nil {
		errResponses.ServerErrorResponse(w, r, err)
		return
	}

	logger.PrintDebug("", map[string]any{
		"url":               r.URL.String(),
		"number of records": len(songs),
		"songs list":        songs,
	})
	// Send a JSON response containing the song info.
	err = jsonutil.WriteJSON(w, http.StatusOK, songs, nil)
	if err != nil {
		errResponses.ServerErrorResponse(w, r, err)
	}
}

// @Summary get text
// @Tags songs
// @Description get song's text
// @Accept json
// @Produce json
// @Param  id path uint true "song id"
// @Param  verse   query uint  false  "verse number, default 1"
// @Success 200 {object} model.SongText
// @Failure 400 {object} model.ErrRes
// @Failure 422 {object} model.ErrRes
// @Failure 500 {object} model.ErrRes
// @Router       /songs/{id}/text [get]
func (h *Handler) listSongTextHandler(w http.ResponseWriter, r *http.Request) {
	var filters model.SongTextFilters
	qs := r.URL.Query()
	v := validator.New()

	id := readIDFromPath(r, v)
	if !v.Valid() {
		errResponses.NotFoundResponse(w, r)
		return
	}

	filters.ID = uint64(id)
	filters.Verse = readUint(qs, "verse", 1, v)

	if delivery.ValidateSongTextFilters(v, filters); !v.Valid() {
		errResponses.FailedValidationResponse(w, r, v.Errors)
		return
	}

	logger.PrintDebug("", map[string]any{
		"method":  r.Method,
		"url":     r.URL.String(),
		"filters": filters,
	})

	verse, err := h.service.GetText(filters)
	if err != nil {
		errResponses.ServerErrorResponse(w, r, err)
		return
	}

	logger.PrintDebug("", map[string]any{
		"url":   r.URL.String(),
		"verse": verse,
	})

	// Send a JSON response containing the verse.
	err = jsonutil.WriteJSON(w, http.StatusOK, verse, nil)
	if err != nil {
		errResponses.ServerErrorResponse(w, r, err)
	}
}

// @Summary add
// @Tags songs
// @Description add songs info
// @Accept json
// @Produce json
// @Param  songs body model.SongsInput  true  "song collection"
// @Success 200 {object} model.Songs
// @Failure 400 {object} model.ErrRes
// @Failure 422 {object} model.ErrRes
// @Failure 500 {object} model.ErrRes
// @Router       /songs [post]
func (h *Handler) addSongInfoHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Groups []string `json:"groups"`
		Songs  []string `json:"songs"`
	}

	err := jsonutil.ReadJSON(w, r, &input)
	if err != nil {
		errResponses.BadRequestResponse(w, r, err)
		return
	}

	logger.PrintDebug("", map[string]any{
		"method":       r.Method,
		"url":          r.URL.String(),
		"input groups": input.Groups,
		"input songs":  input.Songs,
	})

	// validate
	v := validator.New()
	if delivery.ValidateSongs(v, input.Groups, input.Songs); !v.Valid() {
		errResponses.FailedValidationResponse(w, r, v.Errors)
		return
	}

	song, err := h.service.InsertSongs(input.Groups, input.Songs)
	if err != nil {
		errResponses.ServerErrorResponse(w, r, err)
		return
	}

	err = jsonutil.WriteJSON(w, http.StatusOK, song, nil)
	if err != nil {
		errResponses.ServerErrorResponse(w, r, err)
	}
}

// @Summary update
// @Tags songs
// @Description update song data by ID
// @Accept json
// @Produce json
// @Param  id   path    uint  true  "song ID"
// @Param  input body   model.SongInput   true  "song info struct"
// @Success 200 {object} model.SongInfo
// @Failure 400 {object} model.ErrRes
// @Failure 404 {object} model.ErrRes
// @Failure 422 {object} model.ErrRes
// @Failure 500 {object} model.ErrRes
// @Router       /songs/{id} [patch]
func (h *Handler) updateSongInfoHandler(w http.ResponseWriter, r *http.Request) {
	v := validator.New()
	id := readIDFromPath(r, v)
	if !v.Valid() {
		errResponses.NotFoundResponse(w, r)
		return
	}

	// Fetch the existing song info from the database
	song, err := h.service.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, db.ErrRecordNotFound):
			errResponses.NotFoundResponse(w, r)
		default:
			errResponses.ServerErrorResponse(w, r, err)
		}
		return
	}

	// Declare an input struct to hold the expected data from the client.
	var input model.SongInput

	logger.PrintDebug("", map[string]any{
		"method": r.Method,
		"url":    r.URL.String(),
		"id":     id,
		"input":  input,
	})

	err = jsonutil.ReadJSON(w, r, &input)
	if err != nil {
		errResponses.BadRequestResponse(w, r, err)
		return
	}
	// Copy the values from the request body
	if input.Group != nil {
		song.Group = *input.Group
	}
	if input.Song != nil {
		song.Song = *input.Song
	}
	if input.ReleaseDate != nil {
		song.ReleaseDate = *input.ReleaseDate
	}
	if input.Text != nil {
		song.Text = *input.Text
	}
	if input.Link != nil {
		song.Link = *input.Link
	}

	// validate
	v = validator.New()
	if delivery.ValidateSongInfo(v, song); !v.Valid() {
		errResponses.FailedValidationResponse(w, r, v.Errors)
		return
	}

	err = h.service.Update(song)
	if err != nil {
		errResponses.ServerErrorResponse(w, r, err)
		return
	}

	logger.PrintDebug("updated", map[string]any{
		"song": song,
	})

	err = jsonutil.WriteJSON(w, http.StatusOK, song, nil)
	if err != nil {
		errResponses.ServerErrorResponse(w, r, err)
	}
}

// @Summary delete
// @Tags songs
// @Description delete song data
// @Accept json
// @Produce json
// @Param  id   path      uint  true  "song ID"
// @Success 200 {object} model.SongInfo
// @Failure 404 {object} model.ErrRes
// @Failure 500 {object} model.ErrRes
// @Router       /songs/{id} [delete]
func (h *Handler) deleteSongInfoHandler(w http.ResponseWriter, r *http.Request) {
	v := validator.New()
	id := readIDFromPath(r, v)
	if !v.Valid() {
		errResponses.NotFoundResponse(w, r)
		return
	}

	logger.PrintDebug("", map[string]any{
		"method": r.Method,
		"url":    r.URL.String(),
		"id":     id,
	})

	err := h.service.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, db.ErrRecordNotFound):
			errResponses.NotFoundResponse(w, r)
		default:
			errResponses.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = jsonutil.WriteJSON(w, http.StatusOK, map[string]string{"message": "song info successfully deleted"}, nil)
	if err != nil {
		errResponses.ServerErrorResponse(w, r, err)
	}
}
