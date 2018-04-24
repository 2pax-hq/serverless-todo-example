package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	mrand "math/rand"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
	"time"
)

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

// SamInvoke uses SAM Local to invoke a Lambda Function handler locally and
// returns the response which can be the successful response or the encoded
// error response.
//
// Requirements:
// - You must have a `sam` binary installed, see https://github.com/awslabs/aws-sam-local
//   for more information.
// - The CloudFormation template must include a `AWS::Serverless::Function`
// 	 resource for the with the name specified in `function`.
// - The Lambda handler specified in the CloudFormation template resource must
//	 be compiled before running the test.
func SamInvoke(template, function, event string) (string, error) {

	mrand.Seed(time.Now().Unix())
	// create a tmp dir in the current directory because it needs to be mounted
	// within the container that SAM Local uses to execute the binary.
	d := fmt.Sprintf("test-tmp-%d", mrand.Int())
	err := os.Mkdir(d, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(d)

	err = ioutil.WriteFile(d+"/event.json", []byte(event), 0644)
	if err != nil {
		return "", err
	}

	cmd := exec.Command(
		"sam",
		"local",
		"invoke",
		function,
		"--template",
		template,
		"--event",
		d+"/event.json",
	)

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	if err = cmd.Run(); err != nil {
		return "", err
	}

	// If error response has been encoded and written to stderr, extract and the
	// JSON string from errorMessage.
	r := regexp.MustCompile(`errorMessage\": \"(.*)\",`)
	match := r.FindStringSubmatch(errb.String())
	if len(match) == 2 {
		// unescape json string
		if len(match) == 2 {
			// cleanup escape symbols
			str := strings.Replace(match[1], "\\", "", -1)
			return str, nil
		}
	}
	return outb.String(), nil
}
