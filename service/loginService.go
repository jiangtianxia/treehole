package service

import (
	"fmt"
	"strconv"
	"time"
	"treehole/dao"
	"treehole/define"
	"treehole/logger"
	"treehole/models"
	"treehole/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// SendEmailCode
// @Summary 发送邮件验证码
// @Tags 登录业务接口
// @Accept application/json
// @Produce application/json
// @Param object body utils.SendCodeForm true "发送参数"
// @Success 200 {object} utils.H
// @Router /sendEmailCode [post]
func SendEmailCode(c *gin.Context) {
	// 1、获取邮箱
	var user utils.SendCodeForm
	if err := c.ShouldBindJSON(&user); err != nil {
		// 获取valadtor.valiadtionErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非valiadtor.ValidationErrors类型错误直接返回
			utils.RespFail(c, int(define.ParamsInvalidCode), err.Error())
			logger.SugarLogger.Error("Params Invalid" + err.Error())
			return
		}
		// validator.ValidationErrors类型错误进行翻译
		hashMap := utils.RemoveTopStruct(errs.Translate(utils.Trans))
		msg := ""
		for _, v := range hashMap {
			msg += v + ","
		}
		logger.SugarLogger.Error("Params Invalid" + msg)
		utils.RespFail(c, int(define.ParamsInvalidCode), msg)
		return
	}

	// 2、获取验证码
	rand := utils.GetRand(6)

	// 3、发送验证码
	err := utils.SendEmailCode(user.Email, rand)
	if err != nil {
		logger.SugarLogger.Error("Send Email Fail Error:", err.Error())
		utils.RespFail(c, int(define.FailCode), "发送邮箱验证码失败")
		return
	}

	// 4、将验证码存入redis，过期时间为5分钟
	err = utils.RDB.Set(c, user.Email, rand, time.Second*300).Err()
	if err != nil {
		logger.SugarLogger.Error("Set Redis Fail Error:", err.Error())
		utils.RespFail(c, int(define.FailCode), "发送邮箱验证码失败")
		return
	}

	utils.RespSuccess(c, "发送邮箱验证码成功", "")
}

// Register
// @Summary 用户注册
// @Tags 登录业务接口
// @Accept application/json
// @Produce application/json
// @Param object body utils.RegisterForm true "发送参数"
// @Success 200 {object} utils.H
// @Router /register [post]
func Register(c *gin.Context) {
	// 1、获取参数
	var user utils.RegisterForm
	if err := c.ShouldBindJSON(&user); err != nil {
		// 获取valadtor.valiadtionErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非valiadtor.ValidationErrors类型错误直接返回
			utils.RespFail(c, int(define.ParamsInvalidCode), err.Error())
			logger.SugarLogger.Error("Params Invalid" + err.Error())
			return
		}
		// validator.ValidationErrors类型错误进行翻译
		hashMap := utils.RemoveTopStruct(errs.Translate(utils.Trans))
		msg := ""
		for _, v := range hashMap {
			msg += v + ","
		}
		logger.SugarLogger.Error("Params Invalid" + msg)
		utils.RespFail(c, int(define.ParamsInvalidCode), msg)
		return
	}

	// 2、验证验证码是否正确
	sysCode, err := utils.RDB.Get(c, user.Email).Result()
	if err != nil {
		utils.RespFail(c, int(define.FailCode), "验证码已过期，请重新获取")
		return
	}
	if sysCode != user.Code {
		utils.RespFail(c, int(define.FailCode), "验证码不正确，请重新输入")
		return
	}

	// 3、判断邮箱是否已存在
	cnt, err := dao.FindUserByEmailCount(user.Email)
	if cnt > 0 {
		logger.SugarLogger.Error("Eamil Exist Error:" + user.Email)
		utils.RespFail(c, int(define.EmailExistCode), define.EmailExistCode.Msg())
		return
	}
	if err != nil {
		logger.SugarLogger.Error("sql Email Error:" + user.Email)
		utils.RespFail(c, int(define.FailCode), "注册失败")
		return
	}

	// 4、判断用户名是否已经存在
	cnt, err = dao.FindUserByNameCount(user.Username)
	if cnt > 0 {
		logger.SugarLogger.Error("Username Exist Error:" + user.Email)
		utils.RespFail(c, int(define.UserExistCode), define.UserExistCode.Msg())
		return
	}
	if err != nil {
		logger.SugarLogger.Error("sql Username Error:" + user.Username)
		utils.RespFail(c, int(define.FailCode), "注册失败")
		return
	}

	// 3、生成identity和对密码加密，插入数据空间
	// md5 对密码进行加密
	password := utils.MakePassword(user.Password)
	// 使用雪花算法生成identidy
	identity, err := utils.GetID()
	if err != nil {
		logger.SugarLogger.Error("create identity Error:" + err.Error())
		utils.RespFail(c, int(define.FailCode), "注册失败")
		return
	}

	u := models.UserBasic{
		Identity: strconv.Itoa(int(identity)),
		Username: user.Username,
		Password: password,
		Email:    user.Email,
	}

	// 4、创建用户
	err = dao.CreateUser(u)
	if err != nil {
		logger.SugarLogger.Error("create User Error:" + err.Error())
		utils.RespFail(c, int(define.FailCode), "注册失败")
		return
	}

	// 使用jwt生成token，返回
	token, err := utils.GenerateToken(u.Identity, u.Username, "")
	if err != nil {
		logger.SugarLogger.Error("Generate Token Error:" + err.Error())
		utils.RespFail(c, int(define.FailCode), "注册失败")
		return
	}

	userInfo := utils.UserInfo{
		Username: user.Username,
		Usericon: "",
		Age:      18,
		Sex:      0,
	}
	data := map[string]interface{}{
		"token":    token,
		"userInfo": userInfo,
	}
	utils.RespSuccess(c, "注册成功", data)
}

