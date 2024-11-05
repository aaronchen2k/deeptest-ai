package zap_server

import (
	"fmt"
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/viper_server"
	_file "github.com/deeptest-com/deeptest-next/pkg/libs/file"
	_logUtils "github.com/deeptest-com/deeptest-next/pkg/libs/log"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init() {
	_logUtils.SetLogger(getLogger())
}

func getLogger() (logger *zap.Logger) {
	viper_server.Init(getViperConfig())

	logger, _ = createLogger()

	return
}

func createLogger() (ret *zap.Logger, err error) {
	level := getLogLevel()
	logPath := getLogPath()

	prodEncoder := getEncoderConfig()

	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= level
	})
	writeSyncer, lowClose, err := zap.Open(logPath)
	if err != nil {
		if lowClose != nil {
			lowClose()
		}
		return
	}

	swSugar := zapcore.NewMultiWriteSyncer(
		writeSyncer,
		zapcore.AddSync(os.Stdout),
	)
	infoCore := zapcore.NewCore(zapcore.NewJSONEncoder(prodEncoder), swSugar, infoPriority)

	ret = zap.New(zapcore.NewTee(infoCore))

	ret = ret.WithOptions(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel))

	return
}

func getLogPath() (ret string) {
	logsDir := filepath.Join(consts.WorkDir, CONFIG.Director)
	_file.InsureDir(logsDir)

	ret = filepath.Join(logsDir, fmt.Sprintf("%s.log", consts.App))

	return
}

func getEncoderConfig() (conf zapcore.EncoderConfig) {
	conf = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  CONFIG.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeTime:     customTimeEncoder,
	}

	switch {
	case CONFIG.EncodeLevel == "LowercaseLevelEncoder":
		conf.EncodeLevel = zapcore.LowercaseLevelEncoder
	case CONFIG.EncodeLevel == "LowercaseColorLevelEncoder":
		conf.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case CONFIG.EncodeLevel == "CapitalLevelEncoder":
		conf.EncodeLevel = zapcore.CapitalLevelEncoder
	case CONFIG.EncodeLevel == "CapitalColorLevelEncoder":
		conf.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		conf.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return conf
}

type StringsArray [][]string

// MarshalLogArray
func (ss StringsArray) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range ss {
		for ii := range ss[i] {
			arr.AppendString(ss[i][ii])
		}
	}
	return nil
}

func getLogLevel() (level zapcore.Level) {
	switch CONFIG.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	return
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(CONFIG.Prefix + " 2006/01/02 15:04:05.000"))
}
