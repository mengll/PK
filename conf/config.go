package conf

import (
	"PK/libs/utils"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
	"os"
	"fmt"
)

//Config 配置
type Config struct {
	content []byte
}

//New 配置实例
func New(filename string) (c *Config) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	c = &Config{content}
	return
}

//Default 默认配置实例

var path string = utils.GetCurrentDirectory()
var Default = New(path + "/config.json")

//Get 读取配置
func (config *Config) Get(name string, v interface{}) (err error) {
	keys := strings.Split(name, ".")
	c := config.content
	ok := false
	for _, k := range keys {
		var m = map[string]json.RawMessage{}
		json.Unmarshal(c, &m)
		c, ok = m[k]
		if !ok {
			return errors.New("config: unmarshal error " + name)
		}
	}
	return json.Unmarshal(c, v)
}

//Get 读取默认配置
func Get(name string, v interface{}) (err error) {
	return Default.Get(name, v)
}

//config Obj
var Configobj map[string]map[string]interface{}

func init(){
	Configobj = make(map[string]map[string]interface{})
	f, err1 := os.OpenFile(path + "/config.json", os.O_RDONLY, 0666)

	if err1 != nil {
		utils.Log(err1)
	}

	err := json.NewDecoder(f).Decode(&Configobj)
	if err != nil {
		fmt.Println(err)
	}
}

//获取对象值
func GetConfig(db,key string) interface{} {

	if obj,ok := Configobj[db];ok{
		if k,o := obj[key];o{
			return k
		}else {
			panic("不存在")
		}
	}else{
		panic("不存在")
	}

	return ""
}
