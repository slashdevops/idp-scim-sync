package aws

import "fmt"

// HTTPResponseError describes a non-2xx response from the AWS SSO SCIM API.
//
// It is part of this package's error contract: callers use errors.AsType to
// recover it and branch on StatusCode. pkg/aws does so for 409 Conflict in
// CreateOrGetUser and CreateOrGetGroup and for 404 Not Found in DeleteUser and
// DeleteGroup, and internal/core does so when classifying state-repository
// failures.
type HTTPResponseError struct {
	// Code is the SCIM scimType value when the response carried a structured
	// SCIM error body, and the raw HTTP status line otherwise.
	Code string `json:"ErrorCode"`

	// Message is the SCIM detail value when present, and the raw response body
	// otherwise.
	Message string `json:"ErrorMessage"`

	// StatusCode is the HTTP status code of the response.
	StatusCode int `json:"StatusCode"`
}

// Error implements the error interface.
func (e *HTTPResponseError) Error() string {
	return fmt.Sprintf("statusCode: %d,  errCode: %s, errMsg: %s", e.StatusCode, e.Code, e.Message)
}
