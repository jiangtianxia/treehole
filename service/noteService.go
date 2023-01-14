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
		Approve:    "1",
		Against:    "0",
		Visit:      0,
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

	votedNote := models.VotedNote{
		NoteIdentity: strconv.Itoa(int(noteid)),
		UserIdentity: userClaim.Identity,
		Isvoted:      "1",
	}
	err = dao.CreateVotedNote(votedNote)
	if err != nil {
		logger.SugarLogger.Error("CreateVotedNote Error:" + err.Error())
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

// ModifyNote
// @Summary 修改帖子
// @Tags 帖子业务接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Authorization"
// @Param object body utils.ModifyNoteForm true "发送参数"
// @Success 200 {object} utils.H
// @Router /note/modifyNote [post]
func ModifyNote(c *gin.Context) {
	// 1、获取参数
	var form utils.ModifyNoteForm
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
		utils.RespFail(c, int(define.FailCode), "修改帖子失败")
		return
	}

	// 3、修改数据库数据
	note := models.NoteBasic{
		Identity: form.NoteIdentity,
		Title:    form.Title,
		Urls:     form.Urls,
		Content:  form.Content,
	}
	err = dao.UpdateNote(note)
	if err != nil {
		logger.SugarLogger.Error("UpdateNote Error" + err.Error())
		utils.RespFail(c, int(define.FailCode), "修改帖子失败")
		return
	}

	// 4、休眠3秒钟，再次删除redis缓存
	time.Sleep(time.Second * 3)
	err = dao.RedisDeleteNote(c, form.NoteIdentity)
	if err != nil {
		logger.SugarLogger.Error("RedisDeleteNote Error" + err.Error())
		utils.RespFail(c, int(define.FailCode), "修改帖子失败")
		return
	}

	utils.RespSuccess(c, "修改帖子成功", "")
}

// GetNoteInfo
// @Summary 获取帖子详细信息
// @Tags 帖子业务接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Authorization"
// @Param object body utils.GetNoteInfoFrom true "发送参数"
// @Success 200 {object} utils.H
// @Router /note/getNoteInfo [post]
func GetNoteInfo(c *gin.Context) {
	// 1、获取参数
	var note utils.GetNoteInfoFrom
	if err := c.ShouldBindJSON(&note); err != nil {
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

	// 获取token上的用户信息
	u, _ := c.Get("user_claims")
	userClaim := u.(*utils.UserClaims)

	// 2、到缓存中查询是否存在该帖子数据
	NoteInfoKey := viper.GetString("redis.KeyNoteInfoHashPrefix") + note.NoteIdentity
	NoteInfo := utils.RDB.HGetAll(c, NoteInfoKey).Val()

	if len(NoteInfo) == 0 {
		// 如果帖子中不存在该数据，则去数据库查询。
		// 然后将信息放到redis缓存当中，同时设置过期时间为半天
		n, err := dao.FindUserNoteByNoteIdentityFind(note.NoteIdentity)
		if err != nil {
			logger.SugarLogger.Error("FindUserNoteByNoteIdentityFind Error:" + err.Error())
			utils.RespFail(c, int(define.FailCode), "获取详细信息失败")
			return
		}
		authorInfo, err := dao.FindByIdentity(n.AuthorIdentity)
		if err != nil {
			logger.SugarLogger.Error("FindByIdentity Error:" + err.Error())
			utils.RespFail(c, int(define.FailCode), "获取详细信息失败")
			return
		}

		// 将数据存入redis
		// 事务操作
		pipeline := utils.RDB.TxPipeline()
		noteInfo := map[string]interface{}{
			"author_identity": n.AuthorIdentity,
			"author_name":     authorInfo.Username,
			"author_icon":     authorInfo.Usericon,
			"author_sex":      authorInfo.Sex,
			"note_identity":   n.NoteIdentity,
			"title":           n.NoteBasic.Title,
			"content":         n.NoteBasic.Content,
			"urls":            n.NoteBasic.Urls,
			"create_time":     n.NoteBasic.CreateTime,
			"visit":           n.NoteBasic.Visit + 1,
			"approve":         n.NoteBasic.Approve,
			"against":         n.NoteBasic.Against,
		}
		pipeline.HMSet(c, NoteInfoKey, noteInfo)
		pipeline.Expire(c, NoteInfoKey, time.Second*define.OneWeekInSeconds/14)
		// 提交事务
		pipeline.Exec(c)

		// 查询是否有点赞数据
		voted := dao.FindVotedNodeByIdentity(userClaim.Identity, n.NoteIdentity)
		if voted.UserIdentity == "" {
			votedNote := models.VotedNote{
				NoteIdentity: n.NoteIdentity,
				UserIdentity: userClaim.Identity,
				Isvoted:      "0",
			}
			dao.CreateVotedNote(votedNote)
			noteInfo["isvoted"] = "0"
		} else {
			noteInfo["isvoted"] = voted.Isvoted
		}

		// 修改数据库访问量
		tmp := models.NoteBasic{
			Identity: n.NoteIdentity,
			Visit:    n.NoteBasic.Visit + 1,
		}
		err = dao.UpdateNote(tmp)
		if err != nil {
			logger.SugarLogger.Error("FindByIdentity Error:" + err.Error())
			utils.RespFail(c, int(define.FailCode), "获取详细信息失败")
			return
		}
		utils.RespSuccess(c, "获取帖子详细信息成功", noteInfo)
		return
	}

	// 如果存在数据，则根据author_id去查询作者信息
	authorInfo, err := dao.FindByIdentity(NoteInfo["author_identity"])
	if err != nil {
		logger.SugarLogger.Error("FindByIdentity Error:" + err.Error())
		utils.RespFail(c, int(define.FailCode), "获取详细信息失败")
		return
	}

	// 修改数据库访问量
	visit, _ := strconv.Atoi(NoteInfo["visit"])
	tmp := models.NoteBasic{
		Identity: NoteInfo["note_identity"],
		Visit:    (visit + 1),
	}
	err = dao.UpdateNote(tmp)
	if err != nil {
		logger.SugarLogger.Error("FindByIdentity Error:" + err.Error())
		utils.RespFail(c, int(define.FailCode), "获取详细信息失败")
		return
	}

	// 查询缓存中改用户是否有点赞
	VotedKey := viper.GetString("redis.KeyNoteVotedZSetPrefix") + NoteInfo["note_identity"]
	flag := utils.RDB.ZScore(c, VotedKey, userClaim.Identity).Val()
	if flag == 0 {
		// 到数据库查询数据
		voted := dao.FindVotedNodeByIdentity(userClaim.Identity, NoteInfo["note_identity"])
		if voted.UserIdentity == "" {
			votedNote := models.VotedNote{
				NoteIdentity: NoteInfo["note_identity"],
				UserIdentity: userClaim.Identity,
				Isvoted:      "0",
			}
			dao.CreateVotedNote(votedNote)
			NoteInfo["isvoted"] = "0"
		} else {
			NoteInfo["isvoted"] = voted.Isvoted
		}
	} else {
		NoteInfo["isvoted"] = strconv.Itoa(int(flag))
	}

	NoteInfo["author_icon"] = authorInfo.Usericon
	NoteInfo["author_name"] = authorInfo.Username
	NoteInfo["author_sex"] = strconv.Itoa(authorInfo.Sex)
	NoteInfo["visit"] = strconv.Itoa(visit + 1)

	// 修改redis缓存量
	utils.RDB.HIncrBy(c, NoteInfoKey, "visit", 1)
	utils.RespSuccess(c, "获取帖子详细信息成功", NoteInfo)
}
