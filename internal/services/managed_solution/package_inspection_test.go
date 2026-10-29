// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package managedsolution

import (
	"archive/zip"
	"os"
	"testing"

	"github.com/microsoft/terraform-provider-power-platform/internal/services/solution"
)

func TestUnitInspectSolutionPackage_ParsesMetadataInterfaceAndDependencies(t *testing.T) {
	zipPath := createTestSolutionZip(t, map[string]string{
		"solution.xml": `<ImportExportXml>
  <SolutionManifest>
    <UniqueName>CodeEditor</UniqueName>
    <LocalizedNames>
      <LocalizedName description="Code Editor" />
    </LocalizedNames>
    <Version>1.2.3.4</Version>
    <Managed>1</Managed>
    <MissingDependencies>
      <MissingDependency>
        <Required type="44" schemaName="dep_one" solution="BaseLib (1.0.0.0)" />
      </MissingDependency>
      <MissingDependency>
        <Required type="66" schemaName="dep_two" solution="BaseLib (1.4.0.0)" />
      </MissingDependency>
      <MissingDependency>
        <Required type="90" schemaName="dep_three" solution="Portal (2.0.0.0)" />
      </MissingDependency>
    </MissingDependencies>
  </SolutionManifest>
</ImportExportXml>`,
		"customizations.xml": `<ImportExportXml>
  <connectionreferences>
    <connectionreference connectionreferencelogicalname="codeeditor_sharepoint">
      <connectorid>/providers/Microsoft.PowerApps/apis/shared_sharepointonline</connectorid>
    </connectionreference>
  </connectionreferences>
</ImportExportXml>`,
		"environmentvariabledefinitions/codeeditor_text/environmentvariabledefinition.xml": `<environmentvariabledefinition schemaname="codeeditor_text">
  <defaultvalue>hello</defaultvalue>
</environmentvariabledefinition>`,
		"environmentvariabledefinitions/codeeditor_secret/environmentvariabledefinition.xml": `<environmentvariabledefinition schemaname="codeeditor_secret"></environmentvariabledefinition>`,
		"environmentvariabledefinitions/codeeditor_secret/environmentvariablevalues.json":    `{"environmentvariablevalues":{"environmentvariablevalue":{"value":"shipped"}}}`,
	})

	pkg, err := inspectSolutionPackage(zipPath)
	if err != nil {
		t.Fatalf("inspectSolutionPackage returned error: %v", err)
	}

	if pkg.UniqueName != "CodeEditor" {
		t.Fatalf("expected unique name CodeEditor, got %s", pkg.UniqueName)
	}
	if pkg.Version != "1.2.3.4" {
		t.Fatalf("expected version 1.2.3.4, got %s", pkg.Version)
	}
	if !pkg.IsManaged {
		t.Fatal("expected package to be managed")
	}
	if pkg.DisplayName != "Code Editor" {
		t.Fatalf("expected display name Code Editor, got %s", pkg.DisplayName)
	}

	if pkg.Dependencies["BaseLib"] != "1.4.0.0" {
		t.Fatalf("expected highest dependency version for BaseLib to be 1.4.0.0, got %s", pkg.Dependencies["BaseLib"])
	}
	if pkg.Dependencies["Portal"] != "2.0.0.0" {
		t.Fatalf("expected Portal dependency version 2.0.0.0, got %s", pkg.Dependencies["Portal"])
	}

	ref := pkg.ConnectionReferences["codeeditor_sharepoint"]
	if ref.ConnectorID != "/providers/Microsoft.PowerApps/apis/shared_sharepointonline" {
		t.Fatalf("expected connector id to be parsed, got %s", ref.ConnectorID)
	}

	if !pkg.EnvironmentVariables["codeeditor_text"].HasDefaultValue {
		t.Fatal("expected codeeditor_text to have a default value")
	}
	if !pkg.EnvironmentVariables["codeeditor_secret"].ContainsPackagedValue {
		t.Fatal("expected codeeditor_secret to be marked as containing a packaged current value")
	}
}

func TestUnitInspectSolutionPackage_FailsWhenDependencySolutionReferenceIsActive(t *testing.T) {
	zipPath := createTestSolutionZip(t, map[string]string{
		"solution.xml": `<ImportExportXml>
  <SolutionManifest>
    <UniqueName>CodeEditor</UniqueName>
    <Version>1.2.3.4</Version>
    <Managed>1</Managed>
    <MissingDependencies>
      <MissingDependency>
        <Required type="12" schemaName="active_placeholder" solution="Active" />
      </MissingDependency>
    </MissingDependencies>
  </SolutionManifest>
</ImportExportXml>`,
		"customizations.xml": `<ImportExportXml></ImportExportXml>`,
	})

	_, err := inspectSolutionPackage(zipPath)
	if err == nil {
		t.Fatal("expected inspectSolutionPackage to fail for Active solution dependency reference")
	}
}

