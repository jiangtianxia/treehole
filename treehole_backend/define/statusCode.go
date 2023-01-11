package define

/**
 * @Author jiang
 * @Description 业务状态码
 * @Date 11:00 2023/1/8
 **/

const (
	OneWeekInSeconds         = 7 * 24 * 3600
	VoteScore        float64 = 500 // 每一票的值500分
)

var DefaultSize = "20" // 默认大小
var DefaultPage = "1"  // 默认页数

type MyCode int64

const (
	SuccessCode         MyCode = 1000
	FailCode            MyCode = 1001
	ParamsInvalidCode   MyCode = 1002
	EmailExistCode      MyCode = 1003
	UserExistCode       MyCode = 1004
	InvalidPasswordCode MyCode = 1005
	EmailNotExistCode   MyCode = 1006
	InvalidTokenCode    MyCode = 1007
	ExpiredTokenCode    MyCode = 1008
	ServerBusyCode      MyCode = 1009
)

var msgFlags = map[MyCode]string{
	SuccessCode:         "成功",
	FailCode:            "失败",
	ParamsInvalidCode:   "请求参数错误",
	EmailExistCode:      "该邮箱已被注册",
	UserExistCode:       "用户已存在",
	InvalidPasswordCode: "用户名或密码错误",
	EmailNotExistCode:   "该邮箱不存在",
	InvalidTokenCode:    "无效的Token",
	ExpiredTokenCode:    "Token已过期",
	ServerBusyCode:      "服务繁忙",
}

func (c MyCode) Msg() string {
	msg, ok := msgFlags[c]
	if ok {
		return msg
	}
	return msgFlags[ServerBusyCode]
}
