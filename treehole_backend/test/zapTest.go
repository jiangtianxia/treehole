package main

import "treehole/logger"

func main() {
	logger.InitLogger()
	defer logger.SugarLogger.Sync() // 刷新流，写日志到文件中
	// 循环写入日志，这里为了达到切割要求循环记录了 100000 次
	for i := 0; i < 100000; i++ {
		logger.SugarLogger.Infof("this is a test")
	}

}
