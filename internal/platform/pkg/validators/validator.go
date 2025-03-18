package validators

import (
	"blog/internal/domain/users"
	"blog/internal/platform/pkg/lang"
	"errors"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Result struct {
	Fails    bool
	Messages TranslatedErrors
}

type Validator struct {
	validator      *validator.Validate
	userRepository users.Repository
}

func NewCustomValidator(u users.Repository) *Validator {
	v := validator.New()
	_ = binding.Validator.Engine().(*validator.Validate)
	lang.SetBundle()
	lang.SetTranslationFilesPath(getTranslationFilePath())
	return &Validator{validator: v, userRepository: u}
}

func (v *Validator) Validate(i interface{}, language string) Result {
	err := v.validator.Struct(i)
	if err != nil {
		res := v.createErrorResult(err)
		res.translatedErr = v.translate(res, language)
		return Result{Fails: res.Fails(), Messages: res.translatedErr}
	}
	return Result{Fails: false}
}

func (v *Validator) RegisterValidation() error {
	if err := v.validator.RegisterValidation("emailOrUsername", v.validateEmailOrUsername); err != nil {
		return err
	}
	if err := v.validator.RegisterValidation("username", v.validateUsername); err != nil {
		return err
	}
	if err := v.validator.RegisterValidation("strongPassword", v.validateStrongPassword); err != nil {
		return err
	}
	if err := v.validator.RegisterValidation("uniqueEmail", v.validateUniqueEmail); err != nil {
		return err
	}
	if err := v.validator.RegisterValidation("uniqueUsername", v.validateUniqueUsername); err != nil {
		return err
	}
	return nil
}

func (v *Validator) createErrorResult(err error) result {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return result{errors: validationErrors}
	}
	return result{}
}

func (v *Validator) translate(r result, language string) TranslatedErrors {
	messages := make(map[string][]string)
	for _, e := range r.errors {
		validatorOperator := e.Tag()
		validateField := e.Field()
		message := lang.TryBy(language, validatorOperator)
		message = strings.Replace(message, "{field}", validateField, -1)
		messages[toCamelCase(validateField)] = append(messages[toCamelCase(validateField)], message)
	}
	return TranslatedErrors{messages: messages}
}

func getTranslationFilePath() string {
	_, currentFilePath, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFilePath)
	return filepath.Join(dir, "translations")
}

func toCamelCase(s string) string {
	var builder strings.Builder
	upperNext := false

	for i, r := range s {
		switch {
		case i == 0:
			builder.WriteRune(unicode.ToLower(r))
		case r == '_', r == ' ', r == '-':
			upperNext = true
		default:
			if upperNext {
				builder.WriteRune(unicode.ToUpper(r))
				upperNext = false
			} else {
				builder.WriteRune(r)
			}
		}
	}

	return builder.String()
}