// GetCapacha
// @Summary 获取验证码
// @Tags 登录业务接口
// @Success 200 {object} utils.H
// @Router /capacha/get [get]
func GetCapacha(c *gin.Context) {
	id, base64, err := utils.MakeCaptcha()
	if err != nil {
		logger.SugarLogger.Error("Make Captcha Error:" + err.Error())
		utils.RespFail(c, int(define.FailCode), "生成图片验证码失败")
		return
	}

	data := map[string]string{
		"captchaId": id,
		"base64":    base64,
	}
	utils.RespSuccess(c, "获取图片验证码成功", data)
}

// VerifyCapacha
// @Summary 验证验证码
// @Tags 登录业务接口
// @Param capachaId query string true "capachaId"
// @Param capachaVal query string true "capachaVal"
// @Success 200 {object} utils.H
// @Router /capacha/verify [get]
func VerifyCapacha(c *gin.Context) {
	id := c.DefaultQuery("capachaId", "")
	capachaVal := c.DefaultQuery("capachaVal", "")

	// id验证码id
	// answer 需要验证的内容
	// clear 校验完是否清除
	if utils.Store.Verify(id, capachaVal, true) {
		utils.RespSuccess(c, "验证图片验证码成功", "")
	} else {
		utils.RespFail(c, int(define.FailCode), "验证图片验证码失败")
	}
}

// Login
// @Summary 用户登录
// @Tags 登录业务接口
// @Accept application/json
// @Produce application/json
// @Param object body utils.LoginForm true "发送参数"
// @Success 200 {object} utils.H
// @Router /login [post]
func Login(c *gin.Context) {
	// 1、获取参数
	var user utils.LoginForm
	if err := c.ShouldBindJSON(&user); err != nil {
		// 获取valadtor.valiadtionErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非valiadtor.ValidationErrors类型错误直接返回
			utils.RespFail(c, int(define.ParamsInvalidCode), err.Error())
			logger.SugarLogger.Error("Params Invalid" + err.Error())
			return
		}
		// validator.ValidationErrors类型错误进行翻译
		hashMap := utils.RemoveTopStruct(errs.Translate(utils.Trans))
		msg := ""
		for _, v := range hashMap {
			msg += v + ","
		}
		logger.SugarLogger.Error("Params Invalid" + msg)
		utils.RespFail(c, int(define.ParamsInvalidCode), msg)
		return
	}

	// 2、判断用户是否存在
	u, err := dao.FindUserByName(user.Username)
	if err != nil {
		logger.SugarLogger.Error("Search UserName Error:" + err.Error())
		utils.RespFail(c, int(define.InvalidPasswordCode), define.InvalidPasswordCode.Msg())
		return
	}

	if u.Identity == "" {
		logger.SugarLogger.Error("Invalid Password Or UserName")
		utils.RespFail(c, int(define.InvalidPasswordCode), define.InvalidPasswordCode.Msg())
		return
	}

	// 判断用户是否已经在一分钟登录了5次
	flag, Minutes := dao.ExceLoginBank(c, u.Identity)
	if flag {
		utils.RespFail(c, int(define.FailCode), fmt.Sprint("1分钟内连续登录失败5次，账号被锁10分钟，请", int(Minutes), "分钟后重试"))
		return
	}

	// 3、判断密码是否正确
	if !utils.ValidPassword(user.Password, u.Password) {

		// 密码不正确
		dao.AddLoginError(c, u.Identity)

		logger.SugarLogger.Error("Invalid Password")
		utils.RespFail(c, int(define.InvalidPasswordCode), "密码不正确，1分钟内连续登录失败5次，该账号将会被锁10分钟")
		return
	}

	// 4、生成token返回
	token, err := utils.GenerateToken(u.Identity, u.Username, u.Usericon)
	if err != nil {
		logger.SugarLogger.Error("Generate Token Error:" + err.Error())
		utils.RespFail(c, int(define.FailCode), "登录失败")
		return
	}

	userInfo := utils.UserInfo{
		Username: u.Username,
		Usericon: u.Usericon,
		Age:      u.Age,
		Sex:      u.Sex,
	}
	data := map[string]interface{}{
		"token":    token,
		"userInfo": userInfo,
	}
	utils.RespSuccess(c, "登录成功", data)
}

