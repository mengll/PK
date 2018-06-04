package v1

import (
	"PK/common"
	"PK/conf"
	"PK/modules/auth/anfeng"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"strconv"
	"time"
	"PK/libs/db"
	"context"
	"math/rand"
	"strings"
)

var (
	Pf_redis      = db.NewRedis()
	CancelFuc = make(map[string] chan int)
	UidRooms = make(map[string]string)
	ShareRoom = make(map[string]bool)
)

func init(){
	Pf_redis.Connect()
}

func (c *Client) Auth(Res common.ResponeDat) error {
	authorizeURL := Call_back_url
	if len(authorizeURL) >= 500 || authorizeURL == "" {
		return errors.New("回调地址出错")
	}
	fmt.Println("【现在运行授权】")
	data := make(map[string]interface{})
	data["url"] = authorizeURL
	Res.ErrorCode = conf.SUCESS_BACK
	Res.Data = data
	Res.Msg = ""
	c.send <- Res

	return nil
}

//处理用户登录
func (c *Client) Login(Res common.ResponeDat, access_token string) (err error) {

	profile, err := anfeng.Default.Profile(access_token)
	if err != nil {
		Res.ErrorCode = conf.FAILED_BACK
		Res.Msg = err.Error()
		c.send <- Res
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("【现在开始登录】")
	udat := new(common.WSDat)
	udat.Uid = strconv.Itoa(profile.UID)
	udat.Avatar = profile.Avatar
	udat.NickName = profile.UserName
	udat.Gender = strconv.Itoa(profile.Gender)
	udat.Brithday = profile.Birthday

	UIDS[c.socket] = udat.Uid

	Res.ErrorCode = conf.SUCESS_BACK
	Res.Msg = ""
	Res.Data = udat
	c.send <- Res

	return nil
}

//进入游戏
func (c *Client) EnetrGame(Res common.ResponeDat, uid, game_id string) error {

	if uid == "" {
		return errors.New("enter game uid is null")
	}

	//保存用户登录信息
	login_key := fmt.Sprintf(conf.CLIENT_LOGIN_KYE, game_id)
	Pf_redis.AddSet(login_key, uid)

	online_key := fmt.Sprintf(conf.ONLINE_KEY, uid)
	Pf_redis.Expire(online_key, time.Second*30)

	//todo 保存进入游戏信息
	if _, ok := Plat_form_user[game_id]; !ok {
		Plat_form_user[game_id] = make(map[string]*websocket.Conn)
	}

	Plat_form_user[game_id][uid] = c.socket
	back_dat := make(map[string]interface{})
	back_dat["online_num"] = Pf_redis.GetSetNum(login_key)
	back_dat["game_id"] = game_id

	//返回登录
	Res.ErrorCode = conf.SUCESS_BACK
	Res.Msg = ""
	Res.Data = back_dat
	c.send <- Res

	return nil
}

//加入房间
func (c *Client) JoinRoom(Res common.ResponeDat, req_data *common.ReqDat) error{

	uid 	:= req_data.Data["uid"].(string)
	game_id := req_data.Data["game_id"].(string)
	room 	:= req_data.Data["room"].(string)
	can_broad_cast,err := Roomer.JoinRoom(room,uid)

	if err != nil{
		Res.ErrorCode = conf.FAILED_BACK
		Res.Msg = err.Error()
		c.send <- Res
	}

	if can_broad_cast {
		dat := SocketData{}
		dat.From = c.socket
		trace_data := make(map[string]interface{})
		trace_data["game_id"] = game_id
		dat.Data = trace_data

		BroadCast(dat) //广播通知当前的玩家，
	}

	return nil
}

//获取要发送给谁
func (c *Client) GetToUser(){

}

//开始匹配
func (c *Client) SearchMatch(dat *common.ReqDat) error{
	if _,ok := dat.Data["game_id"];!ok{
		return errors.New("【game_id is null】")
	}
	game_id := dat.Data["game_id"].(string)
	room_limit,_ := Roomer.RoomLimit(game_id)

	switch room_limit{
	case 2:
		c.PtpMatch(game_id,room_limit)
	}

	//广播
	return nil
}

//p2p返回匹配到的用户的信息列表
func (c *Client) PtpMatch(game_id string ,room_limit int) *SocketData{

	socket_data := new(SocketData)
	gameReady := fmt.Sprintf(conf.GAME_REDAY_LIST, game_id) //所有准备的用户
	Pf_redis.AddSet(gameReady, c.uid)
	fmt.Printf("【uid】%v",c.uid)

	reday_num := Pf_redis.GetSetNum(gameReady)
	now_hour := time.Now().Hour()
	af_time := 10
	rand_int := rand.Intn(time.Now().Nanosecond())
	ra, _ := strconv.Atoi(strconv.Itoa(rand_int)[:1])

	if now_hour < 6 {
		af_time = 10 + ra
	} else {
		af_time = 2 + ra
	}

	dd := []string{}

	ctx_one,cancel_one := context.WithTimeout(context.Background(),time.Second * time.Duration(af_time))
	println("【房间人数限制】",room_limit)
	if reday_num >= room_limit{
		cancel_one()
		ctx,cancel := context.WithTimeout(context.Background(),time.Second *30)

		for{
			select {
			case <-ctx.Done():
				fmt.Println("【结束】")
				err_str := ctx.Err().Error()
				if strings.Contains(err_str, "deadline") {
					fmt.Println("【匹配超时】")
					goto SearchEnd
			}
				break

			default:

				uk := Pf_redis.SPop(gameReady)
				is_exists, err := Pf_redis.EXISTS(fmt.Sprintf(conf.ONLINE_KEY, uk))
				if err != nil {
					fmt.Println(err)
				}

				if uk != "" && is_exists == true {
					//uid
					if _,ok := CancelFuc[uk];ok{
						fmt.Println("【AI取消匹配】")
						CancelFuc[uk] <- 1
					}

					dd = append(dd, uk)
				}

				fmt.Printf("【当前匹配到的用户】%v",dd)
				if len(dd) == room_limit {
					//创建房间
					client_room,_ := Roomer.CreateRoom(game_id,c.uid)
					user_dat := make(map[string]interface{})

					socket_data.From = c.socket
					socket_data.To = []*websocket.Conn{}

					for _, v := range dd {
						if v != c.uid {
							socket_data.To = append(socket_data.To,Plat_form_user[game_id][v])
						}
						user_info := Pf_redis.GetKey(fmt.Sprintf(conf.USER_GAME_KEY, v))
						user_dat[v] = user_info
						Pf_redis.AddSet(client_room, v)
						//保存用户和房间信息
						UidRooms[v] = client_room
					}

					//Res.ErrorCode = conf.SUCESS_BACK
					//Res.Data = map[string]interface{}{"cmd": "start"}
					//Res.Msg = conf.START
					//BroadCast(client_room, game_id, "") //广播通知当前的玩家，

					dt,err := Pf_redis.SMembers(client_room)
					if err != nil{
						fmt.Println(err.Error())
					}

					fmt.Printf("【成功匹配1v1】%v",dt)
					dd = dd[:0]                         //清空
					cancel()

					return socket_data
				}

			}
		}

		SearchEnd:
		//超时匹配
	}else{
		CancelFuc[c.uid] = make (chan int)
		defer close(CancelFuc[c.uid])
		select {
			case <- ctx_one.Done():
				err_str := ctx_one.Err().Error()
				if strings.Contains(err_str, "deadline") {
					fmt.Println("【匹配到机器人】")
				}
			case <-	CancelFuc[c.uid]:
				fmt.Println("【取消机器人匹配】",c.uid)
				cancel_one() //取消当前的匹配
		}

	}

	return socket_data
}




