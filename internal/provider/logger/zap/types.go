package zap

type config interface {
	GetLogFilePath() string

	GetLogLevel() string
	SetLogLevel(string)

	GetLogFileMaxSize() int
	SetLogFileMaxSize(int)

	GetLogFileMaxAge() int
	SetLogFileMaxAge(int)

	GetLogFileMaxBackups() int
	SetLogFileMaxBackups(int)
}
