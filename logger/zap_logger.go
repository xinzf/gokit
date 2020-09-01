package logger

import (
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
	"time"
)

var encoderConfig = zapcore.EncoderConfig{
	TimeKey:       "time",
	LevelKey:      "level",
	NameKey:       "flag",
	CallerKey:     "file",
	MessageKey:    "msg",
	StacktraceKey: "stack",
	LineEnding:    zapcore.DefaultLineEnding,
	EncodeLevel:   zapcore.CapitalLevelEncoder,
	EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006/01/02 15:04:05"))
	},
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

type zapLogger struct {
	options Options
	lg      *zap.SugaredLogger
	atom    zap.AtomicLevel
}

func newZapLogger(opt ...Option) Logger {
	logger := &zapLogger{options: newOptions(opt...)}
	logger.Init()
	return logger
}

func (l *zapLogger) Init(opt ...Option) {

	if len(opt) > 0 {
		l.options = newOptions(opt...)
	}

	var writers = []zapcore.WriteSyncer{}
	writers = append(writers, os.Stdout)
	osfileout := zapcore.AddSync(&lumberjack.Logger{
		Filename:   l.options.Filename,
		MaxAge:     l.options.MaxAge,
		MaxBackups: 3,
		MaxSize:    l.options.MaxSize,
		LocalTime:  true,
		Compress:   false,
	})

	writers = append(writers, osfileout)
	w := zapcore.NewMultiWriteSyncer(writers...)

	atom := zap.NewAtomicLevel()
	atom.SetLevel(l.transform(l.options.Level))

	var enc zapcore.Encoder
	if l.options.LogType == LogText {
		enc = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		enc = zapcore.NewJSONEncoder(encoderConfig)
	}

	core := zapcore.NewCore(enc, w, atom)

	lg := zap.New(
		core,
		zap.AddStacktrace(l.transform(defaultLogStacktrace)),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	lg = lg.With(zap.String(defaultLogProjectKey, l.options.Project))

	l.lg = lg.Sugar()
	l.atom = atom
}

//func (l *zapLogger) Sync() {
//	l.lg.Sync()
//}

func (l *zapLogger) Debug(keysAndValues ...interface{}) {
	l.lg.Debugw("", l.coupArray(keysAndValues)...)
}
func (l *zapLogger) Info(keysAndValues ...interface{}) {
	l.lg.Infow("", l.coupArray(keysAndValues)...)
}
func (l *zapLogger) Warn(keysAndValues ...interface{}) {
	l.lg.Warnw("", l.coupArray(keysAndValues)...)
}
func (l *zapLogger) Error(keysAndValues ...interface{}) {
	l.lg.Errorw("", l.coupArray(keysAndValues)...)
}
func (l *zapLogger) Panic(keysAndValues ...interface{}) {
	l.lg.Panicw("", l.coupArray(keysAndValues)...)
}
func (l *zapLogger) Fatal(keysAndValues ...interface{}) {
	l.lg.Fatalw("", l.coupArray(keysAndValues)...)
}

func (l *zapLogger) Dump(keysAndValues ...interface{}) {
	arr := l.coupArray(keysAndValues)
	for k, v := range arr {
		if k%2 == 0 {
			arr[k] = v
		} else {
			arr[k] = strings.Replace(spew.Sdump(v), "\n", "", -1)
		}
	}
	l.lg.Debugw("Dump", arr...)
}

func (l *zapLogger) Print(v ...interface{}) {
	l.Info("DB", v)
}

func (l *zapLogger) Output(calldepth int, s string) error {
	l.Info("MongoDB", s)
	return nil
}

//拼接完整的数组
func (zapLogger) coupArray(kv []interface{}) []interface{} {
	if len(kv)%2 != 0 {
		kv = append(kv, kv[len(kv)-1])
		kv[len(kv)-2] = "default"
	}
	return kv
}

func (logger *zapLogger) transform(l LogLevel) zapcore.Level {
	mp := make(map[LogLevel]zapcore.Level)
	{
		mp[DebugLevel] = zapcore.DebugLevel
		mp[InfoLevel] = zapcore.InfoLevel
		mp[WarnLevel] = zapcore.WarnLevel
		mp[ErrorLevel] = zapcore.ErrorLevel
		mp[PanicLevel] = zapcore.PanicLevel
		mp[FatalLevel] = zap.FatalLevel
	}

	return mp[l]
}
