package utils

import (
	"hexagon-architecture/internal/errors"
	"regexp"
	"strings"

	"github.com/rl404/fairy/validation"
)

var val validation.Validator

func init() {
	val = validation.New(true)
	val.RegisterModifier("no_space", modNoSpace)
	val.RegisterValidator("alpha", valAlpha)
	val.RegisterValidatorError("required", valErrRequired)
	val.RegisterValidatorError("gte", valErrGTE)
	val.RegisterValidatorError("gt", valErrGT)
	val.RegisterValidatorError("lte", valErrLTE)
	val.RegisterValidatorError("lt", valErrLT)
	val.RegisterValidatorError("len", valErrLen)
	val.RegisterValidatorError("email", valErrEmail)
	val.RegisterValidatorError("url", valErrURL)
	val.RegisterValidatorError("oneof", valErrOneOf)
	val.RegisterValidatorError("iso3166_1_alpha2", valErrISO3166)
	val.RegisterValidatorError("numeric", valErrNumeric)
	val.RegisterValidatorError("alpha", valErrAlpha)
}

func Validate(data interface{}) error {
	return val.Validate(data)
}

func modNoSpace(in string) string {
	return strings.ReplaceAll(in, " ", "")
}

func valAlpha(f interface{}, param ...string) bool {
	return regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(f.(string))
}

func valErrRequired(f string, param ...string) error {
	return errors.ErrRequiredField(camelToSnake(f))
}

func valErrGTE(f string, param ...string) error {
	return errors.ErrGTEField(camelToSnake(f), param[0])
}

func valErrGT(f string, param ...string) error {
	return errors.ErrGTField(camelToSnake(f), param[0])
}

func valErrLTE(f string, param ...string) error {
	return errors.ErrLTEField(camelToSnake(f), param[0])
}

func valErrLT(f string, param ...string) error {
	return errors.ErrLTField(camelToSnake(f), param[0])
}

func valErrLen(f string, param ...string) error {
	return errors.ErrLenField(camelToSnake(f), param[0])
}

func valErrEmail(f string, param ...string) error {
	return errors.ErrEmailField(camelToSnake(f))
}

func valErrURL(f string, param ...string) error {
	return errors.ErrURLField(camelToSnake(f))
}

func valErrOneOf(f string, param ...string) error {
	return errors.ErrOneOfField(camelToSnake(f), param[0])
}

func valErrISO3166(f string, param ...string) error {
	return errors.ErrISO3166Alpha2Field(camelToSnake(f))
}

func valErrNumeric(f string, param ...string) error {
	return errors.ErrNumericField(camelToSnake(f))
}

func valErrAlpha(f string, param ...string) error {
	return errors.ErrAlphaField(camelToSnake(f))
}

func camelToSnake(name string) string {
	if name == "" {
		return ""
	}

	var (
		commonInitialisms         = []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
		commonInitialismsReplacer *strings.Replacer
	)

	commonInitialismsForReplacer := make([]string, 0, len(commonInitialisms))
	for _, initialism := range commonInitialisms {
		commonInitialismsForReplacer = append(commonInitialismsForReplacer, initialism, strings.Title(strings.ToLower(initialism)))
	}
	commonInitialismsReplacer = strings.NewReplacer(commonInitialismsForReplacer...)

	var (
		value                          = commonInitialismsReplacer.Replace(name)
		buf                            strings.Builder
		lastCase, nextCase, nextNumber bool
		curCase                        = value[0] <= 'Z' && value[0] >= 'A'
	)

	for i, v := range value[:len(value)-1] {
		nextCase = value[i+1] <= 'Z' && value[i+1] >= 'A'
		nextNumber = value[i+1] >= '0' && value[i+1] <= '9'

		if curCase {
			if lastCase && (nextCase || nextNumber) {
				buf.WriteRune(v + 32)
			} else {
				if i > 0 && value[i-1] != '_' && value[i+1] != '_' {
					buf.WriteByte('_')
				}
				buf.WriteRune(v + 32)
			}
		} else {
			buf.WriteRune(v)
		}

		lastCase = curCase
		curCase = nextCase
	}

	if curCase {
		if !lastCase && len(value) > 1 {
			buf.WriteByte('_')
		}
		buf.WriteByte(value[len(value)-1] + 32)
	} else {
		buf.WriteByte(value[len(value)-1])
	}

	return buf.String()
}
