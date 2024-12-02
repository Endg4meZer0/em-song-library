package responses

import (
	"fmt"
	"net/http"

	"effective-mobile-song-library/pkg/jsonutil"
	"effective-mobile-song-library/pkg/logger"
)

func LogError(r *http.Request, err error) {
	logger.PrintError(err, map[string]any{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	err := jsonutil.WriteJSON(w, status, message, nil)
	if err != nil {
		LogError(r, err)
		w.WriteHeader(500)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	LogError(r, err)
	message := "the server encountered a problem and could not process your request"
	ErrorResponse(w, r, http.StatusInternalServerError, map[string]map[string]string{"errors": {"message": message}})
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	ErrorResponse(w, r, http.StatusNotFound, map[string]map[string]string{"errors": {"message": message}})
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	ErrorResponse(w, r, http.StatusMethodNotAllowed, map[string]map[string]string{"errors": {"message": message}})
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	ErrorResponse(w, r, http.StatusBadRequest, map[string]map[string]string{"errors": {"message": "bad request", "error": err.Error()}})
}

func FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	errors["message"] = "encountered errors"
	ErrorResponse(w, r, http.StatusUnprocessableEntity, map[string]map[string]string{"errors": errors})
}
