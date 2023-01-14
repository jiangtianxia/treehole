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
	"github.com/spf13/viper"
)

// CreateComment
// @Summary 发送评论
// @Tags 评论业务接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Authorization"
// @Param object body utils.CreateCommentFrom true "发送参数"
// @Success 200 {object} utils.H
// @Router /comment/createComment [post]
func CreateComment(c *gin.Context) {
	// 1、获取参数
	var form utils.CreateCommentFrom
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

	// 3、将数据信息插入数据库
	identity, err := utils.GetID()
	if err != nil {
		logger.SugarLogger.Error("Get Identity Error" + err.Error())
		utils.RespFail(c, int(define.FailCode), "发送评论失败")
		return
	}

	// 判断该帖子是否存在
	n, err := dao.GetNoteList(form.NoteIdentity)
	if err != nil || n.Identity == "" {
		logger.SugarLogger.Error("GetNoteList Error" + err.Error())
		utils.RespFail(c, int(define.FailCode), "发送评论失败")
		return
	}

	tx := utils.DB.Begin()
	//事务一旦开始，不论什么异常最终都会 Rollback
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	comment := models.CommentNote{
		Identity:     strconv.Itoa(int(identity)),
		UserIdentity: userClaim.Identity,
		NoteIdentity: form.NoteIdentity,
		Conetent:     form.Content,
		CreateTime:   time.Now().Format("2006-01-02 15:04"),
	}

	err = dao.CreateComment(comment)
	if err != nil {
		logger.SugarLogger.Error("Create Comment Error" + err.Error())
		utils.RespFail(c, int(define.FailCode), "发送评论失败")
		tx.Rollback()
		return
	}

	// 到数据库查询用户的信息
	userInfo, err := dao.FindByIdentity(userClaim.Identity)
	if err != nil {
		logger.SugarLogger.Error("FindByIdentity Error" + err.Error())
		utils.RespFail(c, int(define.FailCode), "发送评论失败")
		tx.Rollback()
		return
	}

	// 4、将数据存入redis
	commentInfo := map[string]string{
		"comment_identity": comment.Identity,
		"user_identity":    userInfo.Identity,
		"note_identity":    comment.NoteIdentity,
		"usericon":         userInfo.Usericon,
		"username":         userInfo.Username,
		"content":          comment.Conetent,
		"creare_time":      comment.CreateTime,
	}
	err = dao.RedisCreateCommentNote(c, commentInfo, true)
	if err != nil {
		logger.SugarLogger.Error("RedisCreateCommentNote Error" + err.Error())
		utils.RespFail(c, int(define.FailCode), "发送评论失败")
		tx.Rollback()
		return
	}

	tx.Commit()

	utils.RespSuccess(c, "发送评论成功", commentInfo)
}

// GetNoteCommentList
// @Summary 获取帖子评论列表
// @Tags 评论业务接口
// @Param Authorization header string true "Authorization"
// @Param page query int false "page"
// @Param size query int false "size"
// @Param note_identity query string true "node_identity"
// @Success 200 {object} utils.H
// @Router /comment/getNoteCommentList [get]
func GetNoteCommentList(c *gin.Context) {
	// 1、获取参数
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, _ := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	page = (page - 1) * size
	var count int64
	count = 0
	note_identity := c.DefaultQuery("note_identity", "")

	if note_identity == "" {
		logger.SugarLogger.Error("Params Invalid note_identity")
		utils.RespFail(c, int(define.ParamsInvalidCode), define.ParamsInvalidCode.Msg())
		return
	}

	// 判断该帖子是否存在
	n, err := dao.GetNoteList(note_identity)
	if err != nil || n.Identity == "" {
		logger.SugarLogger.Error("GetNoteList Error" + err.Error())
		utils.RespFail(c, int(define.FailCode), "获取帖子评论失败")
		return
	}

	// 2、判断redis中是否存在该key
	// exisit key，不存在，返回0
	ListKey := viper.GetString("redis.KeyCommentListPrefix") + note_identity

	if utils.RDB.Exists(c, ListKey).Val() == 0 {
		// 不存在，则查询数据库评论，同时将数据存入redis缓存中
		// 查询数据库评论数据
		commentList, err := dao.FindCommentByNoteIdentity(note_identity)
		if err != nil {
			logger.SugarLogger.Error("FindCommentByNoteIdentity Error" + err.Error())
			utils.RespFail(c, int(define.FailCode), "获取帖子评论失败")
			return
		}

		// 遍历列表
		List := []map[string]string{}
		for i, comment := range commentList {
			commentInfo := map[string]string{
				"comment_identity": comment.Identity,
				"user_identity":    comment.UserBasic.Identity,
				"note_identity":    comment.NoteIdentity,
				"usericon":         comment.UserBasic.Usericon,
				"username":         comment.UserBasic.Username,
				"content":          comment.Conetent,
				"creare_time":      comment.CreateTime,
			}

			// 我们需要的分页数据
			if i >= page && i < page+size {
				List = append(List, commentInfo)
			}
			count++

			// 将数据存入redis，并设置过期时间为1天
			dao.RedisCreateCommentNote(c, commentInfo, false)
		}

		data := map[string]interface{}{
			"count":       count,
			"commentInfo": List,
		}
		utils.RespSuccess(c, "获取帖子评论信息成功", data)
	} else {
		// 存在，则读取count，然后根据索引读取出所需的数据，最后查询hash得到全部数据
		// llen key
		// 获取长度
		count = utils.RDB.LLen(c, ListKey).Val()

		// 查询评论identity
		CommentList := utils.RDB.LRange(c, ListKey, int64(page), int64(page+size-1)).Val()

		List := []map[string]string{}
		// 循环遍历查询hash数据
		for _, comment := range CommentList {
			HashKey := viper.GetString("redis.KeyCommentHashPrefix") + comment

			List = append(List, utils.RDB.HGetAll(c, HashKey).Val())
		}

		data := map[string]interface{}{
			"count":       count,
			"commentInfo": List,
		}

		utils.RespSuccess(c, "获取帖子评论信息成功", data)
	}
}

