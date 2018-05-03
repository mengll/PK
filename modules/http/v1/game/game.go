package game

import (
	"PK/common"
	"PK/conf"
	"PK/modules/auth/anfeng"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

func GetGameResList(c echo.Context) error {

	return nil
}

/*
返回游戏的服务器列表
*/
func GetGameServerList(c echo.Context) error {
	return nil
}

//返回游戏配置信息
func GameConfig(c echo.Context) (err error) {
	type Data struct {
		Websocket string          `json:"websocket"`
		Weixin    json.RawMessage `json:"weixin"`
	}

	data := new(Data)
	if err = conf.Get("web", &data); err != nil {
		return
	}

	var weixin json.RawMessage

	if weixin, err = anfeng.Default.WeixinSDK(c.FormValue("url")); err != nil {
		fmt.Println(err.Error())
		return
	}

	data.Weixin = weixin

	Res := common.ResponeDat{
		ErrorCode: conf.SUCESS_BACK,
		Data:      data,
	}

	return c.JSON(http.StatusOK, Res)
}
