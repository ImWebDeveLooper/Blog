package validators

import "github.com/go-playground/validator/v10"

type result struct {
	errors        validator.ValidationErrors
	translatedErr TranslatedErrors
}

func (r result) Fails() bool {
	return len(r.errors) > 0
}

type TranslatedErrors struct {
	messages map[string][]string
}

func (t TranslatedErrors) All() map[string][]string {
	return t.messages
}
