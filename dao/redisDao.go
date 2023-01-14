package dao

import (
	"strconv"
	"time"
	"treehole/define"
	"treehole/logger"
	"treehole/models"
	"treehole/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/spf13/viper"
)

/**
 * @Author jiang
 * @Description redis创建帖子信息
 * @Date 17:00 2023/1/9
 **/
func RedisCreateNote(ctx *gin.Context, userid string, username string, noteBasic models.NoteBasic) error {
	// 使用事务操作
	// 共四个表结构
	// 有序集合
	// 1、note:time
	// 2、note:voted:noteid
	// 3、note:score
	// 哈希
	// 1、key：noteid  value：noteInfo
	TimeKey := viper.GetString("redis.KeyNoteTimeZSetPrefix")
	VotedKey := viper.GetString("redis.KeyNoteVotedZSetPrefix") + noteBasic.Identity
	ScoreKey := viper.GetString("redis.KeyNoteScoreZSetPrefix")
	NoteInfoKey := viper.GetString("redis.KeyNoteInfoHashPrefix") + noteBasic.Identity

	// fmt.Println("timeKey" + TimeKey)
	// fmt.Println("VotedKey" + VotedKey)
	// fmt.Println("ScoreKey" + ScoreKey)
	// fmt.Println("NoteInfoKey" + NoteInfoKey)

	// 事务操作
	pipeline := utils.RDB.TxPipeline()

	// 1）添加到时间的Zset，不设置过期时间，用于获得总文章数
	pipeline.ZAdd(ctx, TimeKey, redis.Z{
		Score:  float64(noteBasic.Score),
		Member: noteBasic.Identity,
	})

	// 2）添加分数到投票的Zset，作者默认点赞
	// 设置过期时间，等到用户点赞时，再判断是否需要删除
	pipeline.ZAdd(ctx, VotedKey, redis.Z{
		Score:  1,
		Member: userid,
	})
	pipeline.Expire(ctx, VotedKey, time.Second*define.OneWeekInSeconds)

	// 3）添加到分数的Zset
	pipeline.ZAdd(ctx, ScoreKey, redis.Z{
		Score:  float64(noteBasic.Score),
		Member: noteBasic.Identity,
	})

	// 4）添加帖子详细信息到Hash
	noteInfo := map[string]interface{}{
		"author_identity": userid,
		"note_identity":   noteBasic.Identity,
		"title":           noteBasic.Title,
		"content":         noteBasic.Content,
		"urls":            noteBasic.Urls,
		"create_time":     noteBasic.CreateTime,
		"visit":           noteBasic.Visit,
		"approve":         noteBasic.Approve,
		"against":         noteBasic.Against,
	}
	pipeline.HMSet(ctx, NoteInfoKey, noteInfo)

	// 提交事务
	_, err := pipeline.Exec(ctx)
	// fmt.Println(utils.RDB.HGetAll(ctx, NoteInfoKey).Val())
	return err
}

/**
 * @Author jiang
 * @Description 删除redis中哈希集合与时间和分数有序集合的信息
 * @Date 16:00 2023/1/11
 **/
func RedisDeleteNote(c *gin.Context, noteid string) error {
	// 1、获取全部的key
	TimeKey := viper.GetString("redis.KeyNoteTimeZSetPrefix")
	VotedKey := viper.GetString("redis.KeyNoteVotedZSetPrefix") + noteid
	ScoreKey := viper.GetString("redis.KeyNoteScoreZSetPrefix")
	NoteInfoKey := viper.GetString("redis.KeyNoteInfoHashPrefix") + noteid

	// // 打印key，调试
	// fmt.Println("timeKey" + TimeKey)
	// fmt.Println("VotedKey" + VotedKey)
	// fmt.Println("ScoreKey" + ScoreKey)
	// fmt.Println("NoteInfoKey" + NoteInfoKey)

	// 2、删除信息
	// 事务操作
	pipeline := utils.RDB.TxPipeline()

	// 1）删除哈希集合的key
	pipeline.Del(c, NoteInfoKey)

	// 2）删除有序集合中的成员投票的key
	pipeline.Del(c, VotedKey)

	// 3）删除时间集合中的成员
	pipeline.ZRem(c, TimeKey, noteid)

	// 4）删除分数有序集合中的成员
	pipeline.ZRem(c, ScoreKey, noteid)

	// 提交事务
	_, err := pipeline.Exec(c)

	// // 获取数据，验证是否删除成功
	// fmt.Println(utils.RDB.HGetAll(c, NoteInfoKey).Val())
	return err
}

