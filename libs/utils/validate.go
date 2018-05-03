package utils

//数据验证类
func Checkerror(err error) {
	if err != nil {
		panic(err)
	}
}
