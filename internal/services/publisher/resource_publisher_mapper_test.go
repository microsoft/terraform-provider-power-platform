// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package publisher

import (
	"math"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
)

func TestUnitAddressModelFromDto_IgnoresEmptySlotWithOnlyAddressID(t *testing.T) {
	dto := &publisherDto{
		Address2AddressId: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
	}

	model := addressModelFromDto(2, dto)
	if model != nil {
		t.Fatalf("expected address slot 2 to be ignored when only address id remains, got %#v", model)
	}
}

func TestUnitAddressModelsFromDto_IgnoresPlaceholderSlotWithoutExistingState(t *testing.T) {
	dto := &publisherDto{
		Address1AddressId:          "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		Address1AddressTypeCode:    int64Pointer(1),
		Address1ShippingMethodCode: int64Pointer(1),
	}

	models := addressModelsFromDto(dto, nil)
	if models != nil {
		t.Fatalf("expected placeholder address slot to be ignored, got %#v", models)
	}
}

func TestUnitAddressModelsFromDto_PreservesPlaceholderSlotWhenAlreadyTracked(t *testing.T) {
	dto := &publisherDto{
		Address1AddressId:          "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		Address1AddressTypeCode:    int64Pointer(1),
		Address1ShippingMethodCode: int64Pointer(1),
	}

	existing := []PublisherAddressModel{
		{
			Slot:               types.Int64Value(1),
			AddressTypeCode:    types.Int64Value(1),
			ShippingMethodCode: types.Int64Value(1),
		},
	}

	models := addressModelsFromDto(dto, existing)
	if len(models) != 1 {
		t.Fatalf("expected tracked placeholder address slot to be preserved, got %#v", models)
	}
	if models[0].Slot.ValueInt64() != 1 {
		t.Fatalf("expected preserved address slot 1, got %#v", models[0])
	}
}

func TestUnitSetResourceModelFromDto_PreservesExplicitEmptyTopLevelStrings(t *testing.T) {
	model := ResourceModel{
		Description:          types.StringValue(""),
		EmailAddress:         types.StringValue(""),
		SupportingWebsiteURL: types.StringValue(""),
	}

	setResourceModelFromDto(&model, "00000000-0000-0000-0000-000000000001", &publisherDto{
		Id:                             "11111111-1111-1111-1111-111111111111",
		UniqueName:                     "testpublisher",
		FriendlyName:                   "Test Publisher",
		CustomizationPrefix:            "tp",
		CustomizationOptionValuePrefix: 12345,
	})

	if model.Description.IsNull() || model.Description.ValueString() != "" {
		t.Fatalf("expected empty description to be preserved, got %#v", model.Description)
	}
	if model.EmailAddress.IsNull() || model.EmailAddress.ValueString() != "" {
		t.Fatalf("expected empty email_address to be preserved, got %#v", model.EmailAddress)
	}
	if model.SupportingWebsiteURL.IsNull() || model.SupportingWebsiteURL.ValueString() != "" {
		t.Fatalf("expected empty supporting_website_url to be preserved, got %#v", model.SupportingWebsiteURL)
	}
}

func TestUnitAddressModelsFromDto_PreservesExplicitEmptyAddressList(t *testing.T) {
	models := addressModelsFromDto(&publisherDto{}, []PublisherAddressModel{})
	if models == nil {
		t.Fatal("expected explicit empty address list to be preserved")
	}
	if len(models) != 0 {
		t.Fatalf("expected no address entries, got %#v", models)
	}
}

func TestUnitAddressModelsFromDto_PreservesTrackedEmptyStringAddressFields(t *testing.T) {
	dto := &publisherDto{
		Address1AddressId:          "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		Address1AddressTypeCode:    int64Pointer(1),
		Address1ShippingMethodCode: int64Pointer(1),
	}

	existing := []PublisherAddressModel{
		{
			Slot:               types.Int64Value(1),
			AddressTypeCode:    types.Int64Value(1),
			ShippingMethodCode: types.Int64Value(1),
			Line1:              types.StringValue(""),
		},
	}

	models := addressModelsFromDto(dto, existing)
	if len(models) != 1 {
		t.Fatalf("expected tracked address slot to be preserved, got %#v", models)
	}
	if models[0].Line1.IsNull() || models[0].Line1.ValueString() != "" {
		t.Fatalf("expected empty line1 to be preserved, got %#v", models[0].Line1)
	}
}

func TestUnitDeriveCustomizationOptionValuePrefix_MatchesClientAlgorithm(t *testing.T) {
	got := deriveCustomizationOptionValuePrefix("mf", "")
	if got != 12457 {
		t.Fatalf("expected derived prefix 12457 for 'mf', got %d", got)
	}
}

func TestUnitDeriveCustomizationOptionValuePrefix_UsesPublisherSpecialCase(t *testing.T) {
	got := deriveCustomizationOptionValuePrefix("anything", "d21aab71-79e7-11dd-8874-00188b01e34f")
	if got != 10000 {
		t.Fatalf("expected special-case derived prefix 10000, got %d", got)
	}
}

func TestUnitDeriveCustomizationOptionValuePrefix_UsesPublisherSpecialCaseCaseInsensitive(t *testing.T) {
	got := deriveCustomizationOptionValuePrefix("anything", "D21AAB71-79E7-11DD-8874-00188B01E34F")
	if got != 10000 {
		t.Fatalf("expected case-insensitive special-case derived prefix 10000, got %d", got)
	}
}

