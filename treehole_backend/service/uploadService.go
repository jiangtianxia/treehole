package service

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"treehole/dao"
	"treehole/define"
	"treehole/logger"
	"treehole/models"
	"treehole/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// UploadLocal
// @Summary 上传图片
// @Tags 公共接口
// @Param Authorization header string true "Authorization"
// @Param file formData file true "文件"
// @Param type formData string true "文件类型"
// @Success 200 {object} utils.H
// @Router /uploadLocal [post]
func UploadLocal(c *gin.Context) {
	// 1、获取token上的用户信息
	u, _ := c.Get("user_claims")
	userClaim := u.(*utils.UserClaims)

	// 2、获取发送过来的文件
	srcFile, err := c.FormFile("file")
	if err != nil {
		logger.SugarLogger.Error("Get FormFile Error:" + err.Error())
		utils.RespFail(c, int(define.ParamsInvalidCode), define.ParamsInvalidCode.Msg())
		return
	}
	imageType, err := strconv.Atoi(c.DefaultPostForm("type", "0"))
	if err != nil {
		logger.SugarLogger.Error("Get FormFile Error:" + err.Error())
		utils.RespFail(c, int(define.ParamsInvalidCode), define.ParamsInvalidCode.Msg())
		return
	}

	// 3、读取后缀名以及相关信息
	suffix := path.Ext(srcFile.Filename)
	// 判断后缀是否为图片信息
	if suffix != ".png" && suffix != ".jpeg" && suffix != ".jpg" {
		utils.RespFail(c, int(define.FailCode), "请上传文件图片")
		return
	}
	iconid, err := utils.GetID()
	if err != nil {
		logger.SugarLogger.Error("snowflake Error:" + err.Error())
		utils.RespFail(c, int(define.FailCode), "上传文件失败")
		return
	}

	fileName := fmt.Sprintf("%d%s", int(iconid), suffix)

	// 4、创建文件目录及文件
	dirName := viper.GetString("uplaodBase") + userClaim.Identity
	_, err = os.Stat(dirName)
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(dirName, 0777)
		}
	}

	dstFile := dirName + "/" + fileName
	// fmt.Println(dstFile)

	err = c.SaveUploadedFile(srcFile, dstFile)
	if err != nil {
		logger.SugarLogger.Error("Save Uploaded File Error:" + err.Error())
		utils.RespFail(c, int(define.FailCode), "上传文件失败")
		return
	}

	// 5、将url返回，同时将图片存入图片数据库
	imageBasic := models.ImageBasic{
		Identity: userClaim.Identity,
		Url:      dstFile[1:],
		Type:     imageType,
	}
	err = dao.CreateImage(imageBasic)
	if err != nil {
		logger.SugarLogger.Error("Insert into ImageBasic Error:" + err.Error() + "url：" + dstFile)
		utils.RespFail(c, int(define.FailCode), "上传文件失败")

		// 将文件删除
		os.Remove(dstFile)
		return
	}

	data := map[string]string{
		"url": imageBasic.Url,
	}
	utils.RespSuccess(c, "上传文件成功", data, 0)
}
