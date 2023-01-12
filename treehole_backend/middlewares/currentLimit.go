package middlewares

import (
	"treehole/define"
	"treehole/utils"

	"github.com/gin-gonic/gin"
)

/**
 * @Author jiang
 * @Description 令牌桶限流策略
 * @Date 9:00 2023/1/10
 **/
func CurrentLimit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 判断令牌桶中是否有令牌
		if !utils.Bucket.Allow() {
			// 不允许访问
			ctx.Abort()
			utils.RespFail(ctx, int(define.ServerBusyCode), define.ServerBusyCode.Msg())
			return
		}
		ctx.Next()
	}
}
