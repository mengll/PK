package common

type (
	ReqDat struct {
		Cmd        string                 `json:"cmd"`
		Data       map[string]interface{} `json:"data"`
		MessageId  string                 `json:"message_id"`
		MessageKey string                 `json:"message_key"`
	}

	ResponeDat struct {
		ErrorCode int         `json:"error_code"`
		Data      interface{} `json:"data"`
		Msg       string      `json:"msg"`
		MessageId string      `json:"message_id"`
	}

	UserData struct {
		Uid      string `json:"uid"`
		Gender   string `json:"gender"`
		NickName string `json:"nick_name"`
		Avatar   string `json:"avatar"`
		Brithday string `json:"brithday"`
		Ip       string `json:"ip"`
	}

	PfError struct {
		Msg        string //错误描述
		Event      string //触发世界
		Error_Time int    //错误触发时间
	}

	WSDat struct {
		UserData
		GameId string
	}

	Gmresult struct {
		Uid       string
		Score     int
		GameId    string
		MessageID string
	}

	UserGameResult struct {
		NickName string `json:"nick_name"`
		Avatar   string `json:"avatar"`
		PlayNum  int    `json:"play_num"`
		WinNum   int    `json:"win_num"`
		WinPoint int    `json:"win_point"`
	}

	//卡券的信息
	Cards struct {
		Name  string  `json:"name"`
		Id     int    `json:"id"`
		Num    int    `json:"num"`
	}
)

//命令转化服
const (
	START          = "af01"
	LOGIN          = "af02"
	LOGOUT         = "af03"
	CREATE_ROOM    = "af04"
	SEARCH_MATCH   = "af05"
	GAME_HEART     = "af06"
	JOIN_CANCEL    = "af07"
	ROOM_MESSAGE   = "af08"
	OUT_ROOM       = "af09"
	RECONNECT      = "af10"
	NOW_ONLINE_NUM = "af11"
	JOIN_ROOM      = "af12"
	GAME_RESULT    = "af13"
	AUTHORIZE      = "af14"
	TIME_OUT       = "af15"
	DISCONNECT     = "af16"
	ONLINE         = "af17"
	USER_MESSAGE   = "af18"
	ENTER_GAME     = "af19"

	YOU_WIN   	   = "af20"
	GAME_AI        = "af21"
	AI_REPLAY      = "af22"
	CONFIRM   	   = "af23"
	START_REDAY    = "af24"
	REDAY		   = "af25"


	ONLINE_KEY     = "ONE_LINE:%s"

	//win|lose|draw
	WIN  = "win"
	LOSE = "lose"
	DRAW = "draw"
	ROOM_TYPE_SHARE  = "create"
	ROOM_TYPE_REPLAY = "replay"

	ACCESS_TOKEN     = "pk_access_token"
)


