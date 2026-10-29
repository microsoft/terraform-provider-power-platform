package managedsolution

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestUnitSourcesAreEquivalent_IgnoresQueryAndFragment(t *testing.T) {
	t.Parallel()

	plan := &SourceModel{
		URL: types.StringValue("https://downloads.example.test/source/current.zip?sig=new-token#one"),
	}
	state := &SourceModel{
		URL: types.StringValue("https://downloads.example.test/source/current.zip?sv=2026-01-01&sig=old-token#two"),
	}

	if !sourcesAreEquivalent(plan, state) {
		t.Fatal("expected URLs that differ only by query string or fragment to be treated as equivalent")
	}
}

func TestUnitSourcesAreEquivalent_IgnoresTransportLocatorChangeForExactIdentity(t *testing.T) {
	t.Parallel()

	plan := &SourceModel{
		URL: types.StringValue("https://downloads.example.test/source/current.zip?sig=new-token"),
	}
	state := &SourceModel{
		URL: types.StringValue("http://localhost:45265/archive/legacy.zip"),
	}

	if !sourcesAreEquivalent(plan, state) {
		t.Fatal("expected delivery URLs to be equivalent when managed solution identity and version are unchanged")
	}
}

func TestUnitSourcesAreEquivalent_IgnoresEphemeralWorkspacePathChange(t *testing.T) {
	t.Parallel()

	plan := &SourceModel{
		Path: types.StringValue("/tmp/new-package.zip"),
	}
	state := &SourceModel{
		Path: types.StringValue("/tmp/old-package.zip"),
	}

	if !sourcesAreEquivalent(plan, state) {
		t.Fatal("expected run-local source paths to be equivalent when managed solution identity and version are unchanged")
	}
}

func TestUnitSourcesAreEquivalent_RejectsMissingDeliverySource(t *testing.T) {
	t.Parallel()

	plan := &SourceModel{}
	state := &SourceModel{Path: types.StringValue("/tmp/package.zip")}

	if sourcesAreEquivalent(plan, state) {
		t.Fatal("expected a missing configured delivery source to remain a plan difference")
	}
}

func TestUnitNormalizeSolutionVersion_PadsMissingSegments(t *testing.T) {
	t.Parallel()

	normalized, err := normalizeSolutionVersion("1.3.5")
	if err != nil {
		t.Fatalf("expected version normalization to succeed: %v", err)
	}

	if normalized != "1.3.5.0" {
		t.Fatalf("expected normalized version to be 1.3.5.0, got %s", normalized)
	}
}

func TestUnitReconcileSolutionVersion_PreservesEquivalentDeclaredVersion(t *testing.T) {
	t.Parallel()

	actual := reconcileSolutionVersion(types.StringValue("0.1.39"), "0.1.39.0")

	if actual.ValueString() != "0.1.39" {
		t.Fatalf("expected equivalent declared version to remain 0.1.39, got %s", actual.ValueString())
	}
}

func TestUnitReconcileSolutionVersion_ReportsRemoteVersionDrift(t *testing.T) {
	t.Parallel()

	actual := reconcileSolutionVersion(types.StringValue("0.1.39"), "0.1.40.0")

	if actual.ValueString() != "0.1.40.0" {
		t.Fatalf("expected remote version drift to be reflected as 0.1.40.0, got %s", actual.ValueString())
	}
}
