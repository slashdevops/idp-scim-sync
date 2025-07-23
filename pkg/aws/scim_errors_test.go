package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPResponseError(t *testing.T) {
	t.Run("should format error message correctly", func(t *testing.T) {
		err := &HTTPResponseError{
			StatusCode: 400,
			Code:       "invalidValue",
			Message:    "Invalid filter syntax",
		}

		expected := "statusCode: 400,  errCode: invalidValue, errMsg: Invalid filter syntax"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("should handle empty code and message", func(t *testing.T) {
		err := &HTTPResponseError{
			StatusCode: 500,
			Code:       "",
			Message:    "",
		}

		expected := "statusCode: 500,  errCode: , errMsg: "
		assert.Equal(t, expected, err.Error())
	})

	t.Run("should handle different status codes", func(t *testing.T) {
		testCases := []struct {
			statusCode int
			code       string
			message    string
			expected   string
		}{
			{404, "notFound", "Resource not found", "statusCode: 404,  errCode: notFound, errMsg: Resource not found"},
			{429, "tooManyRequests", "Rate limit exceeded", "statusCode: 429,  errCode: tooManyRequests, errMsg: Rate limit exceeded"},
			{401, "unauthorized", "Invalid token", "statusCode: 401,  errCode: unauthorized, errMsg: Invalid token"},
		}

		for _, tc := range testCases {
			err := &HTTPResponseError{
				StatusCode: tc.statusCode,
				Code:       tc.code,
				Message:    tc.message,
			}
			assert.Equal(t, tc.expected, err.Error())
		}
	})
}

func TestErrorConstants(t *testing.T) {
	t.Run("should have consistent error prefixes", func(t *testing.T) {
		errorVars := []error{
			// SCIM service errors
			ErrURLEmpty,
			ErrCreateGroupRequestEmpty,
			ErrCreateUserRequestEmpty,
			ErrPatchGroupRequestEmpty,
			ErrGroupIDEmpty,
			ErrPatchUserRequestEmpty,
			ErrPutUserRequestEmpty,
			ErrUserExternalIDEmpty,
			ErrGroupDisplayNameEmpty,
			ErrGroupExternalIDEmpty,
			ErrBearerTokenEmpty,
			ErrServiceProviderConfigEmpty,
			// Model validation errors
			ErrUserIDEmpty,
			ErrEmailsTooMany,
			ErrEmailsEmpty,
			ErrFamilyNameEmpty,
			ErrDisplayNameEmpty,
			ErrGivenNameEmpty,
			ErrUserNameEmpty,
			ErrUserUserNameEmpty,
			ErrPrimaryEmailEmpty,
			ErrAddressesTooMany,
			ErrPhoneNumbersTooMany,
			ErrTooManyPrimaryEmails,
			ErrNameEmpty,
			ErrEmailValueEmpty,
		}

		for _, err := range errorVars {
			assert.Contains(t, err.Error(), "aws:", "Error should have 'aws:' prefix: %s", err.Error())
		}
	})

	t.Run("should have descriptive error messages", func(t *testing.T) {
		testCases := []struct {
			err      error
			contains string
		}{
			{ErrURLEmpty, "url may not be empty"},
			{ErrCreateGroupRequestEmpty, "create group request may not be empty"},
			{ErrBearerTokenEmpty, "bearer token may not be empty"},
			{ErrServiceProviderConfigEmpty, "service provider config may not be empty"},
			{ErrNameEmpty, "name may not be nil"},
			{ErrEmailValueEmpty, "email value may not be empty"},
			{ErrUserNameEmpty, "user name may not be empty"},
			{ErrDisplayNameEmpty, "display name may not be empty"},
		}

		for _, tc := range testCases {
			assert.Contains(t, tc.err.Error(), tc.contains)
		}
	})
}
