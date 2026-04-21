package managedsolution

import (
	"testing"

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