// VerifyEmailCode
// @Summary 忘记密码-验证验证码
// @Tags 登录业务接口
// @Accept application/json
// @Produce application/json
// @Param object body utils.VerifyEmailCodeForm true "发送参数"
// @Success 200 {object} utils.H
// @Router /forgetPassword/verifyEmailCode [post]
func VerifyEmailCode(c *gin.Context) {
	// 1、获取参数
	var user utils.VerifyEmailCodeForm
	if err := c.ShouldBindJSON(&user); err != nil {
		// 获取valadtor.valiadtionErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非valiadtor.ValidationErrors类型错误直接返回
			utils.RespFail(c, int(define.ParamsInvalidCode), err.Error())
			logger.SugarLogger.Error("Params Invalid" + err.Error())
			return
		}
		// validator.ValidationErrors类型错误进行翻译
		hashMap := utils.RemoveTopStruct(errs.Translate(utils.Trans))
		msg := ""
		for _, v := range hashMap {
			msg += v + ","
		}
		logger.SugarLogger.Error("Params Invalid" + msg)
		utils.RespFail(c, int(define.ParamsInvalidCode), msg)
		return
	}

	// 2、判断用户该邮箱是否存在
	u, err := dao.FindByEmail(user.Email)
	if err != nil {
		logger.SugarLogger.Error("Generate Token Error:" + err.Error())
		utils.RespFail(c, int(define.EmailNotExistCode), define.EmailNotExistCode.Msg())
		return
	}

	// 3、判断code是否存在且正确
	sysCode, err := utils.RDB.Get(c, user.Email).Result()
	if err != nil {
		utils.RespFail(c, int(define.FailCode), "验证码已过期，请重新获取")
		return
	}
	if sysCode != user.Code {
		utils.RespFail(c, int(define.FailCode), "验证码不正确，请重新输入")
		return
	}

	// 4、生成token并返回
	token, err := utils.GenerateToken(u.Identity, u.Username, u.Usericon)
	if err != nil {
		logger.SugarLogger.Error("Generate Token Error:" + err.Error())
		utils.RespFail(c, int(define.FailCode), "验证失败")
		return
	}

	data := map[string]string{
		"token": token,
	}
	utils.RespSuccess(c, "验证成功", data)
}

// ModifyPassword
// @Summary 密码修改
// @Tags 登录业务接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Authorization"
// @Param object body utils.ModifyPasswordForm true "发送参数"
// @Success 200 {object} utils.H
// @Router /forgetPassword/modifyPassword [post]
func ModifyPassword(c *gin.Context) {
	// 1、获取参数
	var user utils.ModifyPasswordForm
	if err := c.ShouldBindJSON(&user); err != nil {
		// 获取valadtor.valiadtionErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非valiadtor.ValidationErrors类型错误直接返回
			utils.RespFail(c, int(define.ParamsInvalidCode), err.Error())
			logger.SugarLogger.Error("Params Invalid" + err.Error())
			return
		}
		// validator.ValidationErrors类型错误进行翻译
		hashMap := utils.RemoveTopStruct(errs.Translate(utils.Trans))
		msg := ""
		for _, v := range hashMap {
			msg += v + ","
		}
		logger.SugarLogger.Error("Params Invalid" + msg)
		utils.RespFail(c, int(define.ParamsInvalidCode), msg)
		return
	}

	// 2、根据identidy查询用户信息
	u, _ := c.Get("user_claims")
	userClaim := u.(*utils.UserClaims)

	modify, err := dao.FindByIdentity(userClaim.Identity)
	if err != nil {
		logger.SugarLogger.Error("FindByIdentity Error:" + err.Error())
		utils.RespFail(c, int(define.FailCode), "修改失败")
		return
	}

	if modify.Password == "" {
		logger.SugarLogger.Error("FindByIdentity Error:" + err.Error())
		utils.RespFail(c, int(define.FailCode), "修改失败")
		return
	}

	// 3、使用md5对密码进行加密
	Password := utils.MakePassword(user.Password)

	// 4、修改用户信息
	newUser := models.UserBasic{
		Identity: modify.Identity,
		Username: user.Username,
		Password: Password,
		Email:    modify.Email,
	}
	err = dao.ModifyUserInfo(newUser)
	if err != nil {
		logger.SugarLogger.Error("FindByIdentity Error:" + err.Error())
		utils.RespFail(c, int(define.FailCode), "修改失败")
		return
	}

	// 4、生成token并返回
	token, err := utils.GenerateToken(newUser.Identity, newUser.Username, modify.Usericon)
	if err != nil {
		logger.SugarLogger.Error("Generate Token Error:" + err.Error())
		utils.RespFail(c, int(define.FailCode), "修改失败")
		return
	}

	userInfo := utils.UserInfo{
		Username: user.Username,
		Usericon: modify.Usericon,
		Age:      modify.Age,
		Sex:      modify.Sex,
	}
	data := map[string]interface{}{
		"token":    token,
		"userInfo": userInfo,
	}
	utils.RespSuccess(c, "修改成功", data)
}
