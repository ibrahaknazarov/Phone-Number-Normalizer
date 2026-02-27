package main

import (
	"testing"
)

func TestNormalize(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{"1234567890", "1234567890"},
		{"123 456 7891", "1234567891"},
		{"(123) 456 7892", "1234567892"},
		{"(123) 456-7893", "1234567893"},
		{"123-456-7894", "1234567894"},
		{"(123)456-7892", "1234567892"},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			actual := normalize(tc.input)
			if actual != tc.want {
				t.Errorf("got %s; want %s", actual, tc.want)
			}
		})
	}
}

// TestNormalizeEdgeCases tests edge cases for the normalize function.
func TestNormalizeEdgeCases(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", ""},
		{"only digits", "1234567890", "1234567890"},
		{"with +1 prefix", "+11234567890", "11234567890"},
		{"with ext", "123-456-7890 ext 123", "1234567890123"},
		{"with letters", "1-800-FLOWERS", "1800"},
		{"with special chars", "1!2@3#4$5%6^7&8*9(0)", "1234567890"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := normalize(tc.input)
			if actual != tc.want {
				t.Errorf("normalize(%q) = %q; want %q", tc.input, actual, tc.want)
			}
		})
	}
}
