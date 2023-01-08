package main

import (
	"fmt"
	"treehole/logger"
	"treehole/router"
	"treehole/utils"

	"github.com/spf13/viper"
)

func main() {
	// 初始化配置
	utils.InitConfig()
	if err := utils.InitTrans("zh"); err != nil {
		fmt.Println("init trans failed, err:", err)
		return
	}
	fmt.Println("trans inited ......")
	if err := utils.SnowflakeInit(viper.GetUint16("snowflake.machineID")); err != nil {
		fmt.Println("init Snowflake failed, err:", err)
		return
	}
	fmt.Println("Snowflake inited ...... ")
	logger.InitLogger()
	utils.InitMysql()
	utils.InitRedis()
	defer utils.ReidsClose()
	logger.SugarLogger.Info("初始化配置完成")

	// 配置路由
	r := router.Router()
	logger.SugarLogger.Info("配置路由完成")
	r.Run()
}
