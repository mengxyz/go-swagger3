package integration_test

import (
	"encoding/json"
	"testing"

	"github.com/parvez3019/go-swagger3/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_GenericResponseType verifies that generic response types like
// model.ResponseData[model.CreditResponse] are expanded into a proper
// OpenAPI schema instead of being silently dropped.
func Test_GenericResponseType(t *testing.T) {
	p, err := parser.NewParser(
		"test_data_generic",
		"test_data_generic/server/main.go",
		"",
		false,
		false,
		true, // schemaWithoutPkg
	).Init()
	require.NoError(t, err)

	openAPI, err := p.Parse()
	require.NoError(t, err)

	// Serialise to JSON so we can inspect the output easily.
	raw, err := json.MarshalIndent(openAPI, "", "  ")
	require.NoError(t, err)
	spec := string(raw)

	// The expanded generic schema must be present in components.schemas.
	assert.Contains(t, spec, `"ResponseData[CreditResponse]"`,
		"expected expanded generic schema ResponseData[CreditResponse] in components")

	// The $ref inside the 200 response must point at the expanded schema.
	assert.Contains(t, spec, `#/components/schemas/ResponseData[CreditResponse]`,
		"expected $ref pointing to expanded generic schema")

	// The expanded schema must contain the 'data' field coming from CreditResponse.
	assert.Contains(t, spec, `"balance"`,
		"expected 'balance' property inside the expanded generic schema")
	assert.Contains(t, spec, `"message"`,
		"expected 'message' property inside the expanded generic schema")
}
