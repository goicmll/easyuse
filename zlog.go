package habits

import (
	"io"
	"os"
	"path"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logDir string

// 设置工作目录
func SetLogDir(path string) {
	logDir = path
	newLogger()
}

var level int8 = 0

// 设置日志级别
func SetLevel(lv int8) {
	level = lv
	newLogger()
}

var enableFile bool = false

// 设置是否写入文件
func SetEnableFile(enable bool) {
	enableFile = enable
	newLogger()
}

var Logger *zap.Logger
var Sugar *zap.SugaredLogger
var encoderConfig zapcore.EncoderConfig

// 初始化日志对象
func initZapLog() {
	// 设置默认的日志目录
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	SetLogDir(path.Join(wd, "logs"))
	newLogger()
}

// 写入文件
func getFileWriter(filePath string) io.Writer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    100, // 单个文件最大100M
		MaxBackups: 20,  // 多于 60 个日志文件后，清理较旧的日志
		MaxAge:     1,   // 一天一切割
		Compress:   true,
		LocalTime:  true,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// 新建Logger和Sugar对象
func newLogger() {

	encoderConfig = zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "dt"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339))
	}
	var zapCores = make([]zapcore.Core, 1, 4)
	//同时将日志输出到控制台
	zapCores[0] = zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), zapcore.Level(level))
	if enableFile {
		// 创建日志目录
		if !IsDir(logDir) {
			err := os.Mkdir(logDir, 0744)
			if err != nil {
				panic(err)
			}
		}
		// 获取文件的 io.Writer
		var infoWriter io.Writer
		var errorWriter io.Writer
		infoWriter = getFileWriter(path.Join(logDir, "info.log"))
		errorWriter = getFileWriter(path.Join(logDir, "error.log"))
		// 自定义写入 info.log 文件的日志级别
		var infoLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl <= zapcore.WarnLevel && lvl >= zapcore.Level(level)
		})

		// 自定义写入 error.log 文件的日志级别
		var errorLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})

		zapCores = append(
			zapCores,
			//将info及以下写入logPath
			zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(infoWriter), infoLevel),
			//error及以上写入errPath
			zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(errorWriter), errorLevel),
		)
	}
	core := zapcore.NewTee(zapCores...)
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	Sugar = Logger.Sugar()
}

func GinFields(requestID, uri, serviceName string) []zap.Field {
	return []zap.Field{zap.String("RID", requestID), zap.String("URI", uri), zap.String("SN", serviceName)}
}
