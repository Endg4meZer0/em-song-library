package validator

import (
	"regexp"
)

var ReleaseYearRX = regexp.MustCompile(`^[0-9]{4}$`)
var ReleaseYearMonthRX = regexp.MustCompile(`^[0-1][0-9]\.[0-9]{4}$`)
var ReleaseDateRX = regexp.MustCompile(`^[0-3][0-9]\.[0-1][0-9]\.[0-9]{4}$`)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validator) Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
