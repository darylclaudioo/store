package utils_test

import (
	"hexagon-architecture/internal/errors"
	"hexagon-architecture/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyStruct struct {
	NoSpaceMod string `mod:"no_space"`
	Required   string `validate:"required"`
	Gte        int    `validate:"gte=0"`
	Gt         int    `validate:"gt=0"`
	Lte        int    `validate:"lte=0"`
	Lt         int    `validate:"lt=0"`
	Len        string `validate:"len=5"`
	Email      string `validate:"email"`
	Url        string `validate:"url"`
	OneOf      string `validate:"oneof=a b c"`
	Iso        string `validate:"iso3166_1_alpha2"`
	Numeric    string `validate:"numeric"`
	Alpha      string `validate:"alpha"`
}

func TestValidate(t *testing.T) {
	expectedStruct := dummyStruct{
		NoSpaceMod: "abcdef",
		Required:   "abcdef",
		Gte:        1,
		Gt:         1,
		Lte:        -1,
		Lt:         -1,
		Len:        "12345",
		Email:      "string@a.com",
		Url:        "https://a.com",
		OneOf:      "a",
		Iso:        "ID",
		Numeric:    "12345",
		Alpha:      "abcde",
	}

	tests := map[string]struct {
		input    dummyStruct
		expected error
	}{
		"no-space-mod": {
			input: dummyStruct{
				NoSpaceMod: "ab cd ef",
				Required:   "abcdef",
				Gte:        1,
				Gt:         1,
				Lte:        -1,
				Lt:         -1,
				Len:        "12345",
				Email:      "string@a.com",
				Url:        "https://a.com",
				OneOf:      "a",
				Iso:        "ID",
				Numeric:    "12345",
				Alpha:      "abcde",
			},
			expected: nil,
		},
		"required": {
			input: dummyStruct{
				NoSpaceMod: "abcdef",
				Required:   "",
				Gte:        1,
				Gt:         1,
				Lte:        -1,
				Lt:         -1,
				Len:        "12345",
				Email:      "string@a.com",
				Url:        "https://a.com",
				OneOf:      "a",
				Iso:        "ID",
				Numeric:    "12345",
				Alpha:      "abcde",
			},
			expected: errors.ErrRequiredField("required"),
		},
		"gte": {
			input: dummyStruct{
				NoSpaceMod: "abcdef",
				Required:   "abcdef",
				Gte:        -1,
				Gt:         1,
				Lte:        -1,
				Lt:         -1,
				Len:        "12345",
				Email:      "string@a.com",
				Url:        "https://a.com",
				OneOf:      "a",
				Iso:        "ID",
				Numeric:    "12345",
				Alpha:      "abcde",
			},
			expected: errors.ErrGTEField("gte", "0"),
		},
		"gt": {
			input: dummyStruct{
				NoSpaceMod: "abcdef",
				Required:   "abcdef",
				Gte:        1,
				Gt:         0,
				Lte:        -1,
				Lt:         -1,
				Len:        "12345",
				Email:      "string@a.com",
				Url:        "https://a.com",
				OneOf:      "a",
				Iso:        "ID",
				Numeric:    "12345",
				Alpha:      "abcde",
			},
			expected: errors.ErrGTField("gt", "0"),
		},
		"lte": {
			input: dummyStruct{
				NoSpaceMod: "abcdef",
				Required:   "abcdef",
				Gte:        1,
				Gt:         1,
				Lte:        1,
				Lt:         -1,
				Len:        "12345",
				Email:      "string@a.com",
				Url:        "https://a.com",
				OneOf:      "a",
				Iso:        "ID",
				Numeric:    "12345",
				Alpha:      "abcde",
			},
			expected: errors.ErrLTEField("lte", "0"),
		},
		"lt": {
			input: dummyStruct{
				NoSpaceMod: "abcdef",
				Required:   "abcdef",
				Gte:        1,
				Gt:         1,
				Lte:        -1,
				Lt:         0,
				Len:        "12345",
				Email:      "string@a.com",
				Url:        "https://a.com",
				OneOf:      "a",
				Iso:        "ID",
				Numeric:    "12345",
				Alpha:      "abcde",
			},
			expected: errors.ErrLTField("lt", "0"),
		},
		"len": {
			input: dummyStruct{
				NoSpaceMod: "abcdef",
				Required:   "abcdef",
				Gte:        1,
				Gt:         1,
				Lte:        -1,
				Lt:         -1,
				Len:        "1234",
				Email:      "string@a.com",
				Url:        "https://a.com",
				OneOf:      "a",
				Iso:        "ID",
				Numeric:    "12345",
				Alpha:      "abcde",
			},
			expected: errors.ErrLenField("len", "5"),
		},
		"email": {
			input: dummyStruct{
				NoSpaceMod: "abcdef",
				Required:   "abcdef",
				Gte:        1,
				Gt:         1,
				Lte:        -1,
				Lt:         -1,
				Len:        "12345",
				Email:      "string",
				Url:        "https://a.com",
				OneOf:      "a",
				Iso:        "ID",
				Numeric:    "12345",
				Alpha:      "abcde",
			},
			expected: errors.ErrEmailField("email"),
		},
		"url": {
			input: dummyStruct{
				NoSpaceMod: "abcdef",
				Required:   "abcdef",
				Gte:        1,
				Gt:         1,
				Lte:        -1,
				Lt:         -1,
				Len:        "12345",
				Email:      "string@a.com",
				Url:        "a.com",
				OneOf:      "a",
				Iso:        "ID",
				Numeric:    "12345",
				Alpha:      "abcde",
			},
			expected: errors.ErrURLField("url"),
		},
		"one-of": {
			input: dummyStruct{
				NoSpaceMod: "abcdef",
				Required:   "abcdef",
				Gte:        1,
				Gt:         1,
				Lte:        -1,
				Lt:         -1,
				Len:        "12345",
				Email:      "string@a.com",
				Url:        "https:a.com",
				OneOf:      "1",
				Iso:        "ID",
				Numeric:    "12345",
				Alpha:      "abcde",
			},
			expected: errors.ErrOneOfField("one_of", "a b c"),
		},
		"iso": {
			input: dummyStruct{
				NoSpaceMod: "abcdef",
				Required:   "abcdef",
				Gte:        1,
				Gt:         1,
				Lte:        -1,
				Lt:         -1,
				Len:        "12345",
				Email:      "string@a.com",
				Url:        "https:a.com",
				OneOf:      "a",
				Iso:        "IDS",
				Numeric:    "12345",
				Alpha:      "abcde",
			},
			expected: errors.ErrISO3166Alpha2Field("iso"),
		},
		"numeric": {
			input: dummyStruct{
				NoSpaceMod: "abcdef",
				Required:   "abcdef",
				Gte:        1,
				Gt:         1,
				Lte:        -1,
				Lt:         -1,
				Len:        "12345",
				Email:      "string@a.com",
				Url:        "https:a.com",
				OneOf:      "a",
				Iso:        "ID",
				Numeric:    "asd",
				Alpha:      "abcde",
			},
			expected: errors.ErrNumericField("numeric"),
		},
		"alpha": {
			input: dummyStruct{
				NoSpaceMod: "abcdef",
				Required:   "abcdef",
				Gte:        1,
				Gt:         1,
				Lte:        -1,
				Lt:         -1,
				Len:        "12345",
				Email:      "string@a.com",
				Url:        "https:a.com",
				OneOf:      "a",
				Iso:        "ID",
				Numeric:    "12345",
				Alpha:      "123",
			},
			expected: errors.ErrAlphaField("alpha"),
		},
		"ok": {
			input: dummyStruct{
				NoSpaceMod: "abcdef",
				Required:   "abcdef",
				Gte:        1,
				Gt:         1,
				Lte:        -1,
				Lt:         -1,
				Len:        "12345",
				Email:      "string@a.com",
				Url:        "https://a.com",
				OneOf:      "a",
				Iso:        "ID",
				Numeric:    "12345",
				Alpha:      "abcde",
			},
			expected: nil,
		},
	}

	for _, test := range tests {
		err := utils.Validate(&test.input)
		assert.Equal(t, test.expected, err)
		if err == nil {
			assert.Equal(t, expectedStruct, test.input)
		}
	}
}
