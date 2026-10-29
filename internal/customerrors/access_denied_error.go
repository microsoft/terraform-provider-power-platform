// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package customerrors

import (
	"fmt"
	"strings"
)

var _ error = AccessDeniedError{}

// AccessDeniedError reports an HTTP 403 from a service. The body is retained because it carries the
// only explanation of which principal or privilege was rejected.
type AccessDeniedError struct {
	Url  string
	Body []byte
}

func (e AccessDeniedError) Error() string {
	if body := strings.TrimSpace(string(e.Body)); body != "" {
		return fmt.Sprintf("access denied to resource at '%s'. Please validate your permissions: %s", e.Url, body)
	}
	return fmt.Sprintf("access denied to resource at '%s'. Please validate your permissions", e.Url)
}

func NewAccessDeniedError(url string, body []byte) error {
	return AccessDeniedError{
		Url:  url,
		Body: body,
	}
}
