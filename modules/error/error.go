package error

import (
	"os"
	"fmt"
	"PK/libs/utils"
)

type GpError struct {
	Error error
	Msg   string
	Date  string
}


//错误打印输出
func CheckError(err error,msg string){
	if err != nil{
		fmt.Println(err.Error(),msg)
		os.Exit(1)
	}
}

//创建新错误
func NewGpError(err error) GpError{
	gp := GpError{}
	gp.Error = err
	gp.Date = utils.GetNowTime()
	return gp

}
