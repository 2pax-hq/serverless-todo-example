package errors_test

import (
	"encoding/json"
	goerrors "errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/smalleats/serverless-todo-example/errors"
)

func TestWrap(t *testing.T) {
	exp := `{
      "code": "UH_OH",
      "public_message": "something went wrong",
      "private_message": "simulated error"
    }`

	err := errors.Wrap(
		goerrors.New("simulated error"),
		"UH_OH",
		"something went wrong",
	)
	res := err.Error()
	equal, e := EqualJSON(exp, res)
	if e != nil {
		t.Errorf("unexpected error: %v", e)
	} else if !equal {
		t.Errorf("expexted:\n\t%s\ngot:\n%s\t", exp, res)
	}
}

func TestWrapNil(t *testing.T) {
	err := errors.Wrap(
		nil,
		"UH_OH",
		"something went wrong",
	)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWrapf(t *testing.T) {
	exp := `{
      "code": "UH_OH",
      "public_message": "something went wrong: 🙈",
      "private_message": "simulated error"
    }`

	err := errors.Wrapf(
		goerrors.New("simulated error"),
		"UH_OH",
		"something went wrong: %s", "🙈",
	)
	res := err.Error()
	equal, e := EqualJSON(exp, res)
	if e != nil {
		t.Errorf("unexpected error: %v", e)
	} else if !equal {
		t.Errorf("expexted:\n\t%s\ngot:\n%s\t", exp, res)
	}
}

func TestWrapfNil(t *testing.T) {
	err := errors.Wrapf(
		nil,
		"UH_OH",
		"something went wrong: %s", "🙈",
	)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

var withCodeTests = []struct {
	err      error
	code     string
	expected string
}{
	{
		err:  goerrors.New("simulated error"),
		code: errors.CodeApplicationError,
		expected: `{
      "code": "APPLICATION_ERROR",
      "public_message": "Application error",
      "private_message": "simulated error"
    }`,
	},
	{
		err:  goerrors.New("simulated error"),
		code: errors.CodeAccessDenied,
		expected: `{
      "code": "ACCESS_DENIED",
      "public_message": "Access denied",
      "private_message": "simulated error"
    }`,
	},
	{
		err:  goerrors.New("simulated error"),
		code: errors.CodeBadInput,
		expected: `{
      "code": "BAD_INPUT",
      "public_message": "Bad input",
      "private_message": "simulated error"
    }`,
	},
	{
		err:  goerrors.New("simulated error"),
		code: errors.CodeInvalidInput,
		expected: `{
      "code": "INVALID_INPUT",
      "public_message": "Invalid input",
      "private_message": "simulated error"
    }`,
	},
}

func TestWithCode(t *testing.T) {
	for _, test := range withCodeTests {
		err := errors.WithCode(test.err, test.code)
		res := err.Error()

		equal, e := EqualJSON(test.expected, res)
		if e != nil {
			t.Errorf("unexpected error: %v", e)
		} else if !equal {
			t.Errorf("expexted:\n\t%s\ngot:\n%s\t", test.expected, res)
		}
	}
}

func TestWithCodefNil(t *testing.T) {
	err := errors.WithCode(nil, errors.CodeApplicationError)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// EqualJSON is a convenience method to compare the JSON structure of two
// strings.
func EqualJSON(lhs, rhs string) (bool, error) {
	var l, r interface{}
	if err := json.Unmarshal([]byte(lhs), &l); err != nil {
		return false, fmt.Errorf("error parsing left string: %s", err.Error())
	}
	if err := json.Unmarshal([]byte(rhs), &r); err != nil {
		return false, fmt.Errorf("error parsing right string: %s", err.Error())
	}
	return reflect.DeepEqual(l, r), nil
}
