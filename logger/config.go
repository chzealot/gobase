package logger

type DebugModeType int

const (
	DebugModeFromEnv DebugModeType = iota
	DebugModeOff
	DebugModeOn
)

type Config struct {
	AppName   string
	DebugMode DebugModeType
}
