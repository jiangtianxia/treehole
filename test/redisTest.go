package main

// import (
// 	"treehole/utils"

// 	"github.com/go-redis/redis"
// )

// type Z struct {
// 	Score  int
// 	Member string
// }

// func main() {
// 	pipeline := utils.RDB.TxPipeline()
// 	pipeline.ZAdd("jianzg", redis.Z{ // 作者默认投赞成票
// 		Score:  1,
// 		Member: 124,
// 	})

// 	_, err = pipeline.Exec()
// }
