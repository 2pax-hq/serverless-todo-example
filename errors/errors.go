// Package errors provides a boundary error type to be returned by cmds executed
// within the AWS Lambda environment.
//
// See the Lambda Go integration library: https://github.com/aws/aws-lambda-go/
//
// The Lambda Go integration library wraps a Function error's string representation
// in its own error type, so in order to preserve the original error structure
// the we use a JSON string as our Error type's string representation,
// resulting in a JSON within JSON response payload.
// https://github.com/aws/aws-lambda-go/blob/master/lambda/messages/messages.go#L32
//
// In cases where the Lambda is used as an API Gateway integration this allows us to
// deserialise our JSON string in a response mapping template and extract the error code,
// public message and private message separately.
package errors

import (
	"encoding/json"
	"fmt"
	"log"
)

type lambdaError struct {
	code    string
	message string

	origErr error
}

func (e lambdaError) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		log.Println("cannot marshal Error:", e)
		panic(err)
	}
	return string(b[:])
}

func (e lambdaError) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Code           string `json:"code"`
		PublicMessage  string `json:"public_message"`
		PrivateMessage string `json:"private_message"`
	}{
		Code:           e.code,
		PublicMessage:  e.message,
		PrivateMessage: e.origErr.Error(),
	})
}

// Generic error codes for errors returned by a Lambda function. The code can be
// used by API Gateway to reliably map an error response.
const (
	// CodeApplicationError is a catch-all for internal errors.
	// API Gateway mapping:     500 Internal server error
	CodeApplicationError = "APPLICATION_ERROR"

	// CodeAccessDenied represents an authorization error.
	// API Gateway mapping:     403 forbidden
	CodeAccessDenied = "ACCESS_DENIED"

	// CodeBadInput represents a bad Lambda input error.
	// API Gateway mapping:     400 Bad request
	CodeBadInput = "BAD_INPUT"

	// CodeInvalidInput represents a bad Lambda input value error.
	// API Gateway mapping:     422 Unprocessable entity
	CodeInvalidInput = "INVALID_INPUT"
)

var defaultMessages = map[string]string{
	CodeApplicationError: "Application error",
	CodeAccessDenied:     "Access denied",
	CodeBadInput:         "Bad input",
	CodeInvalidInput:     "Invalid input",
}

// Wrap returns a Lambda error for the original error with a specified code
// and message.
func Wrap(err error, code, msg string) error {
	if err == nil {
		return nil
	}
	return lambdaError{
		code:    code,
		message: msg,
		origErr: err,
	}
}

// Wrapf returns a Lambda error for the original error with a specified code and
// string format.
func Wrapf(err error, code, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return lambdaError{
		code:    code,
		message: fmt.Sprintf(format, args...),
		origErr: err,
	}
}

// WithCode returns a Lambda error for the original error with a specified code
// and derives the message from the code. If there's no message defined for the
// code the public message will be empty.
func WithCode(err error, code string) error {
	if err == nil {
		return nil
	}
	return lambdaError{
		code:    code,
		message: defaultMessages[code],
		origErr: err,
	}
}
