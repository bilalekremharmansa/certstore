package shell

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"bilalekrem.com/certstore/internal/logging"
	"bilalekrem.com/certstore/internal/pipeline/action"
	"bilalekrem.com/certstore/internal/pipeline/context"
)

const (
	ARGS_COMMAND string = "command"
)

type shellAction struct {
}

func NewShellAction() *shellAction {
	return &shellAction{}
}

func (shellAction) Run(ctx *context.Context, args map[string]string) error {
	validate(args)

	// ----

	commandStr := args[ARGS_COMMAND]
	splitted := strings.Split(commandStr, " ")

	cmd := splitted[0]
	cmdArgs := splitted[1:]

	logging.GetLogger().Debugf("executing command: [%s], [%s]", cmd, cmdArgs)
	result := exec.Command(cmd, cmdArgs...)

	output, _ := result.Output()
	logging.GetLogger().Debugf("command output: [%s], exit code: %d", string(output), result.ProcessState.ExitCode())

	if result.ProcessState.Success() {
		return nil
	}

	return errors.New(fmt.Sprintf("Error occurred while executing command [%s], exit code: %d",
		cmd, result.ProcessState.ExitCode()))
}

func validate(args map[string]string) error {
	err := action.ValidateRequiredArgs(args, ARGS_COMMAND)
	if err != nil {
		return err
	}

	return nil
}
