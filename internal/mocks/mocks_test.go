// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package mocks_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

//go:noinline
func testNameCaller() string {
	return mocks.TestName()
}

func TestUnitTestName(t *testing.T) {
	if got := testNameCaller(); got != "testNameCaller" {
		t.Fatalf("expected %q, got %q", "testNameCaller", got)
	}
}

//go:noinline
func aTestNameComfortablyLongerThanTheSixtyFourCharacterEntraNicknameLimit() string {
	return mocks.TestNameForEntra()
}

//go:noinline
func aTestNameComfortablyLongerThanTheSixtyFourCharacterEntraNicknameLimitToo() string {
	return mocks.TestNameForEntra()
}

func TestUnitTestNameForEntra(t *testing.T) {
	short := testNameForEntraCaller()
	if short != "testNameForEntraCaller" {
		t.Fatalf("a name within the limit must be returned unchanged, got %q", short)
	}

	long := aTestNameComfortablyLongerThanTheSixtyFourCharacterEntraNicknameLimit()
	if len(long) > mocks.ENTRA_MAIL_NICKNAME_MAX_LENGTH {
		t.Fatalf("expected at most %d characters, got %d (%q)", mocks.ENTRA_MAIL_NICKNAME_MAX_LENGTH, len(long), long)
	}

	// Names sharing the first 64 characters must not collapse onto the same nickname.
	if other := aTestNameComfortablyLongerThanTheSixtyFourCharacterEntraNicknameLimitToo(); other == long {
		t.Fatalf("two distinct test names produced the same nickname %q", long)
	}
}

//go:noinline
func testNameForEntraCaller() string {
	return mocks.TestNameForEntra()
}

func TestUnitTestsEntraLicesingGroupName(t *testing.T) {
	if got := mocks.TestsEntraLicesingGroupName(); got == "" {
		t.Fatal("expected group name to be non-empty")
	}
}

func TestUnitProviderFactories(t *testing.T) {
	if len(mocks.TestUnitTestProtoV6ProviderFactories) == 0 {
		t.Fatal("expected unit test provider factories to be registered")
	}
	if len(mocks.TestAccProtoV6ProviderFactories) == 0 {
		t.Fatal("expected acceptance provider factories to be registered")
	}
	for name, factory := range mocks.TestUnitTestProtoV6ProviderFactories {
		if name == "" || factory == nil {
			t.Fatal("expected factory entry to be valid")
		}
	}
	for name, factory := range mocks.TestAccProtoV6ProviderFactories {
		if name == "" || factory == nil {
			t.Fatal("expected factory entry to be valid")
		}
	}
}

func TestUnitActivateEnvironmentHttpMocks(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	getURLs := []string{
		"https://org000001.crm4.dynamics.com/api/data/v9.2/transactioncurrencies",
		"https://org000001.crm4.dynamics.com/api/data/v9.2/organizations",
		"https://org000001.crm.dynamics.com/api/data/v9.2/WhoAmI",
		"https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/tenant?api-version=2021-04-01",
	}

	for _, url := range getURLs {
		resp, err := http.Get(url)
		if err != nil {
			t.Fatalf("expected GET to succeed for %s, got error: %v", url, err)
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200 for %s, got %d", url, resp.StatusCode)
		}
	}

	resp, err := http.Post("https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/validateEnvironmentDetails?api-version=2021-04-01", "application/json", nil)
	if err != nil {
		t.Fatalf("expected POST to succeed, got error: %v", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 for validateEnvironmentDetails, got %d", resp.StatusCode)
	}
}
