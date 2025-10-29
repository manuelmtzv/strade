package i18n

import (
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func NewBundle() *i18n.Bundle {
	bundle := i18n.NewBundle(language.Spanish)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	bundle.MustLoadMessageFile("i18n/active.es.json")
	bundle.MustLoadMessageFile("i18n/active.en.json")

	return bundle
}
