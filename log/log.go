package log

import (
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var sugaredLogger *zap.SugaredLogger

/*
1. logFile:         指定日志的输出文件,
2. logLevel:        指定日志的输出级别,一共有debug,info,warn,error,dPanic,panic,fatal七个日志级别,默认为info
3. logOutFormat:    指定日志的输出格式,可以是json或console,默认输出方式为console
4. maxSize:         指定日志文件的存储大小,以M为单位,值为int类型
5. maxBackups:      指定日志文件的最大备份数量
6. maxAge:          指定日志备份文件最大保存的天数
7. localTime:       指定日志备份文件是否以当前计算机时间格式命令，该值为bool值, true或false,
8. compress:        指定日志备份文件是否使用gzip方式压缩，该值为bool值, true或false
*/

func InitLoggerFromParams(logFile string, logLevel string, logOutFormat string, maxSize int, maxBackups int, maxAge int, localTime bool, compress bool) {

	logCore := setZapLoggerCore(logFile, logLevel, logOutFormat, maxSize, maxAge, maxBackups, localTime, compress)

	logger = zap.New(logCore, zap.AddCaller(), zap.AddCallerSkip(1))

	sugaredLogger = logger.Sugar()

}

func InitFromConfig(configType string, filePath string) {

}

// 定义日志打印的格式
func setZapLoggerEncoder(logOutFormat string) zapcore.Encoder {

	var encoder zapcore.Encoder

	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.TimeKey = "time"
	
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	encoderConfig.EncodeDuration = func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendInt64(int64(d) / 1000000)
	}

	/*
		encoderConfig := zapcore.EncoderConfig{
			MessageKey:   "msg",
			LevelKey:     "level",
			TimeKey:      "time",
			NameKey:      "mqtt",
			CallerKey:    "file",
			EncodeLevel:  zapcore.CapitalColorLevelEncoder, //将日志级别转换成大写（INFO，WARN，ERROR等）
			EncodeCaller: zapcore.ShortCallerEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("2006-01-02 15:04:05"))
			}, //输出的时间格式
			EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendInt64(int64(d) / 1000000)
			},
		}
	*/

	switch logOutFormat {
	case "json":
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	case "console":
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	return encoder
}

// 定义日志的writer，这里我们使用lumberjack
func setZapLoggerLogFileWriter(logFile string, maxSize int, maxAge int, maxbackups int, localTime bool, compress bool) zapcore.WriteSyncer {

	logFileWriterConfig := &lumberjack.Logger{

		Filename:   logFile,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxbackups,
		LocalTime:  localTime,
		Compress:   compress,
	}

	logFileWriterSyncer := zapcore.AddSync(logFileWriterConfig)

	return logFileWriterSyncer

}

// 定义日志输出的级别
func getZapLoggerLevel(logLever string) zapcore.Level {

	var level zapcore.Level
	switch logLever {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dPanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	return level
}

// 定义日志的core文件
func setZapLoggerCore(logFile string, logLevel string, logOutFormat string, maxSize int, maxAge int, maxBackups int, localTime bool, compress bool) zapcore.Core {

	writeSyncer := setZapLoggerLogFileWriter(logFile, maxSize, maxAge, maxBackups, localTime, compress)
	loggerLever := getZapLoggerLevel(logLevel)
	encoder := setZapLoggerEncoder(logOutFormat)

	core := zapcore.NewCore(encoder, writeSyncer, loggerLever)

	return core
}

// Debug 使用方法：log.Debug("test")
func Debug(args ...interface{}) {

	sugaredLogger.Debug(args...)

}

func Info(args ...interface{}) {

	sugaredLogger.Info(args...)

}

func Warn(args ...interface{}) {

	sugaredLogger.Warn(args...)

}

func Error(args ...interface{}) {

	sugaredLogger.Error(args...)
}

func DPanic(args ...interface{}) {

	sugaredLogger.Panic(args)

}

func Panic(args ...interface{}) {

	sugaredLogger.Panic(args...)

}

func Fatal(args ...interface{}) {

	sugaredLogger.Fatal(args...)

}

// format output logs

func Debugf(template string, args ...interface{}) {

	sugaredLogger.Debugf(template, args...)

}

func Infof(template string, args ...interface{}) {

	sugaredLogger.Infof(template, args...)

}

func Warnf(template string, args ...interface{}) {

	sugaredLogger.Warnf(template, args...)

}

func Errorf(template string, args ...interface{}) {

	sugaredLogger.Errorf(template, args...)

}

func DPanicf(template string, args ...interface{}) {

	sugaredLogger.DPanicf(template, args...)

}

func Panicf(template string, args ...interface{}) {

	sugaredLogger.Panicf(template, args...)

}

func Fatalf(template string, args ...interface{}) {

	sugaredLogger.Fatalf(template, args...)

}

func Debugw(msg string, keysAndValues ...interface{}) {

	sugaredLogger.Debugw(msg, keysAndValues...)

}

func Infow(msg string, keysAndValues ...interface{}) {

	sugaredLogger.Infow(msg, keysAndValues...)

}

func Warnw(msg string, keysAndValues ...interface{}) {

	sugaredLogger.Warnw(msg, keysAndValues...)

}

func Errorw(msg string, keysAndValues ...interface{}) {

	sugaredLogger.Errorw(msg, keysAndValues...)

}

func DPanicw(msg string, keysAndValues ...interface{}) {

	sugaredLogger.DPanicw(msg, keysAndValues...)

}

func Panicw(msg string, keysAndValues ...interface{}) {

	sugaredLogger.Panicw(msg, keysAndValues...)

}

func Fatalw(msg string, keysAndValues ...interface{}) {

	sugaredLogger.Fatalw(msg, keysAndValues...)

}

func Sync() {

	sugaredLogger.Sync()

}
