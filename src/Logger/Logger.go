package Logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	Discard = New(log.New(io.Discard, "", log.LstdFlags), Config{})
	Default = New(log.New(os.Stderr, "\r\n", log.LstdFlags), Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  Warn,
		IgnoreRecordNotFoundError: true,
		Colorful:                  true,
	})
)

// Colors
const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

const (
	// Silent silent log level
	Silent LogLevel = iota + 1
	// Error error log level
	Error
	// Warn warn log level
	Warn
	// Info info log level
	Info
	// Debug debug log level
	Debug
)

// LogLevel log level
type LogLevel int

var LogPath interface{}

// Writer log writer interface
type Writer interface {
	Printf(string, ...interface{})
}

// Config logger config
type Config struct {
	SlowThreshold             time.Duration
	Colorful                  bool
	IgnoreRecordNotFoundError bool
	ParameterizedQueries      bool
	LogLevel                  LogLevel
	Debugger                  bool
	PanelLogDir               string
}

// Interface logger interface
type Interface interface {
	Info(int, string, ...interface{})
	Success(int, string, ...interface{})
	Error(int, string, ...interface{})
	Debug(int, string, ...interface{})
}

// New initialize logger
func New(writer Writer, config Config) Interface {
	var (
		infoStr    = "%s [INFO] "
		successStr = "%s [SUCCESS] "
		errStr     = "%s [ERROR] "
		debugStr   = "%s [DEBUG] "
	)

	if config.Colorful {
		infoStr = Magenta + "[%d] " + Reset + Cyan + "[INFO] " + Reset
		successStr = Magenta + "[%d] " + Reset + Green + "[SUCCESS] " + Reset
		errStr = Magenta + "[%d] " + Reset + Red + "[ERROR] " + Reset
		debugStr = Magenta + "[%d] " + Reset + Yellow + "[DEBUG] " + Reset
	}

	return &logger{
		Writer:     writer,
		Config:     config,
		infoStr:    infoStr,
		successStr: successStr,
		errStr:     errStr,
		debugStr:   debugStr,
	}
}

type logger struct {
	Writer
	Config
	infoStr, successStr, errStr, debugStr string
}

// Info print info
func (l logger) Info(code int, msg string, data ...interface{}) {
	dt := time.Now()
	currentDate := dt.Format("2006-01-02 15:04:05")
	msgStr := fmt.Sprintf(msg, data...)
	str := "[" + currentDate + "] [INFO] [" + strconv.Itoa(code) + "] " + msgStr + "\n"
	logPath := "info.log"
	if LogPath != nil {
		logPath = LogPath.(string)
	}
	logPath = l.PanelLogDir + "/" + logPath
	infoFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		// fmt.Println(err)
		os.Exit(0)
	}
	defer func() {
		err := infoFile.Close()
		if err != nil {
			// fmt.Println(err)
			os.Exit(0)
		}
	}()

	_, err = infoFile.WriteString(str)
	if err != nil {
		fmt.Println(err)
	}

	if l.Debugger {
		var attr []interface{}
		attr = append(attr, code)
		if len(data) > 0 {
			attr = append(attr, data...)
		}
		l.Printf(l.infoStr+msg, attr...)
	}
}

// Debug print debug
func (l logger) Debug(code int, msg string, data ...interface{}) {
	if l.Debugger {
		dt := time.Now()
		currentDate := dt.Format("2006-01-02 15:04:05")
		msgStr := fmt.Sprintf(msg, data...)
		str := "[" + currentDate + "] [DEBUG] [" + strconv.Itoa(code) + "] " + msgStr + "\n"

		logPath := "debug.log"
		if LogPath != nil {
			logPath = LogPath.(string)
		}
		filePath := l.PanelLogDir + "/" + logPath
		infoFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		defer func() {
			err := infoFile.Close()
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
		}()

		_, err = infoFile.WriteString(str)
		if err != nil {
			fmt.Println(err)
		}
		var attr []interface{}
		attr = append(attr, code)
		if len(data) > 0 {
			attr = append(attr, data...)
		}
		l.Printf(l.debugStr+msg, attr...)
	}
}

// Success print warn messages
func (l logger) Success(code int, msg string, data ...interface{}) {
	dt := time.Now()
	currentDate := dt.Format("2006-01-02 15:04:05")
	msgStr := fmt.Sprintf(msg, data...)
	str := "[" + currentDate + "] [SUCCESS] [" + strconv.Itoa(code) + "] " + msgStr + "\n"

	logPath := "success.log"
	if LogPath != nil {
		logPath = LogPath.(string)
	}
	logPath = l.PanelLogDir + "/" + logPath
	infoFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer func() {
		err := infoFile.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}()

	_, err = infoFile.WriteString(str)
	if err != nil {
		fmt.Println(err)
	}
	var attr []interface{}
	attr = append(attr, code)
	if len(data) > 0 {
		attr = append(attr, data...)
	}
	l.Printf(l.successStr+msg, attr...)
}

// Error print error messages
func (l logger) Error(code int, msg string, data ...interface{}) {
	dt := time.Now()
	currentDate := dt.Format("2006-01-02 15:04:05")
	msgStr := fmt.Sprintf(msg, data...)
	str := "[" + currentDate + "] [ERROR] [" + strconv.Itoa(code) + "] " + msgStr + "\n"

	logPath := "error.log"
	if LogPath != nil {
		logPath = LogPath.(string)
	}
	logPath = l.PanelLogDir + "/" + logPath
	errorFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	defer func() {
		err := errorFile.Close()
		if err != nil {
			os.Exit(0)
		}
	}()

	_, err = errorFile.WriteString(str)
	if err != nil {
		fmt.Println(err)
	}
	var attr []interface{}
	attr = append(attr, code)
	if len(data) > 0 {
		attr = append(attr, data...)
	}
	l.Printf(l.errStr+msg, attr...)
	fmt.Fprintln(os.Stdout, fmt.Sprintf(msgStr))
}
