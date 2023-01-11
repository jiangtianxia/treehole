package utils

import (
	"net/http"
	"treehole/define"

	"github.com/gin-gonic/gin"
)

/**
 * @Author jiang
 * @Description 返回结构包
 * @Date 11:00 2023/1/8
 **/
type H struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Resp(ctx *gin.Context, code int, msg string, data interface{}) {
	h := H{
		Code: code,
		Msg:  msg,
		Data: data,
	}

	ctx.JSON(http.StatusOK, h)
}

func RespSuccess(ctx *gin.Context, msg string, data interface{}) {
	Resp(ctx, int(define.SuccessCode), msg, data)
}

func RespFail(ctx *gin.Context, code int, msg string) {
	Resp(ctx, code, msg, "")
}

/**
 * @Author jiang
 * @Description 用户信息返回参数
 * @Date 22:00 2023/1/9
 **/
type UserInfo struct {
	Username string `json:"username"` // 用户名
	Usericon string `json:"usericon"` // 头像
	Age      int    `json:"age"`      // 年龄
	Sex      int    `json:"sex"`      // 性别  0：无性别 1：男 2：女
}

// /**
//  * @Author jiang
//  * @Description 用户信息返回参数
//  * @Date 22:00 2023/1/9
//  **/
// type NoteInfo struct {
// 	AuthorIdentity string `json:"author_identity"`
// 	NoteIdentity   string `json:"note_identity"`
// 	Title          string `json:"title"`
// 	Conetnt        string `json:"content"`
// 	Urls           string `json:"urls"`
// 	CreateTime     string `json:"create_time"`
// 	Visit          int    `json:"visit"`
// }
