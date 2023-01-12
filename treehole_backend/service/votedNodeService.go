package service

import (
	"treehole/dao"
	"treehole/define"
	"treehole/logger"
	"treehole/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// VotedNote
// @Summary 点赞或踩帖子
// @Tags 帖子业务接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Authorization"
// @Param object body utils.VotedNoteFrom true "发送参数"
// @Success 200 {object} utils.H
// @Router /note/votedNote [post]
func VotedNote(c *gin.Context) {
	// 1、获取参数
	var form utils.VotedNoteFrom
	if err := c.ShouldBindJSON(&form); err != nil {
		// 获取valadtor.valiadtionErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非valiadtor.ValidationErrors类型错误直接返回
			utils.RespFail(c, int(define.ParamsInvalidCode), err.Error())
			logger.SugarLogger.Error("Params Invalid" + err.Error())
			return
		}
		// validator.ValidationErrors类型错误进行翻译
		hashMap := utils.RemoveTopStruct(errs.Translate(utils.Trans))
		msg := ""
		for _, v := range hashMap {
			msg += v + ","
		}
		logger.SugarLogger.Error("Params Invalid" + msg)
		utils.RespFail(c, int(define.ParamsInvalidCode), msg)
		return
	}

	// 2、获取token上的用户信息
	u, _ := c.Get("user_claims")
	userClaim := u.(*utils.UserClaims)

	// 3、判断redis中是否存在该投票key
	flag := dao.RedisVotedExists(c, form.NoteIdentity)

	if flag == 0 {
		// 1）不存在，查询分数排名是否在2000名以内，不在则删除成员和哈希key，修改mysql中的投票数，不计分数
		err := dao.RedisScoreNotExists(c, form, *userClaim)
		if err != nil {
			logger.SugarLogger.Error("RedisScoreExists Error:" + err.Error())
			utils.RespFail(c, int(define.FailCode), "点赞或踩失败")
			return
		}
	} else {
		// 2）存在，修改投票数据，再修改分数，最后修改数据库
		err := dao.RedisScoreExists(c, form, *userClaim)
		if err != nil {
			logger.SugarLogger.Error("RedisScoreExists Error:" + err.Error())
			utils.RespFail(c, int(define.FailCode), "点赞或踩失败")
			return
		}
	}

	utils.RespSuccess(c, "点赞或踩成功", "")
}