// GetCommentList
// @Summary 获取评论记录
// @Tags 评论业务接口
// @Param Authorization header string true "Authorization"
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {object} utils.H
// @Router /comment/getCommentList [get]
func GetCommentList(c *gin.Context) {
	// 1、获取参数
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, _ := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	page = (page - 1) * size
	var count int

	// 2、获取token上的用户信息
	u, _ := c.Get("user_claims")
	userClaim := u.(*utils.UserClaims)

	// 3、查询数据库上的评论信息
	CommentList, err := dao.FindCommentByUserIdentity(userClaim.Identity)
	if err != nil {
		logger.SugarLogger.Error("FindCommentByUserIdentity Error：", err.Error())
		utils.RespFail(c, int(define.FailCode), "获取评论记录失败")
		return
	}

	count = len(CommentList)
	max := page + size
	if page+size > count {
		max = count
	}
	List := CommentList[page:max]

	data := map[string]interface{}{
		"count":       count,
		"commentInfo": List,
	}
	utils.RespSuccess(c, "获取评论记录成功", data)
}

// DeleteNoteComment
// @Summary 删除评论
// @Tags 评论业务接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Authorization"
// @Param object body utils.DeleteNoteCommentFrom true "发送参数"
// @Success 200 {object} utils.H
// @Router /comment/deleteNoteComment [post]
func DeleteNoteComment(c *gin.Context) {
	// 1、获取参数
	var form utils.DeleteNoteCommentFrom
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

	// 2、判断该帖子是否存在
	n, err := dao.GetNoteList(form.NoteIdentity)
	if err != nil || n.Identity == "" {
		logger.SugarLogger.Error("GetNoteList Error" + err.Error())
		utils.RespFail(c, int(define.FailCode), "删除评论失败")
		return
	}

	// 3、删除缓存数据
	err = dao.RedisDeleteComment(c, form.CommentIdentity, form.NoteIdentity)
	if err != nil {
		logger.SugarLogger.Error("RedisDeleteComment Error：" + err.Error())
		utils.RespFail(c, int(define.FailCode), "删除评论失败")
		return
	}

	// 4、删除数据库数据
	err = dao.DeleteNoteComment(form.CommentIdentity, form.NoteIdentity)
	if err != nil {
		logger.SugarLogger.Error("RedisDeleteComment Error：" + err.Error())
		utils.RespFail(c, int(define.FailCode), "删除评论失败")
		return
	}

	// 5、再次删除数据库数据
	time.Sleep(time.Second * 3)
	err = dao.RedisDeleteComment(c, form.CommentIdentity, form.NoteIdentity)
	if err != nil {
		logger.SugarLogger.Error("RedisDeleteComment Error：" + err.Error())
		utils.RespFail(c, int(define.FailCode), "删除评论失败")
		return
	}

	utils.RespSuccess(c, "删除评论成功", "")
}
