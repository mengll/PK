package conf

type Dbdat struct {
	Host     string
	User     string
	PassWord string
	Port     string
	DataBase string
	DB       int
}

//链接pg数据 pg10 与天眼同库
var PgConfAdt Dbdat = Dbdat{
	Host		: "192.168.1.241", //"192.168.1.53",
	PassWord	: "adttianyan",
	Port		: "4453",//"5432",
	User		: "regina",
	DataBase	: "game_platform",
}

//redis 数据配置
var RedisConf Dbdat = Dbdat{
	Host        : "192.168.1.246",
	Port        : "6379",
	PassWord    : "",
	DB          : 1,
}