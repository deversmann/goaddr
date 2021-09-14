package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
)

// colors
var (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
)

func init() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
		Gray = ""
		White = ""
	}
}

// logging levels
var (
	Debug = log.New(io.Discard, fmt.Sprint("[", Red, "DBG", Reset, "] "), log.Ldate|log.Ltime|log.Lshortfile)
	Info  = log.New(os.Stderr, fmt.Sprint("[", Green, "INF", Reset, "] "), log.Ldate|log.Ltime|log.Lshortfile)
)
