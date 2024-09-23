// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package customerrors

import (
	"fmt"
)

var _ error = UnexpectedHttpStatusCodeError{}

type UnexpectedHttpStatusCodeError struct {
	ExpectedStatusCodes []int
	StatusCode          int
	StatusText          string
	Body                []byte
}

func (e UnexpectedHttpStatusCodeError) Error() string {
	return fmt.Sprintf("Unexpected HTTP status code. Expected: %v, received: [%d] %s | %s", e.ExpectedStatusCodes, e.StatusCode, e.StatusText, e.Body)
}

func NewUnexpectedHttpStatusCodeError(expectedStatusCodes []int, statusCode int, statusText string, body []byte) error {
	return UnexpectedHttpStatusCodeError{
		ExpectedStatusCodes: expectedStatusCodes,
		StatusCode:          statusCode,
		StatusText:          statusText,
		Body:                body,
	}
}
