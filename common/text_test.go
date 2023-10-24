package common

import "testing"

func TestRemoveNonNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello, 123!", "123"},
		{"This_is_a_test", ""},
		{"SpecialChars#$@123", "123"},
		{"   ", ""},                  // Test for an empty string
		{"", ""},                     // Test for an empty string
		{"1234567890", "1234567890"}, // Test with all numeric characters
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := RemoveNonNumeric(test.input)
			if result != test.expected {
				t.Errorf("Expected: %s, Got: %s", test.expected, result)
			}
		})
	}
}
