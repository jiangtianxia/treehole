package define

/**
 * @Author jiang
 * @Description 业务状态码
 * @Date 11:00 2023/1/8
 **/

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

	CodeUserExist       MyCode = 10019
	CodeUserNotExist    MyCode = 10020
	CodeInvalidPassword MyCode = 10021
	CodeServerBusy      MyCode = 10022

	CodeInvalidToken MyCode = 1023
	// CodeInvalidAuthFormat MyCode = 1007
	// CodeNotLogin MyCode = 1008
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

	CodeUserExist:       "用户名重复",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",

	CodeInvalidToken: "无效的Token",
	// CodeInvalidAuthFormat: "认证格式有误",
	// CodeNotLogin: "未登录",
}

func (c MyCode) Msg() string {
	msg, ok := msgFlags[c]
	if ok {
		return msg
	}
	return msgFlags[CodeServerBusy]
}
