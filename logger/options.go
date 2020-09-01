package logger

type Options struct {
	Filename string
	Level    LogLevel
	MaxSize  int
	MaxAge   int
	LogType  LogType
	Project  string
}

type Option func(*Options)

func newOptions(opt ...Option) Options {
	opts := Options{}

	for _, o := range opt {
		o(&opts)
	}

	if opts.Filename == "" {
		opts.Filename = defaultLogFilename
	}
	if opts.Level == "" {
		opts.Level = defaultLogLevel
	}
	if opts.MaxSize == 0 {
		opts.MaxSize = defaultLogMaxSize
	}
	if opts.MaxAge == 0 {
		opts.MaxAge = defaultLogMaxAge
	}
	if opts.LogType == "" {
		opts.LogType = defaultLogType
	}
	return opts
}

func Filename(name string) Option {
	return func(options *Options) {
		options.Filename = name
	}
}

func Level(level LogLevel) Option {
	return func(options *Options) {
		options.Level = level
	}
}

func MaxSize(size int) Option {
	return func(options *Options) {
		options.MaxSize = size
	}
}

func MaxAge(age int) Option {
	return func(options *Options) {
		options.MaxAge = age
	}
}

func Type(logType LogType) Option {
	return func(options *Options) {
		options.LogType = logType
	}
}

func ProjectName(name string) Option {
	return func(options *Options) {
		options.Project = name
	}
}
