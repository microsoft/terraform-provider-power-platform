// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package customerrors

import (
	"fmt"
)

var _ error = (*UrlFormatError)(nil)

type UrlFormatError struct {
	Url string
	Err error
}

func NewUrlFormatError(url string, err error) error {
	return &UrlFormatError{
		Err: err,
		Url: url,
	}
}

func (e *UrlFormatError) Error() string {
	errorMsg := ""
	if e.Err != nil {
		errorMsg = e.Err.Error()
	}

	return fmt.Sprintf("Request url must be an absolute url: '%s' : '%s'", e.Url, errorMsg)
}

func (e *UrlFormatError) Unwrap() error {
	return e.Err
}
