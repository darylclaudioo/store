package errors

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidDBFormat      = errors.New("invalid db address")
	ErrNotFoundBoilerplate  = errors.New("not found boilerplate")
	ErrNotFoundProduct      = errors.New("not found product")
	ErrInvalidRequestFormat = errors.New("invalid request format")
	ErrInternalDB           = errors.New("internal database error")
	ErrInternalElastic      = errors.New("internal elastic error")
	ErrInternalCache        = errors.New("internal cache error")
	ErrInternalServer       = errors.New("internal server error")
)

func ErrRequiredField(str string) error {
	return fmt.Errorf("required field %s", str)
}

func ErrGTField(str, value string) error {
	return fmt.Errorf("field %s must be greater than %s", str, value)
}

func ErrGTEField(str, value string) error {
	return fmt.Errorf("field %s must be greater than or equal %s", str, value)
}

func ErrLTField(str, value string) error {
	return fmt.Errorf("field %s must be lower than %s", str, value)
}

func ErrLTEField(str, value string) error {
	return fmt.Errorf("field %s must be lower than or equal %s", str, value)
}

func ErrLenField(str, value string) error {
	return fmt.Errorf("field %s length must be %s", str, value)
}

func ErrISO3166Alpha2Field(str string) error {
	return fmt.Errorf("field %s must be in ISO 3166-1 alpha-2 format", str)
}

func ErrEmailField(str string) error {
	return fmt.Errorf("field %s must be in email format", str)
}

func ErrURLField(str string) error {
	return fmt.Errorf("field %s must be in URL format", str)
}

func ErrInvalidFormatField(str string) error {
	return fmt.Errorf("invalid format field %s", str)
}

func ErrOneOfField(str, value string) error {
	return fmt.Errorf("field %s must be one of %s", str, strings.Join(strings.Split(value, " "), "/"))
}

func ErrNumericField(str string) error {
	return fmt.Errorf("field %s must contain number only", str)
}

func ErrAlphaField(str string) error {
	return fmt.Errorf("field %s must contain letter only", str)
}

func ErrMinPrice(price float64) error {
	return fmt.Errorf("minimum price is %d", int(price))
}

func ErrMaxPrice(price float64) error {
	return fmt.Errorf("maximum price is %d", int(price))
}

func ErrMinAmount(amount float64) error {
	return fmt.Errorf("minimum amount is %d", int(amount))
}

func ErrMaxAmount(amount float64) error {
	return fmt.Errorf("maximum amount is %d", int(amount))
}

func ErrDatetimeField(str string) error {
	return fmt.Errorf("%s must be in date format (dd/mm/yyyy)", str)
}
