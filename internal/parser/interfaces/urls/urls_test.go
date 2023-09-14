package urls_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bearer/bearer/internal/parser/interfaces/urls"
	"github.com/bearer/bearer/internal/report/values"
	"github.com/bearer/bearer/internal/report/variables"
)

func TestValueIsRelevant(t *testing.T) {
	tests := []struct {
		name     string
		value    *values.Value
		expected bool
	}{
		{"empty", &values.Value{Parts: []values.Part{}}, false},
		{"no_protocol_no_tld", &values.Value{Parts: []values.Part{values.NewStringPart("blah")}}, false},
		{"no_protocol_tld_no_dot", &values.Value{Parts: []values.Part{values.NewStringPart("com")}}, false},
		{"path_like", &values.Value{Parts: []values.Part{values.NewStringPart("../example.com")}}, false},
		{"invalid_char", &values.Value{Parts: []values.Part{values.NewStringPart("http://example^.com")}}, false},
		{"protocol_no_tld", &values.Value{Parts: []values.Part{values.NewStringPart("http://blah")}}, true},
		{"no_protocol_tld", &values.Value{Parts: []values.Part{values.NewStringPart("example.com")}}, true},
		{"allowed_domain", &values.Value{Parts: []values.Part{values.NewStringPart("example.local")}}, true},
		{"protocol_only_variables", &values.Value{Parts: []values.Part{
			values.NewStringPart("https://"),
			values.NewUnknownPart([]values.Part{}),
			values.NewStringPart("."),
			values.NewUnknownPart([]values.Part{}),
		}}, false},
		{"variables_in_value", &values.Value{Parts: []values.Part{
			values.NewVariableReferencePart(variables.VariableEnvironment, "MY_VAR"),
			values.NewStringPart(".local"),
		}}, true},
		{"variable_relevant", &values.Value{Parts: []values.Part{
			values.NewVariableReferencePart(variables.VariableEnvironment, "USER_HOST"),
			values.NewStringPart("/123"),
		}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, urls.ValueIsRelevant(tt.value))
		})
	}
}
