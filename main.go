package main

import (
	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo"
	"PK/modules/http/v1/user"
	"PK/modules/http/v1/game"
	"PK/libs/utils"
	"fmt"
)

func main() {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	//
	//e.GET("/auth/callback", server.AuthCallback)
	fmt.Println(utils.GetNowTime())
	gv1 := e.Group("/v1/")
	gv1.POST("user_game_result",user.GetUserGameInfo)
	gv1.POST("game_result_list",game.GetGameResList)
	e.Logger.Fatal(e.Start(":1323"))

}
