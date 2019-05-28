package logger

import (
	"log"
	"os"
	"testing"
)

func TestError(t *testing.T) {
	lf := logFile{
		logger:         log.New(os.Stdout, "TestLog", 1),
		isDebugEnabled: true,
	}

	lf.Error("ErrorMessage", nil)
}

func TestInfo(t *testing.T) {
	lf := logFile{
		logger:         log.New(os.Stdout, "TestLog", 1),
		isDebugEnabled: true,
	}

	lf.Error("InfoMessage", nil)
}

func TestDebug(t *testing.T) {
	lf := logFile{
		logger:         log.New(os.Stdout, "TestLog", 1),
		isDebugEnabled: true,
	}

	lf.Error("DebugMessage", nil)
}

func TestPrivatewrite(t *testing.T) {
	lf := logFile{
		logger:         log.New(os.Stdout, "TestLog", 1),
		isDebugEnabled: true,
	}
	lf.write("error", "TestErrorMessage", nil)
	lf.write("info", "TestInfoMessage", nil)
	lf.write("debug", "TestDebugMessage", nil)
}
