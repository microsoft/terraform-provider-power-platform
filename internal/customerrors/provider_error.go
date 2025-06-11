// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package customerrors

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (
	ERROR_OBJECT_NOT_FOUND             ErrorCode = "OBJECT_NOT_FOUND"
	ERROR_ENVIRONMENT_URL_NOT_FOUND    ErrorCode = "ENVIRONMENT_URL_NOT_FOUND"
	ERROR_ENVIRONMENTS_IN_ENV_GROUP    ErrorCode = "ENVIRONMENTS_IN_ENV_GROUP"
	ERROR_POLICY_ASSIGNED_TO_ENV_GROUP ErrorCode = "POLICY_ASSIGNED_TO_ENV_GROUP"
	ERROR_ENVIRONMENT_SETTINGS_FAILED  ErrorCode = "ENVIRONMENT_SETTINGS_FAILED"
	ERROR_ENVIRONMENT_CREATION         ErrorCode = "ENVIRONMENT_CREATION"
)

// Sentinel errors for use with errors.Is().
var (
	ErrObjectNotFound            = ProviderError{ErrorCode: ERROR_OBJECT_NOT_FOUND}
	ErrEnvironmentUrlNotFound    = ProviderError{ErrorCode: ERROR_ENVIRONMENT_URL_NOT_FOUND}
	ErrEnvironmentsInEnvGroup    = ProviderError{ErrorCode: ERROR_ENVIRONMENTS_IN_ENV_GROUP}
	ErrPolicyAssignedToEnvGroup  = ProviderError{ErrorCode: ERROR_POLICY_ASSIGNED_TO_ENV_GROUP}
	ErrEnvironmentSettingsFailed = ProviderError{ErrorCode: ERROR_ENVIRONMENT_SETTINGS_FAILED}
	ErrEnvironmentCreation       = ProviderError{ErrorCode: ERROR_ENVIRONMENT_CREATION}
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

// Is implements the Is method for error equality checking with errors.Is().
func (e ProviderError) Is(target error) bool {
	if t, ok := target.(ProviderError); ok {
		return e.ErrorCode == t.ErrorCode
	}
	return false
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
