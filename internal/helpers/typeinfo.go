// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers

import (
	"fmt"
)

// TypeInfo represents a managed object type in the provider such as a resource or data source.
// Resource and data source types can inherit from TypeInfo to provide a consistent way to reference the type.
type TypeInfo struct {
	ProviderTypeName string
	TypeName         string
}

// FullTypeName returns the full type name in the format provider_type
func (t *TypeInfo) FullTypeName() string {
	return fmt.Sprintf("%s_%s", t.ProviderTypeName, t.TypeName)
}
