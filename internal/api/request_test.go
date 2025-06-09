package api

import (
	"context"
	"testing"

	"github.com/microsoft/terraform-provider-power-platform/internal/config"
	"github.com/stretchr/testify/require"
)

func TestUnitBuildUserAgent_WithPartnerId(t *testing.T) {
	cfg := config.ProviderConfig{PartnerId: "00000000-0000-0000-0000-000000000001"}
	client := NewApiClientBase(&cfg, NewAuthBase(&cfg))
	ua := client.buildUserAgent(context.Background())
	require.Contains(t, ua, "pid-00000000-0000-0000-0000-000000000001")
}
