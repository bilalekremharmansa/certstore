package shell

import (
	"fmt"
	"os"
	"testing"

	"bilalekrem.com/certstore/internal/assert"
	"bilalekrem.com/certstore/internal/pipeline/context"
)

func TestShellAction(t *testing.T) {
	tempDir, err := os.MkdirTemp("/tmp/", "*_test")
	assert.NotError(t, err, "creating temp dir")
	defer os.Remove(tempDir)

	// ------

	testFileName := "file"
	testFile := fmt.Sprintf("%s/%s", tempDir, testFileName)

	// ------

	args := map[string]string{}
	args["command"] = fmt.Sprintf("/usr/bin/touch %s", testFile)

	action := NewShellAction()
	err = action.Run(context.New(), args)
	assert.NotError(t, err, "running action")

	// -------

	files, err := os.ReadDir(tempDir)
	assert.NotError(t, err, "reading file")

	assert.Equal(t, 1, len(files))

	file := files[0]
	assert.Equal(t, testFileName, file.Name())
}

func TestShellActionWithError(t *testing.T) {
	args := map[string]string{}
	// this command will fail, because there /tmp is a directory.
	args["command"] = "/bin/mkdir /tmp"

	action := NewShellAction()
	err := action.Run(context.New(), args)
	assert.Error(t, err, "running action should've been failed")
}
