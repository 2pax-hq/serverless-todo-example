package main

import (
	"encoding/json"
	"testing"

	"github.com/smalleats/serverless-todo-example/internal"
	"github.com/smalleats/serverless-todo-example/todo"
)

const template = "../../cloudformation.yaml"

var cmdTests = []struct {
	description string
	event       string
	expected    string
}{
	{
		description: "return error for empty note",
		event: `{
			"note": ""
		}`,
		expected: `{
			"code": "INVALID_INPUT",
			"public_message": "Invalid input",
			"private_message": "validation failed: missing note value"
		}`,
	},
	{
		description: "return error for empty note",
		event:       `{}`,
		expected: `{
			"code": "INVALID_INPUT",
			"public_message": "Invalid input",
			"private_message": "validation failed: missing note value"
		}`,
	},
}

func TestCommand(t *testing.T) {
	for _, tt := range cmdTests {
		t.Log(tt.description)

		res, err := internal.SamInvoke(template, "AddTask", tt.event)
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

func TestCommandSuccess(t *testing.T) {
	t.Log("create task")

	res, err := internal.SamInvoke(template, "AddTask", `{"note": "foo"}`)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	tsk := todo.Task{}
	err = json.Unmarshal([]byte(res), &tsk)
	if err != nil {
		t.Error("unmarshalling error:", err)
	}

	if tsk.Note != "foo" {
		t.Errorf("expected note:\nfoo\n\ngot response:\n%v\n", res)
	}

}
