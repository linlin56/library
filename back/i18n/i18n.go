// Package i18n provides the API messages in English (default) and French.
// The language is picked from the Accept-Language HTTP header.
package i18n

import (
	"errors"
	"fmt"
	"strings"
)

type Lang string

const (
	EN Lang = "en"
	FR Lang = "fr"
)

// FromHeader returns the first supported language found in an
// Accept-Language header (e.g. "fr-FR,fr;q=0.9,en;q=0.8"). English by default.
func FromHeader(header string) Lang {
	for _, part := range strings.Split(header, ",") {
		tag := strings.ToLower(strings.TrimSpace(strings.SplitN(part, ";", 2)[0]))
		switch {
		case strings.HasPrefix(tag, "fr"):
			return FR
		case strings.HasPrefix(tag, "en"):
			return EN
		}
	}
	return EN
}

// T translates a catalogue key into the requested language.
// An unknown key is returned as is.
func T(lang Lang, key string, args ...any) string {
	msg, ok := messages[key]
	if !ok {
		return key
	}
	format, ok := msg[lang]
	if !ok {
		format = msg[EN]
	}
	if len(args) == 0 {
		return format
	}
	return fmt.Sprintf(format, args...)
}

// Error is a business error whose message can be shown to the user
// in their language.
type Error struct {
	Key  string
	Args []any
}

// NewError creates a translatable business error.
func NewError(key string, args ...any) *Error {
	return &Error{Key: key, Args: args}
}

func (e *Error) Error() string {
	return T(EN, e.Key, e.Args...)
}

// ErrorMessage translates err if it is a business error (*Error).
// Otherwise, it returns the generic fallbackKey message so that no internal detail is exposed.
func ErrorMessage(lang Lang, err error, fallbackKey string) string {
	var e *Error
	if errors.As(err, &e) {
		return T(lang, e.Key, e.Args...)
	}
	return T(lang, fallbackKey)
}
