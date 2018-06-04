package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"runtime"
	"math/rand"
	"reflect"
	"sync"
	uuid "github.com/satori/go.uuid"
)

//时间转化函数
const (
	format_date = "2006-04-03 15:04:05"
	ROOM_TYPE_SHARE = "create"
)

//urlu转码操作
func Urlencode(str string) string {
	return url.PathEscape(str)
}

//urldecode反码
func Urldecode(str string) string {
	vk, errp := url.PathUnescape(str)
	if errp != nil {
		fmt.Println("转化错误", errp)
	}
	return vk
}

//md5 32 小写
func M5(str string) string {
	md5 := md5.New()
	md5.Write([]byte(str))
	mds := hex.EncodeToString(md5.Sum(nil)) // 转
	return strings.ToLower(mds)
}

//简单异或操作
func SimpleorX(info, key []byte) string {
	if len(info) == 0 || len(key) == 0 {
		return ""
	}
	var back []byte
	k_len := len(key)
	d_len := len(info)

	for i := 0; i < d_len; i++ {
		key_len := i % k_len
		back = append(back, info[i]^key[key_len])
	}

	return base64.URLEncoding.EncodeToString(back)
}

//数据类型转化 float64 to int
func IntFromFloat64(x float64) int {
	if math.MinInt32 <= x && x <= math.MaxInt32 { // x lies in the integer range
		whole, fraction := math.Modf(x)
		if fraction >= 0.5 {
			whole++
		}
		return int(whole)
	}
	panic(fmt.Sprintf("%g is out of the int32 range", x))
}

//map 转字符串
func MaptoJson(data map[string]interface{}) string {
	configJSON, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return ""
	}
	return string(configJSON) //返回格式化后的字符串的内容0
}

//str to map
func StrToMap(data string) map[string]interface{} {
	var dat map[string]interface{}
	json.Unmarshal([]byte(data), &dat)
	return dat
}

//将str转换为时间格式
func StrToTime(st string) time.Time {
	t, _ := time.ParseInLocation(format_date, st, time.Local) //时间戳转化
	return t
}

//现在的格式化2018-06-04 16:06:27
func GetNowTime() string {
	return time.Now().Format("2006-04-03 15:04:05")
}

//获取当前的时间戳
func GetTimeStamp() int {
	t, err := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		fmt.Println("时间戳转化失败", err.Error())
	}
	return t
}

func TodayLeftTime() int {
	tf := time.Now()
	//tf.Sub()
	return (23 - tf.Hour())* 3600 + (60-tf.Minute()) * 60 + 60 - tf.Second()
}

/*
获取程序运行路径
*/
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("This is an error")
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//字符串的utf 真实长度
func StrLength( str string) int{
	return utf8.RuneCount([]byte(str))
}

//格式化日志输出
func Log(txt interface{}) {
	pc, file, line, ok := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	fmt.Println("【pkgame】", fmt.Sprintf("func = %s,file = %s,line = %d,ok = %v ,val = %v", f.Name(), file, line, ok, txt))
}

//生成随机数
func RandomNum()int{
	m := time.Now().Second()
	return 100 + m + rand.Intn(m)
}

//判断当前是否是分享
func IsShare(room string )int{
	if strings.HasSuffix(room,ROOM_TYPE_SHARE){
		return 1
	}
	return 0
}

//转化成相关的类型
func NewReflectInterface(v interface{}) interface{}{
	tp := reflect.ValueOf(v).Type()
	vb := reflect.New(tp).Interface()
	return vb
}

//单运行程序
type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap(cb func(argvs ...interface{}), argvs ...interface{}) {
	w.Add(1)
	go func() {
		cb(argvs...)
		w.Done()
	}()
}

//生成uuid
func Uuid() string {
	uuid ,err :=  uuid.NewV4()
	if err != nil{
		Log(err)
	}
	return uuid.String()
}
