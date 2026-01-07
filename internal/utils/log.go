package utils

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

const (
	colorReset       = "\033[0m"
	colorRed         = "\033[31m"       // Fatal
	colorDarkOrange  = "\033[38;5;208m" // Error
	colorGreen       = "\033[32m"       // Info
	colorYellow      = "\033[33m"       // Info message
	colorCyan        = "\033[36m"       // Debug
	colorMagenta     = "\033[35m"       // File names
	colorDarkYellow  = "\033[38;5;214m" // New color for Warn
	colorBrightWhite = "\033[97m"       // Message text
)

func isDebugMode() bool {
	debugValue := os.Getenv("DEBUG")
	return debugValue == "1"
}

func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(2) // 2 levels up to get the caller function
	if !ok {
		return "unknown:0"
	}
	filename := filepath.Base(file) // Extract only the file name
	return colorMagenta + filename + ":" + strconv.Itoa(line) + colorReset
}

func Debug(format string, v ...interface{}) {
	if isDebugMode() {
		log.Printf(colorCyan+"[DEBUG] "+colorMagenta+"%s: "+colorBrightWhite+format+colorReset, append([]interface{}{getCallerInfo()}, v...)...)
	}
}

func Info(format string, v ...interface{}) {
	log.Printf(colorGreen+"[INFO] "+colorYellow+format+colorReset, v...)
}

func Warn(format string, v ...interface{}) {
	log.Printf(colorDarkYellow+"[WARN] "+colorMagenta+"%s: "+colorBrightWhite+format+colorReset, append([]interface{}{getCallerInfo()}, v...)...)
}

func Error(format string, v ...interface{}) {
	log.Printf(colorDarkOrange+"[ERROR] "+colorMagenta+"%s: "+colorBrightWhite+format+colorReset, append([]interface{}{getCallerInfo()}, v...)...)
}

func Fatal(format string, v ...interface{}) {
	log.Printf(colorRed+"[FATAL] "+colorMagenta+"%s: "+colorBrightWhite+format+colorReset, append([]interface{}{getCallerInfo()}, v...)...)
	os.Exit(1)
}
