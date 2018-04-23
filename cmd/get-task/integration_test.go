package main

import (
	"testing"

	"github.com/smalleats/serverless-todo-example/internal"
)

const template = "../../cloudformation.yaml"

var cmdTests = []struct {
	description string
	event       string
	expected    string
}{
	{
		description: "return error for invalid task id",
		event: `{
			"task": "foo-bar"
		}`,
		expected: `{
			"code": "TASK_NOT_FOUND",
			"public_message": "Task not found",
			"private_message": "unknown task: foo-bar"
		}`,
	},
	{
		description: "return task",
		event: `{
			"task": "6eb69ac1-14fb-48b5-9c06-a82670342384"
		}`,
		expected: `{
			"id": "6eb69ac1-14fb-48b5-9c06-a82670342384",
			"done": false,
			"note": "Take the umbrella ☔",
			"created_at": "2018-04-23T11:04:00Z",
			"updated_at": "2018-04-23T11:04:00Z"
		}`,
	},
}

func TestCommand(t *testing.T) {
	for _, tt := range cmdTests {

		res, err := internal.SamInvoke(template, "GetTask", tt.event)
		if err != nil {
			t.Error("unexpected error:", err)
		}

		equal, err := internal.EqualJSON(tt.expected, res)
		if err != nil {
			t.Error("unexpected error:", err)
		}

		if !equal {
			t.Errorf("expected:\n%v\n\ngot:\n%v\n", tt.expected, res)
		}
	}
}
