// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers_test

import (
	"os"
	"testing"

	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

func TestUnitCalculateSHA256(t *testing.T) {
	t.Parallel()

	tdir := t.TempDir()
	file1 := tdir + "/test.txt"
	file2 := tdir + "/test2.txt"
	file3 := tdir + "/test3.txt"
	file4 := tdir + "/test4.txt"
	file5 := tdir + "/test5"

	err := os.WriteFile(file1, []byte("same"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(file2, []byte("same"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(file3, []byte("different"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Mkdir(file5, 0644)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("TestUnitCalculateSHA256_SameFile", func(t *testing.T) {
		t.Parallel()

		// Test code here
		f1, err := helpers.CalculateSHA256(file1)
		if err != nil {
			t.Fatal(err)
		}

		f1b, err := helpers.CalculateSHA256(file1)
		if err != nil {
			t.Fatal(err)
		}

		if f1 != f1b {
			t.Errorf("Expected %s to equal %s", f1, f1b)
		}
	})

	t.Run("TestUnitCalculateSHA256_SameContent", func(t *testing.T) {
		t.Parallel()

		// Test code here
		f1, err := helpers.CalculateSHA256(file1)
		if err != nil {
			t.Fatal(err)
		}

		f2, err := helpers.CalculateSHA256(file2)
		if err != nil {
			t.Fatal(err)
		}

		if f1 != f2 {
			t.Errorf("Expected %s to equal %s", f1, f2)
		}
	})

	t.Run("TestUnitCalculateSHA256_DifferentContent", func(t *testing.T) {
		t.Parallel()

		// Test code here
		f1, err := helpers.CalculateSHA256(file1)
		if err != nil {
			t.Fatal(err)
		}

		f3, err := helpers.CalculateSHA256(file3)
		if err != nil {
			t.Fatal(err)
		}

		if f1 == f3 {
			t.Errorf("Expected %s to not equal %s", f1, f3)
		}
	})

	t.Run("TestUnitCalculateSHA256_FileDoesNotExist", func(t *testing.T) {
		t.Parallel()

		// Test code here
		f4, err := helpers.CalculateSHA256(file4)
		if err != nil {
			t.Fatal(err)
		}

		if f4 != "" {
			t.Errorf("Expected %s to be empty", f4)
		}
	})

	t.Run("TestUnitCalculateSHA256_FileNameIsADirectory", func(t *testing.T) {
		t.Parallel()

		// Test code here
		_, err := helpers.CalculateSHA256(file5)
		if err == nil {
			t.Error("Expected an error but got nil error")
		}
	})
}
