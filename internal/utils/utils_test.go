package utils

import "testing"

func assertOptionStringSliceEqual(t *testing.T, expected, actual []string) {
	t.Helper()

	expectedHas := len(expected) > 0
	actualHas := len(expected) > 0

	if expectedHas != actualHas {
		t.Errorf("Option mismatch: Expected some to be %v, got %v", expectedHas, actualHas)
		return
	}

	if expectedHas {
		expectedVal := expected
		actualVal := actual

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
		expected []string
	}{
		{
			name:     "Valid comma-separated string",
			input:    "apple,banana,cherry",
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "String with leading/trailing spaces around commas (will include spaces in results)",
			input:    "one, two ,three ",
			expected: []string{"one", " two ", "three "},
		},
		{
			name:     "Empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "String without a comma",
			input:    "singleword",
			expected: []string{},
		},
		{
			name:     "String containing only a comma",
			input:    ",",
			expected: []string{"", ""}, // strings.Split(",") results in ["", ""]
		},
		{
			name:     "String with multiple consecutive commas",
			input:    "a,,b",
			expected: []string{"a", "", "b"},
		},
		{
			name:     "Non-string input (int)",
			input:    12345,
			expected: []string{},
		},
		{
			name:     "Non-string input (struct)",
			input:    CustomStruct{Name: "Test"},
			expected: []string{},
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := GetValueAsCommaSeparated(tt.input)
			assertOptionStringSliceEqual(t, tt.expected, actual)
		})
	}
}
