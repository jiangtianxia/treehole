package main

import (
	"fmt"
	"image/color"

	"github.com/mojocn/base64Captcha"
)

/**
 * @Author jiang
 * @Description 生成图片验证码
 * @Date 20:00 2023/1/8
 **/

// 设置自带的 store（可以配置成redis）
var store = base64Captcha.DefaultMemStore

//获取验证码
func MakeCaptcha() (id, b64s string, err error) {
	var driver base64Captcha.Driver
	//配置验证码的参数
	driverString := base64Captcha.DriverString{
		Height:          40,
		Width:           100,
		NoiseCount:      0,
		ShowLineOptions: 2 | 4,
		Length:          4,
		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		BgColor:         &color.RGBA{R: 3, G: 102, B: 214, A: 125},
		Fonts:           []string{"wqy-microhei.ttc"},
	}
	//ConvertFonts 按名称加载字体
	driver = driverString.ConvertFonts()
	//创建 Captcha
	captcha := base64Captcha.NewCaptcha(driver, store)
	//Generate 生成随机 id、base64 图像字符串
	id, b64s, err = captcha.Generate()
	return id, b64s, err

}

//验证验证码
func VerifyCaptcha(id string, VerifyValue string) bool {
	fmt.Println(id, VerifyValue)
	if store.Verify(id, VerifyValue, true) {
		//验证成功
		return true
	} else {
		//验证失败
		return false
	}
}

func main() {
	id, b64s, _ := MakeCaptcha()

	fmt.Println(b64s)
	fmt.Println(VerifyCaptcha(id, b64s))
}
