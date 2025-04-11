// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers_test

import (
	"os"
	"testing"

	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

// FuzzCalculateSHA256 is a fuzz test for the CalculateSHA256 function.
func FuzzCalculateSHA256(f *testing.F) {

	tmp := f.TempDir()
	expected := tmp + "/test.txt"
	err := os.WriteFile(expected, []byte("same"), 0644)
	if err != nil {
		f.Fatal(err)
	}

	// Add initial seed corpus
	f.Add(expected)
	f.Add("") // Empty string
	f.Add("/invalid/path/to/file")

	// Add additional edge cases to the seed corpus
	f.Add("/path/with/illegal|char")
	f.Add("/path/with/<>*?")
	f.Add("/path/with/\\backslashes")
	f.Add("/dev/null")                    // Reserved name on Linux
	f.Add("CON")                          // Reserved name on Windows
	f.Add(string(make([]byte, 300, 300))) // Extremely long path
	f.Add("../relative/path")
	f.Add("./current/dir")
	f.Add(" ")                    // Single space
	f.Add("\n")                   // Newline character
	f.Add("Z:/nonexistent/drive") // Nonexistent drive on Windows
	f.Add("//network/share")
	f.Add("\\\\network\\share")
	f.Add("/dev/random")  // Special device file on Linux
	f.Add("/dev/urandom") // Special device file on Linux

	f.Fuzz(func(t *testing.T, filePath string) {
		// Call the function with the fuzzed input
		_, err := helpers.CalculateSHA256(filePath)

		// Ensure the function does not panic and handles errors gracefully
		if err != nil {
			t.Logf("Expected error for input '%s': %v", filePath, err)
		}
	})
}
