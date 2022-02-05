package windns

import (
	"os/exec"
	"strings"

	"bilalekrem.com/certstore/internal/logging"
)

const POWERSHELL_PATH = "C:\\Windows\\syswow64\\WindowsPowerShell\\v1.0\\powershell.exe"

func runPowershellCmd(cmd string) (string, error) {
	cmdArgs := strings.Split(cmd, " ")

	logging.GetLogger().Infof("executing powershell cmd, %v", cmdArgs)
	result := exec.Command(POWERSHELL_PATH, cmdArgs...)
	output, err := result.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
