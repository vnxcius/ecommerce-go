package validators

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var NameRX = regexp.MustCompile(`^[\p{L} '-]+$`)
var NumberRX = regexp.MustCompile(`^[0-9]+$`)

// Valid() retorna true se FieldErrors map não conter nenhuma entrada.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// Cria um flash de erro.
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// Adiciona um campo de erro.
func (v *Validator) AddFieldError(key, message string) {
	// Inicia o map primeiro, se já não foi inicializado.
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField() adiciona uma mensagem de erro ao FieldErrors map
// somente se a validação não estiver ok.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank() retorna true se o valor não for uma string vazia.
func NotBlank(value interface{}) bool {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) != ""
	default:
		// Caso não reconheça o tipo, pode-se retornar false ou gerar um erro
		return false
	}
}

func NotBlankInt(value interface{}) bool {
	switch v := value.(type) {
	case int:
		return v != 0
	case int32:
		return v != 0
	case int64:
		return v != 0
	case float32:
		return v != 0
	case float64:
		return v != 0
	default:
		return false
	}
}

// MaxChars() retorna true se o valor conter não mais que N caracteres.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func MaxImages(value []string, n int) bool {
	return len(value) <= n
}

// PermittedValue() retorna true se o valor estiver em uma
// lista de valores específicos permitidos.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// MinChars() retorna true se o valor conter pelo menos N caracteres
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// Matches() retorna true se o valor corresponder ao padrão definido
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func IsEqual(value1 string, value2 string) bool {
	return value1 == value2
}
