package logger

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Level byte

const (
	LevelDebug Level = iota
	LevelInfo
	LevelError
	LevelFatal
	LevelOff
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

var l *Logger

func init() {
	l = New(os.Stdout, LevelDebug)
}

func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

func PrintDebug(message string, properties map[string]any) {
	l.print(LevelDebug, message, properties)
}
func PrintInfo(message string, properties map[string]any) {
	l.print(LevelInfo, message, properties)
}
func PrintError(err error, properties map[string]any) {
	l.print(LevelError, err.Error(), properties)
}
func PrintFatal(err error, properties map[string]any) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

func (l *Logger) print(level Level, message string, properties map[string]any) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	out := struct {
		Level      string         `json:"level"`
		Time       string         `json:"time"`
		Message    string         `json:"message,omitempty"`
		Properties map[string]any `json:"properties,omitempty"`
		Trace      string         `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	if level >= LevelFatal {
		out.Trace = string(debug.Stack())
	}

	var line []byte

	line, err := json.Marshal(out)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message: " + err.Error())
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))
}

func (l *Logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}
