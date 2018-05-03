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
)
