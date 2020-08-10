package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	logger *zap.Logger
)

// Init 初始化zap日志
func Init() {
	writeSyncer := getLogWriter() // 日志输出位置相关
	errWriteSyncer := getErrLogWriter()
	encoder := getEncoder() // 日志的格式相关

	level := zap.AtomicLevel{} // 解析字符串格式的级别
	if err := level.UnmarshalText([]byte(viper.GetString("log.level"))); err != nil {
		level = zap.NewAtomicLevel() // 默认用info
	}
	//core := zapcore.NewCore(encoder, writeSyncer, level)
	// 根据app的模式把日志输出到不同的位置
	var core zapcore.Core
	if viper.GetString("app.mode") == gin.DebugMode {
		// consoleEncoder 一个往终端输出日志的配置对象
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		// NewTee 可以指定多个日志配置
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, level),
			// 创建一个将debug级别以上的日志输出到终端的配置信息
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
			// 将error级别以上的日志输出到err文件
			zapcore.NewCore(encoder, errWriteSyncer, zapcore.ErrorLevel),
		)
	} else {
		core = zapcore.NewCore(encoder, writeSyncer, level)
	}

	logger := zap.New(core, zap.AddCaller()) // 根据上面的配置创建logger
	zap.ReplaceGlobals(logger)               // 替换zap库里全局的logger,后序在代码中使用zap.L()来记录日志
	//sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	//encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 人类可读的时间格式
	//encoderConfig.EncodeTime = zapcore.EpochTimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	//return zapcore.NewConsoleEncoder(encoderConfig)  // 可读日志
	return zapcore.NewJSONEncoder(encoderConfig) // json格式日志
}

//func getLogWriter() zapcore.WriteSyncer {
//	file, _ := os.Create("./test.log")
//	return zapcore.AddSync(file)
//}

// lumberjack 添加日志切割归档功能
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   viper.GetString("log.filename"), // 日志文件的位置
		MaxSize:    viper.GetInt("log.max_size"),    // 日志文件在进行切割前最大大小 单位：MB
		MaxBackups: viper.GetInt("log.max_backups"), // 保留旧文件的最大个数
		MaxAge:     viper.GetInt("log.max_age"),     // 保留旧文件的最大天数 单位：天
		Compress:   viper.GetBool("log.compress"),   // 是否压缩/归档旧文件
	}
	return zapcore.AddSync(lumberJackLogger)
}

func getErrLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   viper.GetString("log.filename") + ".err",
		MaxSize:    viper.GetInt("log.max_size"),    // 日志文件大小 单位：MB
		MaxBackups: viper.GetInt("log.max_backups"), // 备份数量
		MaxAge:     viper.GetInt("log.max_age"),     // 备份时间 单位：天
		Compress:   viper.GetBool("log.compress"),   // 是否压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}
