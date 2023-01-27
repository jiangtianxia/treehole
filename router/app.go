package router

import (
	"net/http"
	"treehole/middlewares"
	"treehole/service"

	"github.com/gin-gonic/gin"

	docs "treehole/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()

	// 设置成发布模式
	gin.SetMode(gin.ReleaseMode)

	// swagger 配置
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//加载静态资源，一般是上传的资源，例如用户上传的图片
	r.StaticFS("/upload", http.Dir("upload"))

	// 路由配置
	v1 := r.Group("/api/v1", middlewares.CurrentLimit())
	{
		/*
		* 公共接口
		 */
		v1.GET("/hello", service.Hello)

		// 图片上传
		v1.POST("/uploadLocal", middlewares.AuthUserCheck(), service.UploadLocal)

		// 搜索帖子信息
		v1.GET("/searchNotes", service.SearchNotes)

		// 按照热度或时间读取帖子信息
		v1.GET("/searchNotesScoreOrTime", service.SearchNotesScoreOrTime)

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

		/*
		* 帖子业务接口
		 */
		note := v1.Group("/note", middlewares.AuthUserCheck())
		{
			// 创建帖子
			note.POST("/createNote", service.CreateNote)

			// 获取发布帖子列表
			note.GET("/getNoteList", service.GetNoteList)

			// 删除帖子
			note.POST("/deleteNote", service.DeleteNote)

			// 获取帖子详细信息
			note.POST("/getNoteInfo", service.GetNoteInfo)

			// 修改帖子
			note.POST("/modifyNote", service.ModifyNote)

			// 点赞帖子
			note.POST("/votedNote", service.VotedNote)
		}

		/*
		* 评论业务接口
		 */
		comment := v1.Group("/comment", middlewares.AuthUserCheck())
		{
			// 发送评论
			comment.POST("/createComment", service.CreateComment)

			// 获取文章评论
			comment.GET("/getNoteCommentList", service.GetNoteCommentList)

			// 获取评论记录
			comment.GET("/getCommentList", service.GetCommentList)

			// 删除评论记录
			comment.POST("/deleteNoteComment", service.DeleteNoteComment)
		}

		/*
		* 聊天室业务接口
		 */
		chat := v1.Group("chat", middlewares.AuthUserCheck())
		{
			// 进入聊天室建立聊天
			chat.GET("/websocket/message", service.WebsocketMessage)

			// 获取当前在线人数
			chat.GET("/websocket/getOnlineList", service.GetOnlistList)

			// 获取聊天记录
			chat.GET("/websocket/getMessageList", service.GetMessageList)
		}
	}

	return r
}
