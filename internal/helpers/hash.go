// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// CalculateSHA256 computes the SHA256 checksum of the file specified by filePath.
// If the file does not exist, it returns an empty string without an error.
// For other errors (e.g., permission issues), it returns the error.
func CalculateSHA256(filePath string) (string, error) {
	// Attempt to open the file directly
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist; return empty checksum
			return "", nil
		}
		// An unexpected error occurred; return it
		return "", fmt.Errorf("failed to open file '%s': %w", filePath, err)
	}
	// Ensure the file is closed when the function exits
	defer file.Close()

	// Create a new SHA256 hash instance
	hash := sha256.New()

	// Copy the file content into the hash
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to read file '%s': %w", filePath, err)
	}

	// Compute the final hash value and encode it as a hexadecimal string
	return hex.EncodeToString(hash.Sum(nil)), nil
}
