// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package customerrors

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (
	ERROR_OBJECT_NOT_FOUND          ErrorCode = "OBJECT_NOT_FOUND"
	ERROR_ENVIRONMENT_URL_NOT_FOUND ErrorCode = "ENVIRONMENT_URL_NOT_FOUND"
)

var _ error = ProviderError{}

type ProviderError struct {
	ErrorCode ErrorCode
	Err       error
}

func (e ProviderError) Error() string {
	if e.Err == nil {
		return string(e.ErrorCode)
	}

	return fmt.Sprintf("%s: %s", e.ErrorCode, e.Err.Error())
}

func Unwrap(err error) error {
	if e, ok := err.(ProviderError); ok {
		return errors.Unwrap(e.Err)
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
		Err:       fmt.Errorf(format, args...),
		ErrorCode: errorCode,
	}
}

func WrapIntoProviderError(err error, errorCode ErrorCode, msg string) error {
	if err == nil {
		return ProviderError{
			Err:       fmt.Errorf("%s", msg),
			ErrorCode: errorCode,
		}
	}
	return ProviderError{
		Err:       fmt.Errorf("%s: [%w]", msg, err),
		ErrorCode: errorCode,
	}
}
