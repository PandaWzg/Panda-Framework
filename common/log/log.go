package log

import (
	"fmt"
	rotatelog "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
	"Panda/conf"
)

var logger *zap.Logger
/**
 * 初始化日志
 * logPath 日志文件路径
 * loglevel 日志级别
 * maxSize 每个日志文件保存的最大尺寸 单位：M
 * maxBackups 日志文件最多保存多少个备份
 * maxAge 文件最多保存多少天
 * compress 是否压缩
 * serviceName 服务名 经常
 */
func Init(conf *conf.Cfg) *zap.Logger {
	//hook := &lumberjack.Logger{
	//	Filename:   conf.Config.Site.Logfile + logPath + ".log",
	//	//Filename:   "log/" + logPath, // ⽇志⽂件路径
	//	MaxSize:    1,    // megabytes
	//	MaxBackups: 2,       // 最多保留3个备份
	//	MaxAge:     15,       //days
	//	Compress:   false,    // 是否压缩 disabled by default
	//}
	logPath := conf.LogFile + ".log"
	hook, _ := rotatelog.New(
		logPath + ".%Y%m%d",
		rotatelog.WithLinkName(logPath),
		rotatelog.WithMaxAge(time.Hour*24*30),
		rotatelog.WithRotationTime(time.Hour),
	)
	w := zapcore.AddSync(hook)
	var level zapcore.Level
	switch conf.LogLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = timeEncoder
	//encoderConfig.EncodeCaller = shortCallerWithClassFunctionEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		w,
		level,
	)
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel))
	logger.Info("DefaultLogger init success")
	return logger
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}


func Debug(args ...interface{}) {
	logger.Debug(fmt.Sprint(args...))
}

func Infof(format string, v ...interface{}) {
	logger.Info(fmt.Sprintf(format, v...))
}

func Warnf(format string, v ...interface{}) {
	logger.Warn(fmt.Sprintf(format, v...))
}

func Errorf(format string, v ...interface{}) {
	logger.Error(fmt.Sprintf(format, v...))
}

func Debugf(format string, v ...interface{}) {
	logger.Debug(fmt.Sprintf(format, v...))
}

func Fatalf(format string, v ...interface{}) {
	logger.Fatal(fmt.Sprintf(format, v...))
}

func Info(args ...interface{}) {
	logger.Info(fmt.Sprint(args...))
}

func Warn(args ...interface{}) {
	logger.Warn(fmt.Sprint(args...))
}

func Error(args ...interface{}) {
	logger.Error(fmt.Sprint(args...))
}

func Fatal(args ...interface{}) {
	logger.Fatal(fmt.Sprint(args...))
}

func Panic(args ...interface{}) {
	logger.Panic(fmt.Sprint(args...))
}