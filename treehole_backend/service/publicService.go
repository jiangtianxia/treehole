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

// SearchNotes
// @Summary 搜索帖子
// @Tags 公共接口
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Success 200 {object} utils.H
// @Router /searchNotes [get]
func SearchNotes(c *gin.Context) {
	// 1、获取参数
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, _ := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	page = (page - 1) * size
	var count int64
	keyword := c.Query("keyword")

	// 2、查询信息
	tx := dao.SearchNotes(keyword)
	list := []*models.NoteBasic{}
	err := tx.Count(&count).Omit("content").Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		logger.SugarLogger.Error("Search Notes Error:" + err.Error())
		utils.RespFail(c, int(define.FailCode), "查询失败")
		return
	}

	// 3、循环遍历，帖子信息，获取作者id等
	data := map[string]interface{}{
		"count": count,
	}
	noteList := make([]map[string]string, 0)
	for _, note := range list {
		author, err := dao.FindUserNoteByNoteIdentity(note.Identity)
		if err != nil {
			logger.SugarLogger.Error("FindUserNoteByNote Identity:" + err.Error())
			utils.RespFail(c, int(define.FailCode), "查询失败")
			return
		}
		temp := map[string]string{
			"author_identity": author.AuthorIdentity,
			"note_identity":   note.Identity,
			"title":           note.Title,
			"content":         note.Content,
			"urls":            note.Urls,
			"create_time":     note.CreateTime,
			"visit":           strconv.Itoa(note.Visit),
			"approve":         strconv.Itoa(note.Approve),
			"against":         strconv.Itoa(note.Against),
		}
		noteList = append(noteList, temp)
	}
	data["noteInfo"] = noteList

	utils.RespSuccess(c, "查询成功", data)
}

// GetNoteInfo
// @Summary 获取帖子详细信息
// @Tags 公共接口
// @Accept application/json
// @Produce application/json
// @Param object body utils.GetNoteInfoFrom true "发送参数"
// @Success 200 {object} utils.H
// @Router /getNoteInfo [post]
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

	NoteInfo["author_icon"] = authorInfo.Usericon
	NoteInfo["author_name"] = authorInfo.Username
	NoteInfo["author_sex"] = strconv.Itoa(authorInfo.Sex)
	NoteInfo["visit"] = strconv.Itoa(visit + 1)

	// 修改redis缓存量
	utils.RDB.HIncrBy(c, NoteInfoKey, "visit", 1)
	utils.RespSuccess(c, "获取帖子详细信息成功", NoteInfo)
}

// SearchNotesScoreOrTime
// @Summary 按照热度或时间获取帖子信息
// @Tags 公共接口
// @Param page query int false "page"
// @Param size query int false "size"
// @Param type query int true "type"
// @Success 200 {object} utils.H
// @Router /searchNotesScoreOrTime [get]
func SearchNotesScoreOrTime(c *gin.Context) {
	// 1、获取参数
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, _ := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	t, _ := strconv.Atoi(c.DefaultQuery("type", "1"))
	page = (page - 1) * size

	if page+size > 5000 {
		utils.RespFail(c, int(define.FailCode), "只能查询排名前5000的数据")
	}

	// 2、到redis中分数的有序集合中获取帖子id
	var noteidList []string
	var count int64
	if t == 1 {
		ScoreKey := viper.GetString("redis.KeyNoteScoreZSetPrefix")
		noteidList = utils.RDB.ZRevRange(c, ScoreKey, int64(page), (int64(page) + int64(size) - 1)).Val()
		count = utils.RDB.ZCard(c, ScoreKey).Val()
	} else {
		TimeKey := viper.GetString("redis.KeyNoteTimeZSetPrefix")
		noteidList = utils.RDB.ZRevRange(c, TimeKey, int64(page), (int64(page) + int64(size) - 1)).Val()
		count = utils.RDB.ZCard(c, TimeKey).Val()
	}

	noteList := make([]map[string]string, 0)
	// 3、for遍历帖子列表，获取哈希集合中的信息
	for _, noteid := range noteidList {
		NoteInfoKey := viper.GetString("redis.KeyNoteInfoHashPrefix") + noteid

		noteInfo := utils.RDB.HGetAll(c, NoteInfoKey).Val()
		noteList = append(noteList, noteInfo)
	}

	data := map[string]interface{}{
		"count":    count,
		"noteInfo": noteList,
	}
	utils.RespSuccess(c, "获取帖子列表成功", data)
}
