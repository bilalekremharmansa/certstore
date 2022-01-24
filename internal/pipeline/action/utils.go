package action

import (
	"errors"
	"fmt"

	"bilalekrem.com/certstore/internal/pipeline/context"
)

func ValidateRequiredArgs(actionArgs map[string]string, args ...string) error {
	for _, arg := range args {
		_, exists := actionArgs[arg]
		if !exists {
			return errors.New(fmt.Sprintf("required argument: %s", arg))
		}
	}

	return nil
}

func ValidateContextObjectExists(context *context.Context, keys ...context.Key) error {
	for _, key := range keys {
		val := context.GetValue(key)
		if val == nil {
			return errors.New(fmt.Sprintf("required context object: %s", key))
		}
	}

	return nil
}
