package log

import (
	"fmt"
	"runtime"
)

//格式化日志输出
func Log(txt interface{}) {
	pc, file, line, ok := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	fmt.Println("【pkgame】", fmt.Sprintf("func = %s,file = %s,line = %d,ok = %v ,val = %v", f.Name(), file, line, ok, txt))
}
