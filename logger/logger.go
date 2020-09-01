package logger

var (
	DefaultLogger = newZapLogger()
)

type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	PanicLevel LogLevel = "panic"
	FatalLevel LogLevel = "fatal"
)

//默认参数
const (
	defaultLogFilename   string   = "./log/default.log" //日志保存路径 //需要设置程序当前运行路径
	defaultLogLevel      LogLevel = DebugLevel          //日志记录级别
	defaultLogMaxSize    int      = 512                 //日志分割的尺寸 MB
	defaultLogMaxAge     int      = 30                  //分割日志保存的时间 day
	defaultLogStacktrace LogLevel = PanicLevel          //记录堆栈的级别
	defaultLogProjectKey string   = "project"           //
	defaultLogType                = LogText
)

type LogType string

const (
	LogJson LogType = "json"
	LogText LogType = "text"
)

type Logger interface {
	Init(opt ...Option)
	Debug(keysAndValues ...interface{})
	Info(keysAndValues ...interface{})
	Warn(keysAndValues ...interface{})
	Error(keysAndValues ...interface{})
	Panic(keysAndValues ...interface{})
	Fatal(keysAndValues ...interface{})
	Dump(keysAndValues ...interface{})
	Print(v ...interface{})
	Output(calldepth int, s string) error
}
