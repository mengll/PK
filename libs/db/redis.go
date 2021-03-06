package db

import (
	"smallgamepk.qcwanwan.com/config"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type (
	GsRedis interface {
		Connect()
		GetSetNum(string) int
		HadSet(key, val string) bool
		AddSet(key, val string) error
		SetKey(k string, v interface{}) error
		DelSet(key, val string) error
		GetKey(k string) string
		SPop(k string) string
		DelKey(k string) error
		SMembers(k string) ([]string, error)
		Expire(k string, t time.Duration) error
		EXISTS(k string) (bool, error)
		Incr(k string) error
	}

	GsRedisManage struct {
		Address   string `json:"address"`
		Passworld string `json:"password"`
		DB        int    `json:"db"`
		RS        *redis.Client
	}
)

func (this *GsRedisManage) Connect() {
	this.RS = redis.NewClient(&redis.Options{
		Addr:     this.Address,
		Password: this.Passworld, // 设置Redis的链接的链接方法
		DB:       this.DB,        // use default DB
	})
}

//从Redis的集合中移除数据
func (this *GsRedisManage) DelSet(key, val string) error {
	return this.RS.SRem(key, val).Err()
}

//设置key
func (this *GsRedisManage) SetKey(k string, v interface{}) error {
	return this.RS.Set(k, v, 0).Err()
}

//添加到集合中
func (this *GsRedisManage) AddSet(key, val string) error {
	return this.RS.SAdd(key, val).Err()
}

//判断是否存在
func (this *GsRedisManage) HadSet(key, val string) bool {
	return this.RS.SIsMember(key, val).Val()
}

//获取数据
func (this *GsRedisManage) GetKey(k string) string {
	return this.RS.Get(k).Val()
}

//随机删除并返回
func (this *GsRedisManage) SPop(k string) string {
	return this.RS.SPop(k).Val()
}

//删除键
func (this *GsRedisManage) DelKey(k string) error {
	return this.RS.Del(k).Err()
}

//获取集合中的内容
func (this *GsRedisManage) SMembers(k string) ([]string, error) {
	set_val := this.RS.SMembers(k)
	return set_val.Val(), set_val.Err()
}

//获取集合中的数量
func (this *GsRedisManage) GetSetNum(key string) int {
	checkNumTmp := this.RS.SCard(key).Val()
	dt := strconv.FormatInt(checkNumTmp, 10)
	dd, err := strconv.Atoi(dt)
	if err != nil {
		return 0
	}
	return dd
}

func (this *GsRedisManage) Incr(k string) error{
	return this.RS.Incr(k).Err()
}

//设置可以的实效时间
func (this *GsRedisManage) Expire(k string, t time.Duration) error {
	this.SetKey(k, 1)
	return this.RS.Expire(k, t).Err()
}

//判断当前的可以使是否存在
func (this *GsRedisManage) EXISTS(k string) (bool, error) {

	resut, err := this.RS.Exists(k).Result()

	if err != nil {
		return false, err
	}

	if resut > 0 {
		return true, nil
	}

	return false, nil

}


func(this *GsRedisManage)Clearn(){
	this.RS.Del("*")
}

//NewRedis 生成当前redis 对象
func NewRedis() GsRedis {
	redis_client := new(GsRedisManage)
	if err := config.Get("redis", redis_client); err != nil {
		panic(err)
	}
	return redis_client
}
