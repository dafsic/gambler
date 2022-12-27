package utils

import (
	"fmt"
	"runtime"
)

// LineNo 返回调用此函数的代码所在函数、文件、行号
func LineNo() string {
	function := "xxx"
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	function = runtime.FuncForPC(pc).Name()
	return fmt.Sprintf(" -> %s():%s:%d", function, file, line)
	//return fmt.Sprintf(" -> %s:%d", file, line)
}
