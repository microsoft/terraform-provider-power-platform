// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform_helpers

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (
	ERROR_OBJECT_NOT_FOUND ErrorCode = "OBJECT_NOT_FOUND"
)

type providerError struct {
	error
	errorCode ErrorCode
}

func (e providerError) Error() string {
	return fmt.Sprintf("%s: %s", e.errorCode, e.error.Error())
}

func Unwrap(err error) error {
	if e, ok := err.(providerError); ok {
		return errors.Unwrap(e.error)
	}

	return errors.Unwrap(err)
}

func Code(err error) ErrorCode {
	if err == nil {
		return ""
	}

	if e, ok := err.(providerError); ok {
		return e.errorCode
	}

	return ""
}

func NewProviderError(errorCode ErrorCode, format string, args ...interface{}) error {
	return providerError{
		error:     fmt.Errorf(format, args...),
		errorCode: errorCode,
	}
}

func WrapIntoProviderError(err error, errorCode ErrorCode, msg string) error {
	return providerError{
		error:     fmt.Errorf("%s: [%w]", msg, err),
		errorCode: errorCode,
	}
}
