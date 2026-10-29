// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package managedsolution

import (
	"strings"
	"testing"

	"github.com/microsoft/terraform-provider-power-platform/internal/services/solution"
)

func TestUnitValidateDependenciesIgnoresBuiltInSolutions(t *testing.T) {
	required := map[string]string{
		"BaseCustomControlsCore": "9.0.2606.1003",
		"BaseLib":                "2.0.0.0",
	}

	installed := []solution.SolutionDto{
		{Name: "BaseCustomControlsCore", Version: "9.0.2605.4006"},
		{Name: "BaseLib", Version: "2.0.0.0"},
	}

	if err := validateDependencies(required, installed); err != nil {
		t.Fatalf("expected built-in dependency version drift to be ignored, got %v", err)
	}
}

func TestUnitValidateDependenciesStillFailsForCustomDependencies(t *testing.T) {
	required := map[string]string{
		"BaseCustomControlsCore": "9.0.2606.1003",
		"BaseLib":                "2.0.0.0",
	}

	installed := []solution.SolutionDto{
		{Name: "BaseCustomControlsCore", Version: "9.0.2605.4006"},
		{Name: "BaseLib", Version: "1.0.0.0"},
	}

	err := validateDependencies(required, installed)
	if err == nil {
		t.Fatal("expected custom dependency version drift to fail")
	}
	if !strings.Contains(err.Error(), `required dependency "BaseLib" >= 2.0.0.0 but 1.0.0.0 is installed`) {
		t.Fatalf("expected BaseLib validation error, got %v", err)
	}
	if strings.Contains(err.Error(), "BaseCustomControlsCore") {
		t.Fatalf("expected built-in dependency to be omitted from validation errors, got %v", err)
	}
}
