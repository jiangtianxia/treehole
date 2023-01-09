package router

import (
	"treehole/middlewares"
	"treehole/service"

	"github.com/gin-gonic/gin"

	docs "treehole/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()

	// swagger 配置
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 路由配置
	v1 := r.Group("/api/v1")
	{
		/*
		* 公共接口
		 */
		v1.GET("/hello", service.Hello)

		// 图片上传
		v1.POST("/uploadLocal", middlewares.AuthUserCheck(), service.UploadLocal)

		/*
		* 登录业务接口
		 */
		// 发送邮件验证码
		v1.POST("/sendEmailCode", service.SendEmailCode)

		// 用户注册
		v1.POST("/register", service.Register)

		// 图片验证码业务
		captcha := v1.Group("/capacha")
		{
			// 获取图片验证码
			captcha.GET("/get", service.GetCapacha)

			// 验证码的校验
			captcha.GET("/verify", service.VerifyCapacha)
		}

		// 用户登录
		v1.POST("/login", service.Login)

		// 找回密码业务
		forget := v1.Group("/forgetPassword")
		{
			// 1、获取邮箱验证码
			// 使用前面接口

			// 2、验证邮箱验证码和邮箱信息
			forget.POST("/verifyEmailCode", service.VerifyEmailCode)

			// 3、修改用户密码及信息
			forget.POST("/modifyPassword", middlewares.AuthUserCheck(), service.ModifyPassword)
		}

		/*
		* 用户业务接口
		 */
		user := v1.Group("/user", middlewares.AuthUserCheck())
		{
			// 1、获取当前用户信息
			user.GET("/getUserInfo", service.GetUserInfo)

			// 2、修改用户信息
			// 上传图片：使用公共接口的上传文件，然后得到url
			// 修改年龄，性别，用户名等信息
			user.POST("/modifyUserInfo", service.ModifyUserInfo)

			// 3、更换密码
			// 旧密码、新密码、新确认密码
			user.POST("/userModifyPassword", service.UserModifyPassword)
		}

	}

	return r
}
