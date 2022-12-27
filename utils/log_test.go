package utils

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	tl := NewLogger(os.Stdout, "TEST", LDebug, Ldefault)
	tl.Errorf("this should be print msg -> %s\n", "hahaha")
	tl.Debugf("this shouldn't print msg -> %s\n", "ememem")
}
