package managedsolution

import (
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSourcesAreEquivalent_IgnoresQueryAndFragment(t *testing.T) {
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

func TestSourcesAreEquivalent_DetectsTransportChange(t *testing.T) {
	t.Parallel()

	plan := &SourceModel{
		URL: types.StringValue("https://downloads.example.test/source/current.zip?sig=new-token"),
	}
	state := &SourceModel{
		URL: types.StringValue("http://localhost:45265/archive/legacy.zip"),
	}

	if sourcesAreEquivalent(plan, state) {
		t.Fatal("expected URLs with different transport and path values to be treated as different")
	}
}

func TestSourcesAreEquivalent_DetectsPathChange(t *testing.T) {
	t.Parallel()

	plan := &SourceModel{
		Path: types.StringValue("/tmp/new-package.zip"),
	}
	state := &SourceModel{
		Path: types.StringValue("/tmp/old-package.zip"),
	}

	if sourcesAreEquivalent(plan, state) {
		t.Fatal("expected different local source paths to be treated as different")
	}
}

func TestNormalizeSolutionVersion_PadsMissingSegments(t *testing.T) {
	t.Parallel()

	normalized, err := normalizeSolutionVersion("1.3.5")
	if err != nil {
		t.Fatalf("expected version normalization to succeed: %v", err)
	}

	if normalized != "1.3.5.0" {
		t.Fatalf("expected normalized version to be 1.3.5.0, got %s", normalized)
	}
}

func TestShouldPreserveStateSource_KeepsFreshTokenizedURL(t *testing.T) {
	t.Parallel()

	now := time.Unix(1_776_703_000, 0).UTC()
	source := &SourceModel{
		URL: types.StringValue(powerPackDownloadURL("CodeEditor", "1.0.3.1", now.Add(15*time.Minute))),
	}

	if !shouldPreserveStateSource(source, now) {
		t.Fatal("expected fresh PowerPack tokenized URL to be preserved from state")
	}
}

func TestShouldPreserveStateSource_DropsExpiredTokenizedURL(t *testing.T) {
	t.Parallel()

	now := time.Unix(1_776_703_000, 0).UTC()
	source := &SourceModel{
		URL: types.StringValue(powerPackDownloadURL("smartGridPCF", "1.0.45.1", now.Add(-time.Minute))),
	}

	if shouldPreserveStateSource(source, now) {
		t.Fatal("expected expired PowerPack tokenized URL to be refreshed from plan")
	}
}

func TestShouldPreserveStateSource_DropsNearlyExpiredTokenizedURL(t *testing.T) {
	t.Parallel()

	now := time.Unix(1_776_703_000, 0).UTC()
	source := &SourceModel{
		URL: types.StringValue(powerPackDownloadURL("EditableTable", "1.1.1.1", now.Add(4*time.Minute))),
	}

	if shouldPreserveStateSource(source, now) {
		t.Fatal("expected nearly expired PowerPack tokenized URL to be refreshed from plan")
	}
}

func powerPackDownloadURL(name string, version string, expiresAt time.Time) string {
	payload := fmt.Sprintf(
		`{"Package":%q,"Version":%q,"ExpiresAtUnixTimeSeconds":%d}`,
		name,
		version,
		expiresAt.Unix(),
	)
	token := fmt.Sprintf(
		"%s.%s",
		base64.RawURLEncoding.EncodeToString([]byte(payload)),
		base64.RawURLEncoding.EncodeToString([]byte("signature")),
	)

	return fmt.Sprintf(
		"https://func.example.net/api/packages/%s/%s/download?token=%s",
		name,
		version,
		token,
	)
}
