package utils

import (
	"os"

	"bilalekrem.com/certstore/internal/logging"
)

func ValidateNotError(err error) {
	if err != nil {
		Error("error occurred: %v", err)
	}
}

func Error(template string, args ...interface{}) {
	if len(args) > 0 {
		logging.GetLogger().Errorf(template, args)
	} else {
		logging.GetLogger().Error(template)
	}
	os.Exit(1)
}