// 判断是否存在该投票key值
func RedisVotedExists(c *gin.Context, note_identity string) int64 {
	VotedKey := viper.GetString("redis.KeyNoteVotedZSetPrefix") + note_identity

	return utils.RDB.Exists(c, VotedKey).Val()
}

/**
 * @Author jiang
 * @Description 查询分数排名是否在2000名以内，不在则删除成员和哈希key，修改mysql中的投票数，不计分数
 * @Date 13:00 2023/1/12
 **/
func RedisScoreNotExists(c *gin.Context, form utils.VotedNoteFrom, user utils.UserClaims) error {
	TimeKey := viper.GetString("redis.KeyNoteTimeZSetPrefix")
	ScoreKey := viper.GetString("redis.KeyNoteScoreZSetPrefix")
	NoteInfoKey := viper.GetString("redis.KeyNoteInfoHashPrefix") + form.NoteIdentity
	toVoted, _ := strconv.Atoi(form.Flag)
	Voted, _ := strconv.Atoi(form.Voted)

	// 1、读取分数排名
	flag := utils.RDB.ZRevRank(c, ScoreKey, form.NoteIdentity).Val()
	if flag >= 2000 {
		// 不在，则删除成员和哈希key
		// 事务操作
		pipeline := utils.RDB.TxPipeline()

		// 1）删除哈希集合的key
		pipeline.Del(c, NoteInfoKey)

		// 2）删除分数有序集合中的成员
		pipeline.ZRem(c, ScoreKey, form.NoteIdentity)

		// 提交事务
		_, err := pipeline.Exec(c)

		if err != nil {
			return err
		}
	} else {
		// 存在，则修改哈希集合
		// -1 -> 1
		// -1 -> 0
		// 1 -> 0
		// 1 -> -1
		// 事务操作
		pipeline := utils.RDB.TxPipeline()

		if Voted == -1 {
			pipeline.HIncrBy(c, NoteInfoKey, "against", -1)
		}

		if Voted == 1 {
			pipeline.HIncrBy(c, NoteInfoKey, "approve", -1)
		}

		if toVoted == 1 {
			pipeline.HIncrBy(c, NoteInfoKey, "approve", 1)
		}

		if toVoted == -1 {
			pipeline.HIncrBy(c, NoteInfoKey, "against", 1)
		}

		// 提交事务
		_, err := pipeline.Exec(c)

		if err != nil {
			return err
		}
	}

	// 2、判断时间排名是否在2000名以内
	flag = utils.RDB.ZRevRank(c, TimeKey, form.NoteIdentity).Val()
	if flag >= 2000 {
		err := utils.RDB.ZRem(c, TimeKey, form.NoteIdentity).Err()
		if err != nil {
			return err
		}
	}

	// 3、修改数据库点赞数
	// 判断点赞表中是否存在该记录，不存在，则创建，存在则修改
	tx := utils.DB.Begin()
	//事务一旦开始，不论什么异常最终都会 Rollback
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	voted := FindVotedNodeByIdentity(user.Identity, form.NoteIdentity)
	isvoted := form.Flag
	if voted.NoteIdentity == "" {
		// 不存在，则创建
		v := models.VotedNote{
			NoteIdentity: form.NoteIdentity,
			UserIdentity: user.Identity,
			Isvoted:      isvoted,
		}
		err := CreateVotedNote(v)
		if err != nil {
			logger.SugarLogger.Error("CreateVotedNote Error：" + err.Error())
			tx.Rollback()
			return err
		}
	} else {
		// 存在，则修改
		v := models.VotedNote{
			NoteIdentity: form.NoteIdentity,
			UserIdentity: user.Identity,
			Isvoted:      isvoted,
		}
		err := ModifyVotedNote(v)
		if err != nil {
			logger.SugarLogger.Error("ModifyVotedNote Error：" + err.Error())
			tx.Rollback()
			return err
		}
	}

	// 修改帖子点赞和踩数
	vo, err := GetNoteList(form.NoteIdentity)
	if err != nil {
		logger.SugarLogger.Error("GetNoteList Error：" + err.Error())
		tx.Rollback()
		return err
	}
	if Voted == -1 {
		i, _ := strconv.Atoi(vo.Against)
		vo.Against = strconv.Itoa(i - 1)
	}

	if Voted == 1 {
		i, _ := strconv.Atoi(vo.Approve)
		vo.Approve = strconv.Itoa(i - 1)
	}

	if toVoted == 1 {
		i, _ := strconv.Atoi(vo.Approve)
		vo.Approve = strconv.Itoa(i + 1)
	}

	if toVoted == -1 {
		i, _ := strconv.Atoi(vo.Against)
		vo.Against = strconv.Itoa(i - 1)
	}

	// 修改
	err = UpdateNote(vo)
	if err != nil {
		logger.SugarLogger.Error("UpdateNote Error：" + err.Error())
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

/**
 * @Author jiang
 * @Description 修改投票数据，再修改分数，最后修改数据库
 * @Date 13:00 2023/1/12
 **/
func RedisScoreExists(c *gin.Context, form utils.VotedNoteFrom, user utils.UserClaims) error {
	ScoreKey := viper.GetString("redis.KeyNoteScoreZSetPrefix")
	NoteInfoKey := viper.GetString("redis.KeyNoteInfoHashPrefix") + form.NoteIdentity
	VotedKey := viper.GetString("redis.KeyNoteVotedZSetPrefix") + form.NoteIdentity
	toVoted, _ := strconv.Atoi(form.Flag)
	Voted, _ := strconv.Atoi(form.Voted)
	incr := int64(toVoted - Voted)

	// 事务操作
	pipeline := utils.RDB.TxPipeline()

	// 修改分数
	pipeline.ZIncrBy(c, ScoreKey, float64(incr)*define.VoteScore, form.NoteIdentity)

	// 修改投票数
	pipeline.ZIncrBy(c, VotedKey, float64(incr), user.Identity)

	// 修改哈希集合中的反对和赞成
	if Voted == -1 {
		pipeline.HIncrBy(c, NoteInfoKey, "against", -1)
	}

	if Voted == 1 {
		pipeline.HIncrBy(c, NoteInfoKey, "approve", -1)
	}

	if toVoted == 1 {
		pipeline.HIncrBy(c, NoteInfoKey, "approve", 1)
	}

	if toVoted == -1 {
		pipeline.HIncrBy(c, NoteInfoKey, "against", 1)
	}

	// 提交事务
	_, err := pipeline.Exec(c)

	if err != nil {
		return err
	}

	// 修改数据库数据
	tx := utils.DB.Begin()
	//事务一旦开始，不论什么异常最终都会 Rollback
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	voted := FindVotedNodeByIdentity(user.Identity, form.NoteIdentity)
	isvoted := form.Flag
	if voted.NoteIdentity == "" {
		// 不存在，则创建
		v := models.VotedNote{
			NoteIdentity: form.NoteIdentity,
			UserIdentity: user.Identity,
			Isvoted:      isvoted,
		}
		err := CreateVotedNote(v)
		if err != nil {
			logger.SugarLogger.Error("CreateVotedNote Error：" + err.Error())
			tx.Rollback()
			return err
		}
	} else {
		// 存在，则修改
		v := models.VotedNote{
			UserIdentity: user.Identity,
			NoteIdentity: form.NoteIdentity,
			Isvoted:      isvoted,
		}

		err := ModifyVotedNote(v)
		if err != nil {
			logger.SugarLogger.Error("ModifyVotedNote Error：" + err.Error())
			tx.Rollback()
			return err
		}
	}

	// 修改帖子点赞、踩数和分数
	vo, err := GetNoteList(form.NoteIdentity)
	if err != nil {
		logger.SugarLogger.Error("GetNoteList Error：" + err.Error())
		tx.Rollback()
		return err
	}
	if Voted == -1 {
		i, _ := strconv.Atoi(vo.Against)
		vo.Against = strconv.Itoa(i - 1)
	}

	if Voted == 1 {
		i, _ := strconv.Atoi(vo.Approve)
		vo.Approve = strconv.Itoa(i - 1)
	}

	if toVoted == 1 {
		i, _ := strconv.Atoi(vo.Approve)
		vo.Approve = strconv.Itoa(i + 1)
	}

	if toVoted == -1 {
		i, _ := strconv.Atoi(vo.Against)
		vo.Against = strconv.Itoa(i - 1)
	}
	vo.Score = vo.Score + int(incr)*int(define.VoteScore)

	// 修改数据库
	err = UpdateNote(vo)
	if err != nil {
		logger.SugarLogger.Error("UpdateNote Error：" + err.Error())
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
