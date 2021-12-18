package action

import (
	"fmt"
	"os"
	"testing"
)

func TestShellAction(t *testing.T) {
	tempDir, err := os.MkdirTemp("/tmp/", "*_test")
	if err != nil {
		t.Fatalf("creating temp dir failed %v", err)
	}
	defer os.Remove(tempDir)

	// ------

	testFileName := "file"
	testFile := fmt.Sprintf("%s/%s", tempDir, testFileName)

	// ------

	args := map[string]string{}
	args["command"] = fmt.Sprintf("/usr/bin/touch %s", testFile)

	action := ShellAction{}
	err = action.run(args)
	if err != nil {
		t.Fatalf("running shell action failed %v", err)
	}

	// -------

	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("reading file failed %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected file size in temp dir is 1, found %d", len(files))
	}

	file := files[0]
	if testFileName != file.Name() {
		t.Fatalf("expected file name in temp dir is %s, actual %s", testFileName, file.Name())
	}
}

func TestShellActionWithError(t *testing.T) {
	args := map[string]string{}
	// this command will fail, because there /tmp is a directory.
	args["command"] = "/bin/mkdir /tmp"

	action := ShellAction{}
	err := action.run(args)
	if err == nil {
		t.Fatal("expected to have an error, but did not failed")
	}
}
