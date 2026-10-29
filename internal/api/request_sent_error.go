// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api

// RequestSentError marks an error raised at or after the point the HTTP request may have
// reached the service: the transport failed, the response went missing, or an accepted
// response could not be read or decoded. A caller performing a non-idempotent operation can
// use it to separate these unknown-outcome failures from local errors that provably sent
// nothing, like scope resolution, url validation, token acquisition or request construction.
// It adds no text of its own, so wrapped errors read exactly as they did before.
type RequestSentError struct {
	Err error
}

func (e RequestSentError) Error() string {
	return e.Err.Error()
}

func (e RequestSentError) Unwrap() error {
	return e.Err
}
