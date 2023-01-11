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
		Identity:   strconv.Itoa(int(noteid)),
		Title:      form.Title,
		Content:    form.Content,
		Urls:       form.Urls,
		Score:      int(nowTime),
		Approve:    1,
		Against:    0,
		Visit:      1,
		CreateTime: time.Now().Format("2006-01-02 15:04"),
	}
	err = dao.CreateNote(noteBasic)
	if err != nil {
		logger.SugarLogger.Error("create note Error:" + err.Error())
		tx.Rollback()
		utils.RespFail(c, int(define.FailCode), "发布帖子失败")
		return
	}

	userNote := models.UserNote{
		AuthorIdentity: userClaim.Identity,
		NoteIdentity:   strconv.Itoa(int(noteid)),
	}
	err = dao.CreateUserNode(userNote)
	if err != nil {
		logger.SugarLogger.Error("create usernote Error:" + err.Error())
		tx.Rollback()
		utils.RespFail(c, int(define.FailCode), "发布帖子失败")
		return
	}

	// 5、将数据库存入redis。
	err = dao.RedisCreateNote(c, userClaim.Identity, userClaim.Username, noteBasic)
	if err != nil {
		logger.SugarLogger.Error("Redis Create Note Error:" + err.Error())
		utils.RespFail(c, int(define.FailCode), "发布帖子失败")
		tx.Rollback()
		return
	}

	tx.Commit()

	// 返回信息
	utils.RespSuccess(c, "发布帖子成功", "")
}

// GetNoteList
// @Summary 获取发布帖子列表
// @Tags 帖子业务接口
// @Param Authorization header string true "Authorization"
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {object} utils.H
// @Router /note/getNoteList [get]
func GetNoteList(c *gin.Context) {
	// 1、获取参数
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, _ := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	page = (page - 1) * size
	var count int64

	// 2、根据identidy查询帖子关系表信息
	u, _ := c.Get("user_claims")
	userClaim := u.(*utils.UserClaims)

	tx := dao.FindUserNoteByAuthorIdentity(userClaim.Identity)

	list := []models.UserNote{}
	err := tx.Count(&count).Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		logger.SugarLogger.Error("FindUserNoteByAuthorIdentity Error:" + err.Error())
		utils.RespFail(c, int(define.FailCode), "查询失败")
		return
	}

	data := make(map[string]interface{}, 0)
	data["count"] = count
	data["noteInfo"] = list
	utils.RespSuccess(c, "查询成功", data)
}

// DeleteNote
// @Summary 删除帖子
// @Tags 帖子业务接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Authorization"
// @Param object body utils.GetNoteInfoFrom true "发送参数"
// @Success 200 {object} utils.H
// @Router /note/deleteNote [post]
func DeleteNote(c *gin.Context) {
	// 1、获取参数
	var form utils.GetNoteInfoFrom
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

	// 2、删除redis数据
	err := dao.RedisDeleteNote(c, form.NoteIdentity)
	if err != nil {
		logger.SugarLogger.Error("RedisDeleteNote Error" + err.Error())
		utils.RespFail(c, int(define.FailCode), "删除帖子失败")
		return
	}

	// 3、删除数据库及本地图片信息
	err = DeleteImageAndMysql(form.NoteIdentity)
	if err != nil {
		logger.SugarLogger.Error("DeleteImageAndMysql Error" + err.Error())
		utils.RespFail(c, int(define.FailCode), "删除帖子失败")
		return
	}

	utils.RespSuccess(c, "删除帖子成功", "")
}
