package dao

import (
	"time"
	"treehole/define"
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

	// 2）添加分数到投票的Zset，作者默认投赞成表
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
