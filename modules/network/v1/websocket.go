package v1

import (
	"github.com/gorilla/websocket"
	"fmt"
	"errors"
)

type SocketData struct {
	From *websocket.Conn
	To   []*websocket.Conn
	Data interface{}
}

//创建房间
func BroadCast(data SocketData) error {

	defer func() {
		if err := recover(); err != nil{
			fmt.Println("BroadCast Message error ")
		}
	}()

	if len(data.To) == 0{
		return errors.New("接收方，不能为空")
	}

	for _,ws := range data.To{
		err := ws.WriteJSON(data.Data)

		if err != nil{
			from := data.From
			if from != nil{
				from.WriteJSON("seend erro")
			}

			return err
		}
	}

	return nil
}

//读取信息
func WSread(){
	fmt.Println("This is a file")
}
