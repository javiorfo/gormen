package utils

import (
	"testing"

	"github.com/javiorfo/nilo"
)

func assertOptionStringSliceEqual(t *testing.T, expected, actual nilo.Option[[]string]) {
	t.Helper()

	expectedHas := expected.IsSome()
	actualHas := actual.IsSome()

	if expectedHas != actualHas {
		t.Errorf("Option mismatch: Expected IsSome() to be %v, got %v", expectedHas, actualHas)
		return
	}

	if expectedHas {
		expectedVal := expected.Unwrap()
		actualVal := actual.Unwrap()

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
			expected: nilo.Some([]string{"apple", "banana", "cherry"}),
		},
		{
			name:     "String with leading/trailing spaces around commas (will include spaces in results)",
			input:    "one, two ,three ",
			expected: nilo.Some([]string{"one", " two ", "three "}),
		},
		{
			name:     "Empty string",
			input:    "",
			expected: nilo.None[[]string](),
		},
		{
			name:     "String without a comma",
			input:    "singleword",
			expected: nilo.None[[]string](),
		},
		{
			name:     "String containing only a comma",
			input:    ",",
			expected: nilo.Some([]string{"", ""}), // strings.Split(",") results in ["", ""]
		},
		{
			name:     "String with multiple consecutive commas",
			input:    "a,,b",
			expected: nilo.Some([]string{"a", "", "b"}),
		},
		{
			name:     "Non-string input (int)",
			input:    12345,
			expected: nilo.None[[]string](),
		},
		{
			name:     "Non-string input (struct)",
			input:    CustomStruct{Name: "Test"},
			expected: nilo.None[[]string](),
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: nilo.None[[]string](),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := GetValueAsCommaSeparated(tt.input)
			assertOptionStringSliceEqual(t, tt.expected, actual)
		})
	}
}
