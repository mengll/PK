package main

import (
	"PK/libs/utils"
	"PK/modules/http/v1/game"
	"PK/modules/http/v1/user"
	"PK/modules/network/v1"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
)

func main() {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	path := utils.GetCurrentDirectory()

	fmt.Println(path)
	e.Static("/", path+"/static/bottle/client/") //创建服务
	e.GET("/gameserver", v1.WsPage)
	e.GET("/auth/callback", user.AuthCallback)

	gv1 := e.Group("/v1/")
	gv1.POST("user_game_result", user.GetUserGameInfo)
	gv1.POST("game_result_list", game.GetGameResList)
	gv1.POST("get_server_list", game.GetGameServerList)
	gv1.POST("config", game.GameConfig)
	e.Logger.Fatal(e.Start(":1323"))

}
