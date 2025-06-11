// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package customerrors

import (
	"errors"
	"fmt"

	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
)

type ErrorCode string

// Sentinel errors for use with errors.Is().
var (
	ErrObjectNotFound            = ProviderError{ErrorCode: ErrorCode(constants.ERROR_OBJECT_NOT_FOUND)}
	ErrEnvironmentUrlNotFound    = ProviderError{ErrorCode: ErrorCode(constants.ERROR_ENVIRONMENT_URL_NOT_FOUND)}
	ErrEnvironmentsInEnvGroup    = ProviderError{ErrorCode: ErrorCode(constants.ERROR_ENVIRONMENTS_IN_ENV_GROUP)}
	ErrPolicyAssignedToEnvGroup  = ProviderError{ErrorCode: ErrorCode(constants.ERROR_POLICY_ASSIGNED_TO_ENV_GROUP)}
	ErrEnvironmentSettingsFailed = ProviderError{ErrorCode: ErrorCode(constants.ERROR_ENVIRONMENT_SETTINGS_FAILED)}
	ErrEnvironmentCreation       = ProviderError{ErrorCode: ErrorCode(constants.ERROR_ENVIRONMENT_CREATION)}
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
