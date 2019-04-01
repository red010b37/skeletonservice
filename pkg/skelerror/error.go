package skelerror

import "net/http"

type ServiceError struct {
	key       string // Identifies the error
	err       string // More user friendly error description
	errorCode int    // The status code of the error - we use http status codes as they are universally understood
}

func (e *ServiceError) OverrideErrorMessage(errMsg string) {
	e.err = errMsg
}

func (e *ServiceError) Error() string {
	return e.err
}

func (e *ServiceError) Key() string {
	return e.key
}

func (e *ServiceError) ErrorCode() int {
	return e.errorCode
}

var (
	RequestCastError = ServiceError{
		"REQUEST_CAST_ERROR",
		"There was an error casting the request to the right type",
		http.StatusBadRequest,
	}

	ResponseCastError = ServiceError{
		"RESPONSE_CAST_ERROR",
		"There was an error casting the response to the right type",
		http.StatusBadRequest,
	}

	InvalidCredentials = ServiceError{
		"INVALID_CREDENTIALS",
		"Invalid Credentials provided",
		http.StatusUnauthorized,
	}

	JSONDecodeError = ServiceError{
		"JSON_DECODE_ERROR",
		"Unable to decode JSON",
		http.StatusBadRequest,
	}

	JSONEncodeError = ServiceError{
		"JSON_ENCODE_ERROR",
		"Unable to encode JSON",
		http.StatusServiceUnavailable,
	}

	BodyReadError = ServiceError{
		"BODY_READ_ERROR",
		"Unable read the request body",
		http.StatusInternalServerError,
	}

	IncorrectBillingState = ServiceError{
		"INCORRECT_BILLING_STATE",
		"Billg state was expected to be something else",
		http.StatusBadRequest,
	}

	InvalidToken = ServiceError{
		"INVALID_TOKEN",
		"Token is invalid",
		http.StatusUnauthorized,
	}

	ParamsError = ServiceError{
		"PARAMS_ERROR",
		"",
		http.StatusBadRequest,
	}
)
