// Package log contains a uber-basic logging framework with 2 levels: Info and Debug
// Inspired by:
// 	https://forum.golangbridge.org/t/whats-so-bad-about-the-stdlibs-log-package/1435
//	https://dave.cheney.net/2015/11/05/lets-talk-about-logging
package log

import (
	"fmt"
	"io/ioutil"
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

// logging levels
var (
	Debug *log.Logger
	Info  *log.Logger
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

	Debug = log.New(
		ioutil.Discard, // discard debug by default
		fmt.Sprint("[", Red, "DBG", Reset, "] "),
		log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(
		os.Stderr,
		fmt.Sprint("[", Green, "INF", Reset, "] "),
		log.Ldate|log.Ltime|log.Lshortfile)
}
