package lang

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	bundle               *i18n.Bundle
	defaultLocale        string
	fallbackLocale       string
	translationFilesPath string
)

type Bundle struct {
	bundle               *i18n.Bundle
	defaultLocale        string
	fallbackLocale       string
	translationFilesPath string
}

func SetBundle() {
	if bundle != nil {
		return
	}
	if defaultLocale == "" {
		defaultLocale = "en"
	}
	if fallbackLocale == "" {
		fallbackLocale = "en"
	}
	bundle = i18n.NewBundle(language.Make(defaultLocale))
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
}

func SetTranslationFilesPath(path string) {
	translationFilesPath = path
	langFiles, err := getJSONFilesFromDirectory(path)
	if err != nil {
		panic("could not get json files from path: " + err.Error())
	}
	for _, l := range langFiles {
		bundle.MustLoadMessageFile(l)
	}
}

func SetDefaultLocale(locale string) {
	defaultLocale = locale
}

func SetFallbackLocale(locale string) {
	fallbackLocale = locale
}

func newBundle() *Bundle {
	return &Bundle{
		bundle:               bundle,
		defaultLocale:        defaultLocale,
		fallbackLocale:       fallbackLocale,
		translationFilesPath: translationFilesPath,
	}
}

func TryBy(language, key string) string {
	lang := newBundle()
	return lang.tryBy(language, key)
}

func (b *Bundle) tryBy(language, key string) string {
	var lang string
	if language == "" {
		lang = b.defaultLocale
		if defaultLocale == "" {
			lang = b.fallbackLocale
		}
	}
	localizer := i18n.NewLocalizer(bundle, lang)
	message := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: key,
		},
	})
	return message
}

func getJSONFilesFromDirectory(dir string) ([]string, error) {
	var jsonFiles []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			jsonFiles = append(jsonFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return jsonFiles, nil
}
