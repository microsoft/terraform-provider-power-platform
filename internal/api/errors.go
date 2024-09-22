// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api

import (
	"fmt"
)

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

type UrlFormatError struct {
	error
	Url string
}

func NewUrlFormatError(url string, err error) error {
	return UrlFormatError{
		error: err,
		Url:   url,
	}
}

func (e UrlFormatError) Error() string {
	errorMsg := ""
	if e.error != nil {
		errorMsg = e.error.Error()
	}

	return fmt.Sprintf("Request url must be an absolute url. %s : %s", e.Url, errorMsg)
}