func TestUnitCustomizationOptionValuePrefixFromHash_HandlesMinInt32(t *testing.T) {
	got := customizationOptionValuePrefixFromHash(math.MinInt32)
	if got != 93648 {
		t.Fatalf("expected min-int32 hash to produce 93648, got %d", got)
	}
}

func TestUnitSetDerivedCustomizationOptionValuePrefix_DerivesWhenConfigOmitted(t *testing.T) {
	plan := ResourceModel{
		CustomizationPrefix:            types.StringValue("mf"),
		CustomizationOptionValuePrefix: types.Int64Unknown(),
	}
	config := ResourceModel{
		CustomizationOptionValuePrefix: types.Int64Null(),
	}

	setDerivedCustomizationOptionValuePrefix(&plan, &config, &ResourceModel{}, false)

	if plan.CustomizationOptionValuePrefix.IsUnknown() || plan.CustomizationOptionValuePrefix.IsNull() {
		t.Fatal("expected derived customization option value prefix to be planned")
	}
	if plan.CustomizationOptionValuePrefix.ValueInt64() != 12457 {
		t.Fatalf("expected derived customization option value prefix 12457, got %d", plan.CustomizationOptionValuePrefix.ValueInt64())
	}
}

func TestUnitSetDerivedCustomizationOptionValuePrefix_PreservesExplicitConfigValue(t *testing.T) {
	plan := ResourceModel{
		CustomizationPrefix:            types.StringValue("mf"),
		CustomizationOptionValuePrefix: types.Int64Value(77777),
	}
	config := ResourceModel{
		CustomizationOptionValuePrefix: types.Int64Value(77777),
	}

	setDerivedCustomizationOptionValuePrefix(&plan, &config, &ResourceModel{}, false)

	if plan.CustomizationOptionValuePrefix.ValueInt64() != 77777 {
		t.Fatalf("expected explicit customization option value prefix to be preserved, got %d", plan.CustomizationOptionValuePrefix.ValueInt64())
	}
}

func TestUnitSetDerivedCustomizationOptionValuePrefix_PreservesStateValueAfterCreate(t *testing.T) {
	plan := ResourceModel{
		Id:                             types.StringValue("11111111-1111-1111-1111-111111111111"),
		CustomizationPrefix:            types.StringValue("ab"),
		CustomizationOptionValuePrefix: types.Int64Unknown(),
	}
	config := ResourceModel{
		CustomizationOptionValuePrefix: types.Int64Null(),
	}
	state := ResourceModel{
		Id:                             types.StringValue("11111111-1111-1111-1111-111111111111"),
		CustomizationPrefix:            types.StringValue("old"),
		CustomizationOptionValuePrefix: types.Int64Value(77074),
	}

	setDerivedCustomizationOptionValuePrefix(&plan, &config, &state, true)

	if plan.CustomizationOptionValuePrefix.ValueInt64() != state.CustomizationOptionValuePrefix.ValueInt64() {
		t.Fatalf("expected existing customization option value prefix to be preserved from state, got %d", plan.CustomizationOptionValuePrefix.ValueInt64())
	}
}

func TestUnitGetPublisherIdFromResponse_ParsesCanonicalGuid(t *testing.T) {
	resp := &api.Response{
		HttpResponse: &http.Response{
			Header: http.Header{
				"OData-EntityId": []string{"https://example.crm.dynamics.com/api/data/v9.2/publishers(11111111-1111-1111-1111-111111111111)"},
			},
		},
	}

	got, err := getPublisherIdFromResponse(resp)
	if err != nil {
		t.Fatalf("expected publisher id to be parsed, got error: %v", err)
	}
	if got != "11111111-1111-1111-1111-111111111111" {
		t.Fatalf("expected canonical publisher id, got %q", got)
	}
}

func TestUnitGetPublisherIdFromResponse_RejectsNonCanonicalGuid(t *testing.T) {
	resp := &api.Response{
		HttpResponse: &http.Response{
			Header: http.Header{
				"OData-EntityId": []string{"https://example.crm.dynamics.com/api/data/v9.2/publishers(11111111-1111-1111-1111-11111111111)"},
			},
		},
	}

	_, err := getPublisherIdFromResponse(resp)
	if err == nil {
		t.Fatal("expected malformed publisher id to be rejected")
	}
}

func TestUnitIsGuid_RequiresCanonicalGuidFormat(t *testing.T) {
	if !isGuid("11111111-1111-1111-1111-111111111111") {
		t.Fatal("expected canonical guid to be recognized")
	}
	if isGuid("11111111-1111-1111-1111-11111111111") {
		t.Fatal("expected malformed guid to be rejected")
	}
}

func int64Pointer(value int64) *int64 {
	return &value
}

func TestUnitIsPublisherAlreadyExistsError_MatchesDataverseDuplicateKeyResponse(t *testing.T) {
	err := customerrors.NewUnexpectedHttpStatusCodeError(
		[]int{201, 204},
		412,
		"412 Precondition Failed",
		[]byte(`{"error":{"code":"0x80040237","message":"A record with matching key values already exists."}}`),
	)

	if !isPublisherAlreadyExistsError(err) {
		t.Fatal("expected Dataverse duplicate-key response to be recognized as already exists")
	}
}

func TestUnitIsPublisherAlreadyExistsError_IgnoresOtherUnexpectedStatusErrors(t *testing.T) {
	err := customerrors.NewUnexpectedHttpStatusCodeError(
		[]int{201, 204},
		409,
		"409 Conflict",
		[]byte(`{"error":{"code":"other","message":"conflict"}}`),
	)

	if isPublisherAlreadyExistsError(err) {
		t.Fatal("expected non-412 error to be ignored")
	}
}
