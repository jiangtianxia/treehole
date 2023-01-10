package dao

import (
	"time"
	"treehole/models"
	"treehole/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/spf13/viper"
)

const (
	OneWeekInSeconds         = 7 * 24 * 3600
	VoteScore        float64 = 500 // 每一票的值500分
	PostPerAge               = 20
)

/**
 * @Author jiang
 * @Description redis创建帖子信息
 * @Date 17:00 2023/1/9
 **/
func RedisCreateNote(ctx *gin.Context, userid string, username string, usericon string, noteBasic models.NoteBasic) error {
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
	NoteInfoKey := noteBasic.Identity

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
	pipeline.ZAdd(ctx, VotedKey, redis.Z{
		Score:  1,
		Member: userid,
	})
	pipeline.Expire(ctx, VotedKey, time.Second*OneWeekInSeconds)

	// 3）添加到分数的Zset
	pipeline.ZAdd(ctx, ScoreKey, redis.Z{
		Score:  float64(noteBasic.Score),
		Member: noteBasic.Identity,
	})

	// 4）添加帖子详细信息到Hash
	noteInfo := map[string]interface{}{
		"user_identity": userid,
		"name":          username,
		"icon":          usericon,
		"note_identity": noteBasic.Identity,
		"title":         noteBasic.Title,
		"urls":          noteBasic.Urls,
		"score":         noteBasic.Score,
	}
	pipeline.HMSet(ctx, NoteInfoKey, noteInfo)

	// 提交事务
	_, err := pipeline.Exec(ctx)
	// fmt.Println(utils.RDB.HGetAll(ctx, NoteInfoKey).Val())
	return err
}
