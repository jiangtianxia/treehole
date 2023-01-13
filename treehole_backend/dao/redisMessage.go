package dao

import (
	"context"
	"fmt"
	"time"
	"treehole/define"
	"treehole/models"
	"treehole/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

/**
 * @Author jiang
 * @Description 将用户数据加入在线缓存
 * @Date 15:00 2023/1/13
 **/
func RedisMessageOnlineCreate(c *gin.Context, user models.UserBasic) error {
	// 1、拼接key
	// listKey：treehole:online:list
	// hashKey: treehole:online:hash:
	ListKey := viper.GetString("redis.KeyWebsocketOnlineList")
	HashKey := viper.GetString("redis.KeyWebsocketOnlineHashPrefix") + user.Identity

	// fmt.Println("ListKey：", ListKey)
	// fmt.Println("HashKey：", HashKey)

	// 事务操作
	pipeline := utils.RDB.TxPipeline()
	// 2、在list前插数据
	// lpush key value
	pipeline.LRem(c, ListKey, 0, user.Identity)
	pipeline.LPush(c, ListKey, user.Identity)

	// HMSet key value
	// 3、插入哈希数据
	userInfo := map[string]string{
		"useridentity": user.Identity,
		"usericon":     user.Usericon,
		"username":     user.Username,
	}
	pipeline.HMSet(c, HashKey, userInfo)
	pipeline.Expire(c, HashKey, time.Second*define.OneWeekInSeconds/7*2)

	_, err := pipeline.Exec(c)

	// fmt.Println(utils.RDB.LRange(c, ListKey, 0, -1).Val())
	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>")
	// fmt.Println(utils.RDB.HGetAll(c, HashKey).Val())
	return err
}

/**
 * @Author jiang
 * @Description 将聊天记录加入缓存
 * @Date 15:00 2023/1/13
 **/
func RedisCreateMessage(c *gin.Context, message models.MessageBasic) error {
	// 1、拼接key
	// listKey：treehole:online:list
	// hashKey: treehole:online:hash:
	ListKey := viper.GetString("redis.KeyWebsocketMessageList")
	HashKey := viper.GetString("redis.KeyWebsocketMessageHashPrefix") + message.Identity

	// fmt.Println("ListKey：", ListKey)
	// fmt.Println("HashKey：", HashKey)

	// 事务操作
	pipeline := utils.RDB.TxPipeline()
	// 2、在list前插数据
	// lpush key value
	pipeline.LPush(c, ListKey, message.Identity)
	pipeline.LTrim(c, ListKey, 0, 500)

	// HMSet key value
	// 3、插入哈希数据
	messageInfo := map[string]string{
		"useridentity": message.UserIdentity,
		"usericon":     message.Usericon,
		"username":     message.Username,
		"date":         message.Date,
		"create_time":  message.CreateTime,
	}
	pipeline.HMSet(c, HashKey, messageInfo)
	pipeline.Expire(c, HashKey, time.Second*define.OneWeekInSeconds)

	_, err := pipeline.Exec(c)

	// fmt.Println(utils.RDB.LRange(c, ListKey, 0, -1).Val())
	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>")
	// fmt.Println(utils.RDB.HGetAll(c, HashKey).Val())
	return err
}

/**
 * @Author jiang
 * @Description 将数据回存redis缓存中，同时设置过期时间为1天
 * @Date 17:00 2023/1/13
 **/
func RedisSave(c *gin.Context, key string, info map[string]string) error {
	// 事务操作
	pipeline := utils.RDB.TxPipeline()

	// HMSet key value
	// 插入哈希数据
	pipeline.HMSet(c, key, info)
	pipeline.Expire(c, key, time.Second*define.OneWeekInSeconds/7)

	_, err := pipeline.Exec(c)

	fmt.Println(utils.RDB.HGetAll(c, key).Val())
	return err
}

/**
 * @Author jiang
 * @Description 删除redis缓存
 * @Date 23:00 2023/1/13
 **/

func RedisDeleteMessage(identity string) error {
	// 1、拼接key
	// listKey：treehole:online:list
	// hashKey: treehole:online:hash:
	ListKey := viper.GetString("redis.KeyWebsocketOnlineList")
	HashKey := viper.GetString("redis.KeyWebsocketOnlineHashPrefix") + identity

	var ctx = context.Background()
	// 事务操作
	pipeline := utils.RDB.TxPipeline()
	// 2、在list前插数据
	// lpush key value
	pipeline.LRem(ctx, ListKey, 0, identity)

	pipeline.Del(ctx, HashKey)

	_, err := pipeline.Exec(ctx)

	return err
}
