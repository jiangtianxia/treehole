package utils

/**
 * @Author jiang
 * @Description 请求参数结构体
 * @Date 11:00 2023/1/8
 **/

/**
 * @Author jiang
 * @Description 发送邮箱验证码请求参数
 * @Date 11:00 2023/1/8
 **/
type SendCodeForm struct {
	Email string `json:"email" binding:"required,email=email"`
}

/**
 * @Author jiang
 * @Description 注册请求参数
 * @Date 11:00 2023/1/8
 **/
type RegisterForm struct {
	Username   string `json:"username" binding:"required"`
	Email      string `json:"email" binding:"required,email=email"`
	Code       string `json:"code" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"repassword" binding:"required,eqfield=Password"`
}

/**
 * @Author jiang
 * @Description 登录请求参数
 * @Date 22:00 2023/1/8
 **/
type LoginForm struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"repassword" binding:"required,eqfield=Password"`
}

/**
 * @Author jiang
 * @Description 验证邮件验证码请求参数
 * @Date 22:00 2023/1/8
 **/
type VerifyEmailCodeForm struct {
	Email string `json:"email" binding:"required,email=email"`
	Code  string `json:"code" binding:"required"`
}

/**
 * @Author jiang
 * @Description 修改密码请求参数
 * @Date 22:00 2023/1/8
 **/
type ModifyPasswordForm struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"repassword" binding:"required,eqfield=Password"`
}
