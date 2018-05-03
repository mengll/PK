package user

import (
	"PK/conf"
	"PK/modules/auth/anfeng"
	"fmt"
	"github.com/labstack/echo"
)

func GetUserGameInfo(c echo.Context) error {
	return nil
}

//
func AuthCallback(c echo.Context) error {
	fmt.Println("【授权回调】")
	accessToken, err := anfeng.Default.AccessToken(c.Scheme()+"://"+c.Request().Host+"/auth/callback", "STATE", c.QueryParam("code"))
	if err != nil {
		return err
	}
	var authorizeURL string
	conf.Get("web.authorize_url", &authorizeURL)
	return c.Redirect(302, authorizeURL+accessToken)
}
