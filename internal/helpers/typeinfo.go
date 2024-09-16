// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers

import (
	"fmt"
)

type TypeInfo struct {
	ProviderTypeName string
	TypeName         string
}

func (t *TypeInfo) FullTypeName() string {
	return fmt.Sprintf("%s_%s", t.ProviderTypeName, t.TypeName)
}
