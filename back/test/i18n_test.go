package test

import (
	"errors"
	"fmt"
	"library/i18n"
	"testing"
)

// These tests check the language selection and the translation of the API messages

func TestFromHeader(t *testing.T) {
	cases := map[string]i18n.Lang{
		"":                        i18n.EN,
		"fr":                      i18n.FR,
		"fr-FR,fr;q=0.9,en;q=0.8": i18n.FR,
		"en-US,en;q=0.9":          i18n.EN,
		"de-DE,fr;q=0.8":          i18n.FR,
		"de-DE":                   i18n.EN,
	}
	for header, expected := range cases {
		if got := i18n.FromHeader(header); got != expected {
			t.Errorf("FromHeader(%q): expected %s, got %s", header, expected, got)
		}
	}
}

func TestT_TranslationWithArgs(t *testing.T) {
	if got := i18n.T(i18n.EN, "loans.overdue", 2); got != "You have 2 overdue loan(s). Please return them before borrowing." {
		t.Errorf("Unexpected EN: %s", got)
	}
	if got := i18n.T(i18n.FR, "loans.overdue", 2); got != "Vous avez 2 emprunt(s) en retard, veuillez les retourner avant d'emprunter." {
		t.Errorf("Unexpected FR: %s", got)
	}
}

func TestT_UnknownKey(t *testing.T) {
	if got := i18n.T(i18n.FR, "unknown.key"); got != "unknown.key" {
		t.Errorf("expected the key, got %s", got)
	}
}

func TestErrorMessage(t *testing.T) {
	businessErr := fmt.Errorf("context: %w", i18n.NewError("copies.already_borrowed"))
	if got := i18n.ErrorMessage(i18n.FR, businessErr, "common.internal"); got != "Cet exemplaire est déjà emprunté." {
		t.Errorf("business error wrongly translated: %s", got)
	}
	internalErr := errors.New("pq: connection refused")
	if got := i18n.ErrorMessage(i18n.EN, internalErr, "common.internal"); got != "An internal error occurred." {
		t.Errorf("an internal error must not be exposed: %s", got)
	}
}
