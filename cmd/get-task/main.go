package main

import (
	"context"

	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/smalleats/serverless-todo-example/errors"
	"github.com/smalleats/serverless-todo-example/todo"
)

type request struct {
	Task string `json:"task"`
}

type getter interface {
	Get(string) (todo.Task, error)
}

type handler func(request) (todo.Task, error)

func getHandler(store getter) handler {
	return func(r request) (todo.Task, error) {
		t, err := store.Get(r.Task)

		switch err.(type) {
		case nil:
			return t, nil
		case todo.UnknownTaskError:
			return t, errors.Wrap(err, "TASK_NOT_FOUND", "Task not found")
		default:
			return t, errors.WithCode(err, errors.CodeApplicationError)
		}
	}
}

func fromRaw(h handler) func(context.Context, json.RawMessage) (todo.Task, error) {
	return func(ctx context.Context, event json.RawMessage) (todo.Task, error) {
		var r request
		if err := json.Unmarshal(event, &r); err != nil {
			return todo.Task{}, errors.WithCode(err, errors.CodeBadInput)
		}
		return h(r)
	}
}

func main() {
	store := todo.MockStore{}
	lambda.Start(fromRaw(getHandler(store)))
}