func TestUnitInspectSolutionPackage_FailsWhenDependencySolutionReferenceIsEmpty(t *testing.T) {
	zipPath := createTestSolutionZip(t, map[string]string{
		"solution.xml": `<ImportExportXml>
  <SolutionManifest>
    <UniqueName>CodeEditor</UniqueName>
    <Version>1.2.3.4</Version>
    <Managed>1</Managed>
    <MissingDependencies>
      <MissingDependency>
        <Required type="13" schemaName="empty_placeholder" solution="" />
      </MissingDependency>
    </MissingDependencies>
  </SolutionManifest>
</ImportExportXml>`,
		"customizations.xml": `<ImportExportXml></ImportExportXml>`,
	})

	_, err := inspectSolutionPackage(zipPath)
	if err == nil {
		t.Fatal("expected inspectSolutionPackage to fail for empty solution dependency reference")
	}
}

func TestUnitValidateEnvironmentVariablePackaging_RejectsPackagedValue(t *testing.T) {
	err := validateEnvironmentVariablePackaging(map[string]packageEnvironmentVariable{
		"codeeditor_secret": {
			SchemaName:            "codeeditor_secret",
			ContainsPackagedValue: true,
		},
	}, map[string]bool{})
	if err == nil {
		t.Fatal("expected a packaged current value to be rejected")
	}
}

func TestUnitValidateEnvironmentVariablePackaging_RejectsUnmanagedDefinitionCapture(t *testing.T) {
	err := validateEnvironmentVariablePackaging(map[string]packageEnvironmentVariable{
		"terr_Shared": {
			SchemaName: "terr_Shared",
		},
	}, map[string]bool{"terr_Shared": true})
	if err == nil {
		t.Fatal("expected shipping a definition that exists unmanaged in the target to be rejected")
	}
}

func TestUnitValidateEnvironmentVariablePackaging_AllowsUnresolvedReferences(t *testing.T) {
	err := validateEnvironmentVariablePackaging(map[string]packageEnvironmentVariable{
		"terr_NoDefaultNoValue": {
			SchemaName: "terr_NoDefaultNoValue",
		},
	}, map[string]bool{})
	if err != nil {
		t.Fatalf("a definition with no default and no value resolves to null by platform contract and must be allowed: %v", err)
	}
}

func TestUnitValidateDependencies_FailsWhenInstalledVersionIsTooLow(t *testing.T) {
	err := validateDependencies(map[string]string{
		"CodeEditorBase": "1.2.0.0",
	}, []solution.SolutionDto{
		{
			Name:    "CodeEditorBase",
			Version: "1.1.0.0",
		},
	})
	if err == nil {
		t.Fatal("expected validateDependencies to fail")
	}
}

func TestUnitValidateDependencies_AllowsLowerBuildWhenMajorMinorPatchMatch(t *testing.T) {
	err := validateDependencies(map[string]string{
		"CodeEditorBase": "1.2.3.99",
	}, []solution.SolutionDto{
		{
			Name:    "CodeEditorBase",
			Version: "1.2.3.1",
		},
	})
	if err != nil {
		t.Fatalf("expected validateDependencies to allow lower build when major.minor.patch matches: %v", err)
	}
}

func TestUnitValidateDependencies_FailsWhenPatchIsTooLow(t *testing.T) {
	err := validateDependencies(map[string]string{
		"CodeEditorBase": "1.2.3.0",
	}, []solution.SolutionDto{
		{
			Name:    "CodeEditorBase",
			Version: "1.2.2.999",
		},
	})
	if err == nil {
		t.Fatal("expected validateDependencies to fail when installed patch version is too low")
	}
}

func createTestSolutionZip(t *testing.T, files map[string]string) string {
	t.Helper()

	file, err := os.CreateTemp("", "managed-solution-test-*.zip")
	if err != nil {
		t.Fatalf("failed to create temp zip: %v", err)
	}
	defer file.Close()

	writer := zip.NewWriter(file)
	for name, content := range files {
		entry, err := writer.Create(name)
		if err != nil {
			t.Fatalf("failed to create zip entry %s: %v", name, err)
		}
		if _, err := entry.Write([]byte(content)); err != nil {
			t.Fatalf("failed to write zip entry %s: %v", name, err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close zip writer: %v", err)
	}

	return file.Name()
}
