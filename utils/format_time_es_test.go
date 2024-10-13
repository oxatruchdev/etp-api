package utils

import (
	"testing"
)

func TestFormatTimeInSpanish(t *testing.T) {
	tests := []struct {
		input          string
		format         string
		expectedOutput string
		shouldFail     bool
	}{
		{
			input:          "2024-10-11 19:25:31.872799 -0700 PDT",
			format:         "Monday, 2 January 2006 15:04:05",
			expectedOutput: "viernes, 11 octubre 2024 19:25:31",
			shouldFail:     false,
		},
		{
			input:          "2024-03-03 10:10:10.123456 -0700 PDT",
			format:         "Monday, 2 January 2006 15:04:05",
			expectedOutput: "domingo, 3 marzo 2024 10:10:10",
			shouldFail:     false,
		},
		{
			input:          "2024-03-03 10:10:10.123456 -0700 PDT",
			format:         "January 2 2006",
			expectedOutput: "marzo 3 2024",
			shouldFail:     false,
		},
		{
			input:          "invalid-date",
			format:         "Monday, 2 January 2006 15:04:05",
			expectedOutput: "",
			shouldFail:     true, // Parsing should fail and return an empty string
		},
	}

	for _, test := range tests {
		output := FormatTimeInSpanish(test.input, test.format)
		if test.shouldFail && output != "" {
			t.Errorf("expected failure but got output: %s", output)
		} else if !test.shouldFail && output != test.expectedOutput {
			t.Errorf("for input '%s' expected '%s' but got '%s'", test.input, test.expectedOutput, output)
		}
	}
}
