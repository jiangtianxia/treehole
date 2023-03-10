package main

import (
	"fmt"
	"time"
	"treehole/logger"
	"treehole/router"
	"treehole/service"
	"treehole/utils"

	"github.com/spf13/viper"
)

func main() {
	/*
	* 初始化
	 */
	// 初始化配置
	utils.InitConfig()

	// 初始化翻译器
	if err := utils.InitTrans("zh"); err != nil {
		fmt.Println("init trans failed, err:", err)
		return
	}
	fmt.Println("trans inited ......")

	// 初始化雪花算法
	if err := utils.SnowflakeInit(viper.GetUint16("snowflake.machineID")); err != nil {
		fmt.Println("init Snowflake failed, err:", err)
		return
	}
	fmt.Println("Snowflake inited ...... ")

	// 初始化日志
	logger.InitLogger()
	defer logger.SugarLogger.Sync() // 刷新流，写日志到文件中

	// 初始化Mysql和Redis
	utils.InitMysql()
	utils.InitRedis()
	defer utils.ReidsClose()

	// 初始化令牌桶
	utils.InitCurrentLimit()
	logger.SugarLogger.Info("初始化配置完成")

	// 初始化定时任务
	InitTimer()

	// 配置路由
	r := router.Router()
	logger.SugarLogger.Info("配置路由完成")
	r.Run()
}

// 初始化定时任务
func InitTimer() {
	utils.Timer(time.Second*viper.GetDuration("timeout.DelayHeartbeat"), time.Second*time.Duration(viper.GetInt("timeout.HeartbeatHz")), service.CleanConnection, "")
	fmt.Println("Timer inited ...... ")
	logger.SugarLogger.Info("Timer inited ......")
}
