package utils

import (
	"testing"

	"github.com/javiorfo/nilo"
)

func assertOptionStringSliceEqual(t *testing.T, expected, actual nilo.Option[[]string]) {
	t.Helper()

	expectedHas := expected.IsValue()
	actualHas := actual.IsValue()

	if expectedHas != actualHas {
		t.Errorf("Option mismatch: Expected IsValue() to be %v, got %v", expectedHas, actualHas)
		return
	}

	if expectedHas {
		expectedVal := expected.AsValue()
		actualVal := actual.AsValue()

		if len(expectedVal) != len(actualVal) {
			t.Errorf("Slice length mismatch: Expected %d elements, got %d", len(expectedVal), len(actualVal))
			return
		}

		for i := range expectedVal {
			if expectedVal[i] != actualVal[i] {
				t.Errorf("Slice element mismatch at index %d: Expected %q, got %q", i, expectedVal[i], actualVal[i])
				return
			}
		}
	}
}

func TestGetValueAsCommaSeparated(t *testing.T) {
	type CustomStruct struct {
		Name string
	}

	tests := []struct {
		name     string
		input    any
		expected nilo.Option[[]string]
	}{
		{
			name:     "Valid comma-separated string",
			input:    "apple,banana,cherry",
			expected: nilo.Value([]string{"apple", "banana", "cherry"}),
		},
		{
			name:     "String with leading/trailing spaces around commas (will include spaces in results)",
			input:    "one, two ,three ",
			expected: nilo.Value([]string{"one", " two ", "three "}),
		},
		{
			name:     "Empty string",
			input:    "",
			expected: nilo.Nil[[]string](),
		},
		{
			name:     "String without a comma",
			input:    "singleword",
			expected: nilo.Nil[[]string](),
		},
		{
			name:     "String containing only a comma",
			input:    ",",
			expected: nilo.Value([]string{"", ""}), // strings.Split(",") results in ["", ""]
		},
		{
			name:     "String with multiple consecutive commas",
			input:    "a,,b",
			expected: nilo.Value([]string{"a", "", "b"}),
		},
		{
			name:     "Non-string input (int)",
			input:    12345,
			expected: nilo.Nil[[]string](),
		},
		{
			name:     "Non-string input (struct)",
			input:    CustomStruct{Name: "Test"},
			expected: nilo.Nil[[]string](),
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: nilo.Nil[[]string](),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := GetValueAsCommaSeparated(tt.input)
			assertOptionStringSliceEqual(t, tt.expected, actual)
		})
	}
}
