package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseGenericTypeName(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantBase     string
		wantArgs     []string
		wantIsGeneric bool
	}{
		{
			name:          "simple generic",
			input:         "ResponseData[CreditResponse]",
			wantBase:      "ResponseData",
			wantArgs:      []string{"CreditResponse"},
			wantIsGeneric: true,
		},
		{
			name:          "package-qualified generic",
			input:         "response.ResponseData[model.CreditResponse]",
			wantBase:      "response.ResponseData",
			wantArgs:      []string{"model.CreditResponse"},
			wantIsGeneric: true,
		},
		{
			name:          "multiple type args",
			input:         "Pair[model.User,model.Role]",
			wantBase:      "Pair",
			wantArgs:      []string{"model.User", "model.Role"},
			wantIsGeneric: true,
		},
		{
			name:          "not generic – plain type",
			input:         "model.CreditResponse",
			wantIsGeneric: false,
		},
		{
			name:          "not generic – empty brackets treated as array",
			input:         "[]model.User",
			wantIsGeneric: false,
		},
		{
			name:          "not generic – bracket not at end",
			input:         "ResponseData[T]Extra",
			wantIsGeneric: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base, args, ok := parseGenericTypeName(tt.input)
			assert.Equal(t, tt.wantIsGeneric, ok)
			if tt.wantIsGeneric {
				assert.Equal(t, tt.wantBase, base)
				assert.Equal(t, tt.wantArgs, args)
			}
		})
	}
}

func Test_splitTypeArgs(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"CreditResponse", []string{"CreditResponse"}},
		{"model.User,model.Role", []string{"model.User", "model.Role"}},
		{"Pair[A,B],C", []string{"Pair[A,B]", "C"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := splitTypeArgs(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_substituteTypeParams(t *testing.T) {
	params := map[string]string{"T": "model.CreditResponse"}

	assert.Equal(t, "model.CreditResponse", substituteTypeParams("T", params))
	assert.Equal(t, "[]model.CreditResponse", substituteTypeParams("[]T", params))
	assert.Equal(t, "map[]model.CreditResponse", substituteTypeParams("map[]T", params))
	assert.Equal(t, "string", substituteTypeParams("string", params))
	assert.Equal(t, "T", substituteTypeParams("T", nil))
}
