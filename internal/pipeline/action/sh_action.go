package action

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"bilalekrem.com/certstore/internal/logging"
	"bilalekrem.com/certstore/internal/pipeline/context"
)

type ShellAction struct {
}

func (ShellAction) Run(ctx *context.Context, args map[string]string) error {
	validate(args)

	// ----

	commandStr := args["command"]
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
	_, exist := args["command"]
	if !exist {
		return errors.New("shell action expects a args.command, but could not found!")
	}

	return nil
}
