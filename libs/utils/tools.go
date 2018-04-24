package utils

import (
	"net/url"
	"fmt"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"encoding/base64"
	"math"
	"encoding/json"
	"time"
	"strconv"
)

//时间转化函数
const format_date = "2006-04-03 15:04:05"

//urlu转码操作
func Urlencode(str string) string{
	return url.PathEscape(str)
}

//urldecode反码
func Urldecode(str string) string{
	vk,errp := url.PathUnescape(str)
	if errp != nil{
		fmt.Println("转化错误",errp)
	}
	return vk
}

//md5 32 小写
func M5(str string) string{
	md5 := md5.New()
	md5.Write([]byte(str))
	mds := hex.EncodeToString(md5.Sum(nil)) // 转
	return strings.ToLower(mds)
}

//简单异或操作
func SimpleorX(info,key []byte) string{
	if len(info) ==0 || len(key) == 0{
		return ""
	}
	var back []byte
	k_len := len(key)
	d_len := len(info)

	for i := 0;i < d_len;i++{
		key_len := i % k_len
		back = append(back,info[i] ^ key[key_len])
	}

	return  base64.URLEncoding.EncodeToString(back)
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
	t ,err := strconv.Atoi(strconv.FormatInt(time.Now().Unix(),10))
	if err != nil{
		fmt.Println("时间戳转化失败",err.Error())
	}
	return t
}
