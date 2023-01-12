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

/**
 * @Author jiang
 * @Description 修改用户信息请求参数
 * @Date 15:00 2023/1/9
 **/
type ModifyUserInfoForm struct {
	Username string `json:"username" binding:"required"`
	Age      string `json:"age" binding:"required"`
	Sex      string `json:"sex" binding:"required"`
	Url      string `json:"url" binding:"required"`
}

/**
 * @Author jiang
 * @Description 更换密码请求参数
 * @Date 21:00 2023/1/9
 **/
type UserModifyPasswordForm struct {
	NowPassword string `json:"nowpassword" binding:"required"`
	Password    string `json:"password" binding:"required"`
	RePassword  string `json:"repassword" binding:"required,eqfield=Password"`
}

/**
 * @Author jiang
 * @Description 创建帖子请求参数
 * @Date 16:00 2023/1/9
 **/
type CreateNoteForm struct {
	Title   string `json:"title" binding:"required"`
	Urls    string `json:"urls"`
	Content string `json:"content" binding:"required"`
}

/**
 * @Author jiang
 * @Description 获取帖子详细信息
 * @Date 13:00 2023/1/11
 **/
type GetNoteInfoFrom struct {
	NoteIdentity string `json:"note_identity" binding:"required"`
}

/**
 * @Author jiang
 * @Description 修改帖子请求参数
 * @Date 16:00 2023/1/9
 **/
type ModifyNoteForm struct {
	NoteIdentity string `json:"note_identity" binding:"required"`
	Title        string `json:"title" binding:"required"`
	Urls         string `json:"urls"`
	Content      string `json:"content" binding:"required"`
}

/**
 * @Author jiang
 * @Description 点赞或踩请求参数
 * @Date 17:00 2023/1/12
 **/
type VotedNoteFrom struct {
	NoteIdentity string `json:"note_identity" binding:"required"`
	Voted        string `json:"voted" binding:"required"`
	Flag         string `json:"flag" binding:"required"`
}

/**
 * @Author jiang
 * @Description 发送评论请求参数
 * @Date 21:00 2023/1/12
 **/
type CreateCommentFrom struct {
	NoteIdentity string `json:"note_identity" binding:"required"`
	Content      string `json:"content" binding:"required"`
}

/**
 * @Author jiang
 * @Description 删除评论请求参数
 * @Date 23:30 2023/1/12
 **/
type DeleteNoteCommentFrom struct {
	CommentIdentity string `json:"comment_identity" binding:"required"`
	NoteIdentity    string `json:"note_identity" binding:"required"`
}
