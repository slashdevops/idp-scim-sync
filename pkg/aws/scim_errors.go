package aws

import "fmt"

type HTTPResponseError struct {
	StatusCode int    `json:"StatusCode"`   // Http status code
	Code       string `json:"ErrorCode"`    // Datahub error code
	Message    string `json:"ErrorMessage"` // Error msg of the error code
}

func (e *HTTPResponseError) Error() string {
	return fmt.Sprintf("statusCode: %d,  errCode: %s, errMsg: %s", e.StatusCode, e.Code, e.Message)
}
