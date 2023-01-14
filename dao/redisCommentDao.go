package dao

import (
	"time"
	"treehole/define"
	"treehole/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

/**
 * @Author jiang
 * @Description 创建评论缓存
 * @Date 21:00 2023/1/12
 **/
func RedisCreateCommentNote(c *gin.Context, commentInfo map[string]string, flag bool) error {

	// 1、拼接key
	// listKey：treehole:comment:list:note_identity
	// hashKey: treehole:comment:hash:comment_identity
	ListKey := viper.GetString("redis.KeyCommentListPrefix") + commentInfo["note_identity"]
	HashKey := viper.GetString("redis.KeyCommentHashPrefix") + commentInfo["comment_identity"]

	// fmt.Println("ListKey：", ListKey)
	// fmt.Println("HashKey：", HashKey)

	// 事务操作
	pipeline := utils.RDB.TxPipeline()
	// 2、在list前插数据
	// lpush key value
	pipeline.LPush(c, ListKey, commentInfo["comment_identity"])
	if flag {
		pipeline.Expire(c, ListKey, time.Second*define.OneWeekInSeconds)
	} else {
		pipeline.Expire(c, ListKey, time.Second*define.OneWeekInSeconds/7)
	}

	// HMSet key value
	// 3、插入哈希数据
	pipeline.HMSet(c, HashKey, commentInfo)
	if flag {
		pipeline.Expire(c, HashKey, time.Second*define.OneWeekInSeconds)
	} else {
		pipeline.Expire(c, HashKey, time.Second*define.OneWeekInSeconds/7)
	}

	_, err := pipeline.Exec(c)

	// fmt.Println(utils.RDB.LRange(c, ListKey, 0, -1).Val())
	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>")
	// fmt.Println(utils.RDB.HGetAll(c, HashKey).Val())
	return err
}

/**
 * @Author jiang
 * @Description 删除评论缓存
 * @Date 23:30 2023/1/12
 **/
func RedisDeleteComment(c *gin.Context, comment_identity string, note_identity string) error {
	// 1、拼接key
	// listKey：treehole:comment:list:note_identity
	// hashKey: treehole:comment:hash:comment_identity
	ListKey := viper.GetString("redis.KeyCommentListPrefix") + note_identity
	HashKey := viper.GetString("redis.KeyCommentHashPrefix") + comment_identity

	// 2、删除缓存
	// del key
	// lrem key count VALUE
	// 事务操作
	pipeline := utils.RDB.TxPipeline()
	pipeline.Del(c, HashKey)
	pipeline.LRem(c, ListKey, 1, comment_identity)

	_, err := pipeline.Exec(c)
	return err
}
