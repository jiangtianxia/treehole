package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var SugarLogger *zap.SugaredLogger

func InitLogger() {
	encoder := getEncode()
	writerSyner := getLogWriter()
	core := zapcore.NewCore(encoder, writerSyner, zapcore.DebugLevel)

	// zap.AddCaller() 添加将调用函数信息记录到日志中的功能
	logger := zap.New(core, zap.AddCaller())
	SugarLogger = logger.Sugar()

	fmt.Println("logger inited ......")
	SugarLogger.Info("logger inited")
}

// 编码器
func getEncode() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()

	// 修改时间编码器
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 在日志文件中使用大写字母记录日志级别
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	// NewConsoleEncoder 打印更符合人们观察的方式
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// 打开文件，并返回文件句柄
func getLogWriter() zapcore.WriteSyncer {
	// 获取当前工作目录
	work, _ := os.Getwd()

	dirName := work + viper.GetString("logger.path") + "/" + time.Now().Format("2006-01-02")
	_, err := os.Stat(dirName)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(dirName, 0777)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	file, err := os.Create(dirName + "/log.log")
	if err != nil {
		fmt.Println(err)
	}
	return zapcore.AddSync(file)
}
