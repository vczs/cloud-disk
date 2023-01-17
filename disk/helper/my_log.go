package helper

import (
	"fmt"
	"runtime"
)

func VczsLog(desc string, err error) {
	p, _, line, _ := runtime.Caller(1)
	name := runtime.FuncForPC(p).Name()
	fmt.Printf("vczs_log_err: [%s(%d)](%s): %v\n", name, line, desc, err)
}
