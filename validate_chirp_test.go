// validate_chirp_test.go

package main

import (
	"testing" // importing testing package for unit tests
)

// test cleanProfanity
func TestCleanProfanity(t *testing.T) {
	// build test cases
	testCases := []struct {
		input    string
		expected string
	}{ // }{ -- inits the test values for input vs expected
		{"Hello world", "Hello world"},
		{"Hello kerfuffle world", "Hello **** world"},
		{"KeRfUfFlE is bad", "**** is bad"},
		{"Empty profanity list", "Empty profanity list"},
		{"kErFuFfLe sHaRbErT fOrNaX", "**** **** ****"},
	}

	// loop through the testcases and call the helper function
	for _, tc := range testCases {
		actual := cleanProfanity(tc.input) // get result

		// check result
		if actual != tc.expected {
			// if no match, raise test error
			t.Errorf("cleanProfanity(%q) = %q, want %q",
				tc.input, actual, tc.expected)
		}
		// %q = quoted string, preserves invisible chars like newlines etc.
		// %s doesn't, meaning we can test EXACTLY
	}
}
