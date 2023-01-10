package service

import (
	"strconv"
	"time"
	"treehole/dao"
	"treehole/define"
	"treehole/logger"
	"treehole/models"
	"treehole/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// CresteNote
// @Summary 创建帖子
// @Tags 帖子业务接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Authorization"
// @Param object body utils.CreateNoteForm true "发送参数"
// @Success 200 {object} utils.H
// @Router /note/createNote [post]
func CreateNote(c *gin.Context) {
	// 1、获取参数
	var form utils.CreateNoteForm
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

	// 3、生成帖子identity和获取当前时间
	noteid, err := utils.GetID()
	if err != nil {
		logger.SugarLogger.Error("create identity Error:" + err.Error())
		utils.RespFail(c, int(define.FailCode), "发布帖子失败")
		return
	}

	nowTime := time.Now().Unix() + 2

	// 4、将数据存入数据库
	tx := utils.DB.Begin()
	//事务一旦开始，不论什么异常最终都会 Rollback
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	noteBasic := models.NoteBasic{
		Identity: strconv.Itoa(int(noteid)),
		Title:    form.Title,
		Content:  form.Content,
		Urls:     form.Urls,
		Score:    int(nowTime),
	}
	err = dao.CreateNote(noteBasic)
	if err != nil {
		logger.SugarLogger.Error("create note Error:" + err.Error())
		tx.Rollback()
		utils.RespFail(c, int(define.FailCode), "发布帖子失败")
		return
	}

	userNote := models.UserNote{
		UserIdentity: userClaim.Identity,
		NoteIdentity: strconv.Itoa(int(noteid)),
	}
	err = dao.CreateUserNode(userNote)
	if err != nil {
		logger.SugarLogger.Error("create usernote Error:" + err.Error())
		tx.Rollback()
		utils.RespFail(c, int(define.FailCode), "发布帖子失败")
		return
	}

	// 5、将数据库存入redis。
	// 获取用户的头像url和username
	// uuser, err := dao.FindByIdentity(userClaim.Identity)
	// if err != nil {
	// 	logger.SugarLogger.Error("Find By Identity Error:" + err.Error())
	// 	utils.RespFail(c, int(define.FailCode), "发布帖子失败")
	// 	return
	// }

	err = dao.RedisCreateNote(c, userClaim.Identity, userClaim.Username, userClaim.Usericon, noteBasic)
	if err != nil {
		logger.SugarLogger.Error("Redis Create Note Error:" + err.Error())
		utils.RespFail(c, int(define.FailCode), "发布帖子失败")
		return
	}

	tx.Commit()

	// 返回信息
	utils.RespSuccess(c, "发布帖子成功", "", 0)
}
