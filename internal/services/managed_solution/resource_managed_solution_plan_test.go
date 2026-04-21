package managedsolution

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSourcesAreEquivalent_IgnoresTokenQueryParameter(t *testing.T) {
	t.Parallel()

	plan := &SourceModel{
		URL: types.StringValue("https://func.example.net/api/packages/CodeEditor/1.0.3.1/download?token=new-token"),
	}
	state := &SourceModel{
		URL: types.StringValue("https://func.example.net/api/packages/CodeEditor/1.0.3.1/download?token=old-token"),
	}

	if !sourcesAreEquivalent(plan, state) {
		t.Fatal("expected URLs that differ only by token query value to be treated as equivalent")
	}
}

func TestSourcesAreEquivalent_DetectsTransportChange(t *testing.T) {
	t.Parallel()

	plan := &SourceModel{
		URL: types.StringValue("https://func.example.net/api/packages/CodeEditor/1.0.3.1/download?token=new-token"),
	}
	state := &SourceModel{
		URL: types.StringValue("http://localhost:45265/SCCHousing/362c42e5-845c-4901-aad3-bbcbc0b8e053/project/PowerPlatform/codeeditor/1.0.3.1"),
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
