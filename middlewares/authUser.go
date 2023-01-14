package middlewares

import (
	"time"
	"treehole/dao"
	"treehole/define"
	"treehole/logger"
	"treehole/utils"

	"github.com/gin-gonic/gin"
)

/**
 * @Author jiang
 * @Description 验证用户token信息
 * @Date 23:00 2023/1/8
 **/

// 验证用户token信息
func AuthUserCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.GetHeader("Authorization")
		userClaim, err := utils.AnalyseToken(auth)
		if err != nil {
			ctx.Abort()
			logger.SugarLogger.Error("Unauthorized Authorization")
			utils.RespFail(ctx, int(define.InvalidTokenCode), define.InvalidTokenCode.Msg())
			return
		}

		if userClaim == nil || userClaim.Identity == "" || userClaim.Issuer != "treehole" || userClaim.Username == "" {
			ctx.Abort()
			logger.SugarLogger.Error("Unauthorized User")
			utils.RespFail(ctx, int(define.InvalidTokenCode), define.InvalidTokenCode.Msg())
			return
		}

		// 判断token是否过期
		// fmt.Println("TOKEN:", userClaim.ExpiresAt)
		// fmt.Println("now ", time.Now().Unix())
		if time.Now().Unix() > userClaim.ExpiresAt {
			ctx.Abort()
			logger.SugarLogger.Error("Token Expired")
			utils.RespFail(ctx, int(define.ExpiredTokenCode), define.ExpiredTokenCode.Msg())
			return
		}

		// 判断是否存在该用户
		cnt, err := dao.FindUserByIdentityCount(userClaim.Identity)
		if err != nil {
			ctx.Abort()
			logger.SugarLogger.Error("User Invalid", userClaim)
			utils.RespFail(ctx, int(define.InvalidTokenCode), define.InvalidTokenCode.Msg())
			return
		}
		if cnt <= 0 {
			ctx.Abort()
			logger.SugarLogger.Error("User Invalid", userClaim)
			utils.RespFail(ctx, int(define.InvalidTokenCode), define.InvalidTokenCode.Msg())
			return
		}
		ctx.Set("user_claims", userClaim)
		ctx.Next()
	}
}
