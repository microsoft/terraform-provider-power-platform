package api

import (
	"context"
	"net/http"
	"net/http/httptest"
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

func TestUnitDoRequest_XMsUserAgent_TelemetryNotOptedOut(t *testing.T) {
	var capturedRequest *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedRequest = r
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := config.ProviderConfig{TelemetryOptout: false}
	client := NewApiClientBase(&cfg, NewAuthBase(&cfg))

	token := "test-token"
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	_, err = client.doRequest(context.Background(), &token, req, http.Header{})
	require.NoError(t, err)

	require.Equal(t, "terraform-provider-power-platform", capturedRequest.Header.Get("x-ms-useragent"))
}

func TestUnitDoRequest_XMsUserAgent_TelemetryOptedOut(t *testing.T) {
	var capturedRequest *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedRequest = r
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := config.ProviderConfig{TelemetryOptout: true}
	client := NewApiClientBase(&cfg, NewAuthBase(&cfg))

	token := "test-token"
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	_, err = client.doRequest(context.Background(), &token, req, http.Header{})
	require.NoError(t, err)

	require.Empty(t, capturedRequest.Header.Get("x-ms-useragent"))
}

func TestUnitResponse_MarshallTo(t *testing.T) {
	tests := []struct {
		name        string
		response    *Response
		obj         any
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid pointer to struct",
			response: &Response{
				BodyAsBytes: []byte(`{"name": "test", "value": 123}`),
			},
			obj: &struct {
				Name  string
				Value int
			}{},
			expectError: false,
		},
		{
			name: "Non-pointer value should fail",
			response: &Response{
				BodyAsBytes: []byte(`{"name": "test"}`),
			},
			obj:         struct{ Name string }{},
			expectError: true,
			errorMsg:    "MarshallTo requires a non-nil pointer",
		},
		{
			name: "Nil pointer should fail",
			response: &Response{
				BodyAsBytes: []byte(`{"name": "test"}`),
			},
			obj:         (*struct{ Name string })(nil),
			expectError: true,
			errorMsg:    "MarshallTo requires a non-nil pointer",
		},
		{
			name: "Non-pointer interface should fail",
			response: &Response{
				BodyAsBytes: []byte(`123`),
			},
			obj:         123,
			expectError: true,
			errorMsg:    "MarshallTo requires a non-nil pointer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.response.MarshallTo(tt.obj)

			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
