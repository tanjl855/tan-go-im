package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var Log *zap.Logger

func init() {
	InitLog("", "", "", true)
}

func InitLog(logfilePath string, logLevel string, encodeMode string, isConsole bool) {
	// 1.配置zapcore的编码器
	zapEncode := zapcore.EncoderConfig{
		MessageKey:     "ChatLog",
		LevelKey:       "level",
		TimeKey:        "time", //时间
		NameKey:        "Name",
		CallerKey:      "Caller",
		FunctionKey:    "Function",
		StacktraceKey:  "Stacktrace",
		SkipLineEnding: false,
		LineEnding:     zapcore.DefaultLineEnding,     //输出的分割符
		EncodeLevel:    zapcore.LowercaseLevelEncoder, //序列化字符串的大小写
		//EncodeTime:          zapcore.ISO8601TimeEncoder,     //时间的编码格式
		EncodeTime:          EncodeTime,                     //时间自定义的
		EncodeDuration:      zapcore.SecondsDurationEncoder, //时间显示的位数
		EncodeCaller:        zapcore.ShortCallerEncoder,     //输出的运行文件路径长度
		EncodeName:          zapcore.FullNameEncoder,        //可选的
		NewReflectedEncoder: nil,
		ConsoleSeparator:    "", //控制台格式时，每个字段间的分割符,不配置默认即可
	}

	// 2.日志分割器
	hook := &lumberjack.Logger{
		Filename:   logfilePath,
		MaxSize:    128,
		MaxAge:     30,
		MaxBackups: 2,
		Compress:   false,
	}

	// 3.设置日志
	logLev := zap.NewAtomicLevel()
	logLevel = strings.ToLower(logLevel)
	switch logLevel {
	case "debug":
		isConsole = true
		logLev.SetLevel(zapcore.DebugLevel)
	case "info":
		logLev.SetLevel(zapcore.InfoLevel)
	case "warn":
		logLev.SetLevel(zapcore.WarnLevel)
	case "errors":
		logLev.SetLevel(zapcore.ErrorLevel)
	default:
		isConsole = true
		logLev.SetLevel(zapcore.DebugLevel)
	}
	// 4.设置zap日志输出位置，使用数组的方式便于控制输出到多个位置
	writes := []zapcore.WriteSyncer{}
	if logfilePath != "" {
		writes = append(writes, zapcore.AddSync(hook))
	}
	if isConsole {
		writes = append(writes, zapcore.AddSync(os.Stdout))
	}
	// 设置日志的编码格式json和Console
	var enc zapcore.Encoder
	if encodeMode == "json" {
		enc = zapcore.NewJSONEncoder(zapEncode)
	} else {
		enc = zapcore.NewConsoleEncoder(zapEncode)
	}
	// 5.通过传入的配置实例化core
	core := zapcore.NewCore(enc,
		zapcore.NewMultiWriteSyncer(writes...),
		logLev)
	// 6.构造日志
	// 设置为开发模式会记录panic
	development := zap.Development()
	// 构造一个字段
	zap.Fields(zap.String("im-srv-test-log", "test_1"))
	//通过传入的配置实例化一个日志
	zapLogger := zap.New(core, development)
	zapLogger.Info("初始化日志")
	Log = zapLogger
}

// EncodeTime 自定义时间输出编码器
func EncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 - 15:04:05.000"))
}

func Panic(v ...any) {
	Log.Panic(fmt.Sprint(v...), zap.String("caller", GetCaller()))
}

func Debug(v ...any) {
	Log.Debug(fmt.Sprint(v...), zap.String("caller", GetCaller()))
}

func Info(v ...any) {
	Log.Info(fmt.Sprint(v...), zap.String("caller", GetCaller()))
}

func Error(v ...any) {
	Log.Error(fmt.Sprint(v...), zap.String("caller", GetCaller()))
}

func Warn(v ...any) {
	Log.Warn(fmt.Sprint(v...), zap.String("caller", GetCaller()))
}

func GetCaller() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	return file + ":" + strconv.Itoa(line)
}
