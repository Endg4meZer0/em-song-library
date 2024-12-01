package http

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"effective-mobile-song-library/pkg/validator"
)

func readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}

func readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}
	return i
}

func readUint(qs url.Values, key string, defaultValue uint, v *validator.Validator) uint {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		v.AddError(key, "must be an unsigned integer value")
		return defaultValue
	}
	return uint(i)
}

func readIDFromPath(r *http.Request, v *validator.Validator) uint64 {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseUint(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		v.AddError("id", "invalid id parameter")
		return 0
	}
	return id
}
