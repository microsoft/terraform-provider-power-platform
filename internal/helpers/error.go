// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (
	ERROR_OBJECT_NOT_FOUND            ErrorCode = "OBJECT_NOT_FOUND"
	ERROR_UNEXPECTED_HTTP_RETURN_CODE ErrorCode = "UNEXPECTED_HTTP_RETURN_CODE"
	ERROR_INCORRECT_URL_FORMAT        ErrorCode = "INCORRECT_URL_FORMAT"
	ERROR_ENVIRONMENT_URL_NOT_FOUND   ErrorCode = "ENVIRONMENT_URL_NOT_FOUND"
)

type ProviderError struct {
	error
	ErrorCode ErrorCode
}

func (e ProviderError) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode, e.error.Error())
}

func Unwrap(err error) error {
	if e, ok := err.(ProviderError); ok {
		return errors.Unwrap(e.error)
	}

	return errors.Unwrap(err)
}

func Code(err error) ErrorCode {
	if err == nil {
		return ""
	}

	if e, ok := err.(ProviderError); ok {
		return e.ErrorCode
	}

	return ""
}

func NewProviderError(errorCode ErrorCode, format string, args ...any) error {
	return ProviderError{
		error:     fmt.Errorf(format, args...),
		ErrorCode: errorCode,
	}
}

func WrapIntoProviderError(err error, errorCode ErrorCode, msg string) error {
	if err == nil {
		return ProviderError{
			error:     fmt.Errorf("%s", msg),
			ErrorCode: errorCode,
		}
	}
	return ProviderError{
		error:     fmt.Errorf("%s: [%w]", msg, err),
		ErrorCode: errorCode,
	}
}
