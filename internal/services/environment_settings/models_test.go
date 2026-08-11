// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_settings

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var privacyAndSecurityTestAttrTypes = map[string]attr.Type{
	"blocked_attachment_extensions": types.SetType{ElemType: types.StringType},
}

func TestUnitConvertFromEnvironmentPrivacyAndSecuritySettings_OmittedBlock_OmitsBlockedAttachments(t *testing.T) {
	model := EnvironmentSettingsResourceModel{
		PrivacyAndSecurity: types.ObjectNull(privacyAndSecurityTestAttrTypes),
	}
	dto := &environmentOrgSettingsDto{}

	if err := convertFromEnvironmentPrivacyAndSecuritySettings(context.Background(), model, dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dto.BlockedAttachments != nil {
		t.Errorf("expected BlockedAttachments to be nil for an omitted privacy_and_security block, got %q", *dto.BlockedAttachments)
	}
	if blockedAttachmentsInPatchBody(t, dto) {
		t.Error("expected blockedattachments to be omitted from the PATCH body")
	}
}

func TestUnitConvertFromEnvironmentPrivacyAndSecuritySettings_UnknownBlock_OmitsBlockedAttachments(t *testing.T) {
	model := EnvironmentSettingsResourceModel{
		PrivacyAndSecurity: types.ObjectUnknown(privacyAndSecurityTestAttrTypes),
	}
	dto := &environmentOrgSettingsDto{}

	if err := convertFromEnvironmentPrivacyAndSecuritySettings(context.Background(), model, dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dto.BlockedAttachments != nil {
		t.Errorf("expected BlockedAttachments to be nil for an unknown privacy_and_security block, got %q", *dto.BlockedAttachments)
	}
	if blockedAttachmentsInPatchBody(t, dto) {
		t.Error("expected blockedattachments to be omitted from the PATCH body")
	}
}

func TestUnitConvertFromEnvironmentPrivacyAndSecuritySettings_OmittedAttribute_OmitsBlockedAttachments(t *testing.T) {
	model := EnvironmentSettingsResourceModel{
		PrivacyAndSecurity: types.ObjectValueMust(privacyAndSecurityTestAttrTypes, map[string]attr.Value{
			"blocked_attachment_extensions": types.SetNull(types.StringType),
		}),
	}
	dto := &environmentOrgSettingsDto{}

	if err := convertFromEnvironmentPrivacyAndSecuritySettings(context.Background(), model, dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dto.BlockedAttachments != nil {
		t.Errorf("expected BlockedAttachments to be nil for an omitted blocked_attachment_extensions attribute, got %q", *dto.BlockedAttachments)
	}
	if blockedAttachmentsInPatchBody(t, dto) {
		t.Error("expected blockedattachments to be omitted from the PATCH body")
	}
}

func TestUnitConvertFromEnvironmentPrivacyAndSecuritySettings_EmptySet_ClearsBlockedAttachments(t *testing.T) {
	model := EnvironmentSettingsResourceModel{
		PrivacyAndSecurity: types.ObjectValueMust(privacyAndSecurityTestAttrTypes, map[string]attr.Value{
			"blocked_attachment_extensions": types.SetValueMust(types.StringType, []attr.Value{}),
		}),
	}
	dto := &environmentOrgSettingsDto{}

	if err := convertFromEnvironmentPrivacyAndSecuritySettings(context.Background(), model, dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dto.BlockedAttachments == nil {
		t.Fatal("expected BlockedAttachments to be set for an explicitly empty set")
	}
	if *dto.BlockedAttachments != "" {
		t.Errorf("expected BlockedAttachments to be an empty string, got %q", *dto.BlockedAttachments)
	}
	if !blockedAttachmentsInPatchBody(t, dto) {
		t.Error("expected blockedattachments to be present in the PATCH body")
	}
}

func TestUnitConvertFromEnvironmentPrivacyAndSecuritySettings_PopulatedSet_JoinsWithSemicolons(t *testing.T) {
	model := EnvironmentSettingsResourceModel{
		PrivacyAndSecurity: types.ObjectValueMust(privacyAndSecurityTestAttrTypes, map[string]attr.Value{
			"blocked_attachment_extensions": types.SetValueMust(types.StringType, []attr.Value{
				types.StringValue("exe"),
				types.StringValue("js"),
				types.StringValue("dll"),
			}),
		}),
	}
	dto := &environmentOrgSettingsDto{}

	if err := convertFromEnvironmentPrivacyAndSecuritySettings(context.Background(), model, dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dto.BlockedAttachments == nil {
		t.Fatal("expected BlockedAttachments to be set for a populated set")
	}
	if *dto.BlockedAttachments != "exe;js;dll" {
		t.Errorf("expected BlockedAttachments to be %q, got %q", "exe;js;dll", *dto.BlockedAttachments)
	}
	if !blockedAttachmentsInPatchBody(t, dto) {
		t.Error("expected blockedattachments to be present in the PATCH body")
	}
}

// blockedAttachmentsInPatchBody reports whether the serialized PATCH payload
// contains the blockedattachments field at all.
func blockedAttachmentsInPatchBody(t *testing.T, dto *environmentOrgSettingsDto) bool {
	t.Helper()

	body, err := json.Marshal(dto)
	if err != nil {
		t.Fatalf("failed to marshal dto: %v", err)
	}

	return strings.Contains(string(body), `"blockedattachments"`)
}

func TestUnitConvertBlockedAttachmentsToSetValues(t *testing.T) {
	stringPtr := func(s string) *string { return &s }

	tests := []struct {
		name     string
		input    *string
		expected []string
	}{
		{
			name:     "nil input returns empty slice",
			input:    nil,
			expected: []string{},
		},
		{
			name:     "empty string returns empty slice",
			input:    stringPtr(""),
			expected: []string{},
		},
		{
			name:     "clean list",
			input:    stringPtr("exe;js;dll"),
			expected: []string{"exe", "js", "dll"},
		},
		{
			name:     "messy list with whitespace, trailing and duplicate semicolons",
			input:    stringPtr(" exe; dll ;;js ; "),
			expected: []string{"exe", "dll", "js"},
		},
		{
			name:     "only separators and whitespace returns empty slice",
			input:    stringPtr(" ; ;; "),
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := convertBlockedAttachmentsToSetValues(tt.input)

			if len(values) != len(tt.expected) {
				t.Fatalf("expected %d values %v, got %d values %v", len(tt.expected), tt.expected, len(values), values)
			}
			for i, expected := range tt.expected {
				str, ok := values[i].(types.String)
				if !ok {
					t.Fatalf("expected values[%d] to be types.String, got %T", i, values[i])
				}
				if str.ValueString() != expected {
					t.Errorf("expected values[%d] to be %q, got %q", i, expected, str.ValueString())
				}
			}
		})
	}
}
