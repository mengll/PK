package v1

import (
	"PK/common"
	"PK/conf"
	"PK/modules/auth/anfeng"
	"PK/modules/room"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"net/http"

)

type (
	SocketData struct {
		From *websocket.Conn
		To   []*websocket.Conn
		Data interface{}
	}

	Client struct {
		game_id string
		uid     string
		socket *websocket.Conn
		send   chan interface{}
	}

	ClientManager struct {
		clients    map[*Client]bool
		broadcast  chan *common.ReqDat
		register   chan *Client
		unregister chan *Client
	}
)

var manager = ClientManager{
	broadcast:  make(chan *common.ReqDat),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

//声明回调地址
var (
	Plat_form_user  = make(map[string]map[string]*websocket.Conn) //在线的用户的信息
	Call_back_url string
	UIDS          = make(map[*websocket.Conn]string)
	Roomer        = room.NewRoomManager()
)

//创建房间
func BroadCast(data SocketData) error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("BroadCast Message error ")
		}
	}()

	if len(data.To) == 0 {
		return errors.New("接收方，不能为空")
	}

	for _, ws := range data.To {
		err := ws.WriteJSON(data.Data)

		if err != nil {
			from := data.From
			if from != nil {
				from.WriteJSON("seend erro")
			}

			return err
		}
	}

	return nil
}

func (manager *ClientManager) start() {

	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
			//注册信息配置
			//manager.send(jsonMessage, conn)

		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				println("【channel close ws】")
				//close(conn.send)
				delete(manager.clients, conn)
				//manager.send(jsonMessage, conn)
			}

		case message := <-manager.broadcast:
			for conn := range manager.clients {
				fmt.Println(conn, message)
			}
		}

	}
}

func (manager *ClientManager) send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		if conn != ignore {
			conn.send <- message
		}
	}
}

func (c *Client) read() {

	defer func() {
		if ree := recover();ree != nil{
			fmt.Println("【读取时关闭】")
		}
		manager.unregister <- c
		c.socket.Close() //
	}()

	for {

		dat := &common.ReqDat{}
		err := c.socket.ReadJSON(dat)

		if err != nil {
			fmt.Println("【read data】",err.Error())
			manager.unregister <- c
			c.socket.Close() //
			break
		}

		Res := common.ResponeDat{}
		Res.MessageId = dat.MessageId

		game_id := ""
		uid := ""

		if UIDS[c.socket] != "" {
			uid = UIDS[c.socket]
		}

		if _, ok := dat.Data["game_id"]; ok {
			game_id = dat.Data["game_id"].(string)
		}

		switch dat.Cmd {

		case conf.AUTHORIZE:
			c.Auth(Res)
			//todo 保存用户信息
			//登录
		case conf.LOGIN:

			fmt.Println("【login】")
			if _, ok := dat.Data["access_token"]; !ok {
				fmt.Println(dat)
				continue
			}
			access_token := dat.Data["access_token"].(string)
			c.Login(Res, access_token)

		case conf.ENTER_GAME:
			c.game_id = game_id
			c.uid = uid
			c.EnetrGame(Res, uid, game_id)

			//创建房间
		case conf.CREATE_ROOM:
			_, err := Roomer.CreateRoom(game_id, uid)
			if err != nil {
				fmt.Println("【创建房间失败】")
			}
			//加入房间
		case conf.JOIN_ROOM:
			c.JoinRoom(Res,dat)

		case conf.SEARCH_MATCH:
			fmt.Println("【开始匹配】")
			c.SearchMatch(dat)
		} //end switch

	}
}

func (c *Client) write() {
	defer func() {
		if err := recover();err !=nil{
			fmt.Println("【client write close ws】",err)
		}
		manager.unregister <- c
		c.socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.socket.WriteJSON(message)
				continue
			}
			c.socket.WriteJSON(message)
		}
	}
}

func WsPage(c echo.Context) error {
	fmt.Println("【websocket】")
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Response(), c.Request(), nil)
	if error != nil {
		http.NotFound(c.Response(), c.Request())
		fmt.Println(error.Error())
		return error
	}

	if Call_back_url == "" {
		Call_back_url = anfeng.Default.AuthorizeURL("http://"+c.Request().Host+"/auth/callback", "STATE")
	}

	go manager.start()

	client := &Client{socket: conn, send: make(chan interface{})} //创建一个用户信息
	manager.register <- client

	fmt.Println("【网络通讯】")
	go client.read()
	go client.write()

	return nil
}
