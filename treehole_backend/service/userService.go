package service

import (
	"os"
	"strconv"
	"treehole/dao"
	"treehole/define"
	"treehole/logger"
	"treehole/models"
	"treehole/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// GetUserInfo
// @Summary 获取用户信息
// @Tags 用户业务接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Authorization"
// @Success 200 {object} utils.H
// @Router /user/getUserInfo [get]
func GetUserInfo(c *gin.Context) {
	// 1、获取token上的用户信息
	u, _ := c.Get("user_claims")
	userClaim := u.(*utils.UserClaims)

	// 2、根据identity，查询出用户信息，返回
	user, err := dao.FindByIdentity(userClaim.Identity)
	if err != nil {
		logger.SugarLogger.Error("FindByIdentity Error:" + err.Error())
		utils.RespFail(c, int(define.FailCode), "获取用户信息失败")
		return
	}

	data := map[string]interface{}{
		"username": user.Username,
		"usericon": user.Usericon,
		"sex":      user.Sex,
		"age":      user.Age,
	}
	utils.RespSuccess(c, "获取用户信息成功", data, 0)
}

// ModifyUserInfo
// @Summary 修改用户信息
// @Tags 用户业务接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Authorization"
// @Param object body utils.ModifyUserInfoForm true "发送参数"
// @Success 200 {object} utils.H
// @Router /user/modifyUserInfo [post]
func ModifyUserInfo(c *gin.Context) {
	// 1、获取参数
	var user utils.ModifyUserInfoForm
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

	// 2、获取token信息
	u, _ := c.Get("user_claims")
	userClaim := u.(*utils.UserClaims)

	// 3、判断用户名是否已经存在
	temp, err := dao.FindUserByNameAndIdentity(user.Username, userClaim.Identity)
	if err != nil {
		logger.SugarLogger.Error("FindUserByName Error:"+err.Error()+"username:", user.Username)
		utils.RespFail(c, int(define.FailCode), "修改失败")
		return
	}

	if temp > 0 {
		logger.SugarLogger.Error("User Exist Error username:", user.Username)
		utils.RespFail(c, int(define.UserExistCode), define.UserExistCode.Msg())
		return
	}

	tx := utils.DB.Begin()
	//事务一旦开始，不论什么异常最终都会 Rollback
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 4、修改数据
	sex, _ := strconv.Atoi(user.Sex)
	age, _ := strconv.Atoi(user.Age)
	userBasic := models.UserBasic{
		Identity: userClaim.Identity,
		Username: user.Username,
		Usericon: user.Url,
		Sex:      sex,
		Age:      age,
	}
	err = dao.ModifyUserInfo(userBasic)
	if err != nil {
		logger.SugarLogger.Error("Modify UserInfo Error:"+err.Error()+"username:", user.Username)
		tx.Rollback()
		utils.RespFail(c, int(define.FailCode), "修改失败")
		return
	}

	// 5、查询图片库在头像类但不为头像的图片
	ImageList, err := dao.FindImageByIdentity(userClaim.Identity, userBasic.Usericon)
	if err != nil {
		logger.SugarLogger.Error("FindImageByIdentity Error:"+err.Error()+"username:", user.Username)
		tx.Rollback()
		utils.RespFail(c, int(define.FailCode), "修改失败")
		return
	}

	// 6、删除数据库图片
	err = dao.DeleteImages(ImageList)
	if err != nil {
		logger.SugarLogger.Error("DeleteImages Error:"+err.Error()+"username:", user.Username)
		tx.Rollback()
		utils.RespFail(c, int(define.FailCode), "修改失败")
		return
	}

	tx.Commit()
	// 7、删除本地图片
	for _, image := range ImageList {
		if image.Url != userBasic.Usericon {
			os.Remove(image.Url)
		}
	}
	utils.RespSuccess(c, "修改用户信息成功", "", 0)
}

// UserModifyPassword
// @Summary 更换密码
// @Tags 用户业务接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Authorization"
// @Param object body utils.UserModifyPasswordForm true "发送参数"
// @Success 200 {object} utils.H
// @Router /user/userModifyPassword [post]
func UserModifyPassword(c *gin.Context) {
	// 1、获取参数
	var user utils.UserModifyPasswordForm
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

	// 2、获取token信息
	t, _ := c.Get("user_claims")
	userClaim := t.(*utils.UserClaims)

	// 3、判断用户是否存在
	u, err := dao.FindByIdentity(userClaim.Identity)
	if err != nil {
		logger.SugarLogger.Error("FindByIdentity Error:" + err.Error())
		utils.RespFail(c, int(define.InvalidTokenCode), define.InvalidPasswordCode.Msg())
		return
	}

	if u.Identity == "" {
		logger.SugarLogger.Error("Invalid Password ")
		utils.RespFail(c, int(define.FailCode), "密码错误")
		return
	}

	// 3、判断密码是否正确
	if !utils.ValidPassword(user.NowPassword, u.Password) {
		logger.SugarLogger.Error("Invalid Password")
		utils.RespFail(c, int(define.InvalidPasswordCode), define.InvalidPasswordCode.Msg())
		return
	}

	// 4、修改密码并生成token返回
	newPassword := utils.MakePassword(user.Password)

	// 开启事务来修改
	tx := utils.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 修改密码
	NewUser := models.UserBasic{
		Identity: userClaim.Identity,
		Password: newPassword,
	}

	err = dao.ModifyUserInfo(NewUser)
	if err != nil {
		logger.SugarLogger.Error("Modify UserInfo Error:" + err.Error())
		tx.Rollback()
		utils.RespFail(c, int(define.FailCode), "修改失败")
		return
	}
	token, err := utils.GenerateToken(u.Identity, u.Username)
	if err != nil {
		logger.SugarLogger.Error("Generate Token Error:" + err.Error())
		tx.Rollback()
		utils.RespFail(c, int(define.FailCode), "修改失败")
		return
	}
	tx.Commit()

	data := map[string]string{
		"token": token,
	}
	utils.RespSuccess(c, "修改成功", data, 0)
}
