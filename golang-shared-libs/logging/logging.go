package logging

import (
	"log"
	"os"
	"path/filepath"
)

type Logger struct {
	loggerConsole *log.Logger
	loggerFile    *log.Logger
}

func (l Logger) Print(v ...any) {
	l.loggerFile.Print(v)
	l.loggerConsole.Print(v)
}

func (l Logger) Println(v ...any) {
	l.loggerFile.Println(v)
	l.loggerConsole.Println(v)
}

func (l Logger) Fatalln(v ...any) {
	l.loggerFile.Println(v)
	l.loggerConsole.Fatalln(v)
}

func (l Logger) Fatal(v ...any) {
	l.loggerFile.Print(v)
	l.loggerConsole.Fatal(v)
}

func NewLogger(filename string) *Logger {
	abs, err := filepath.Abs(filename)
	if err != nil {
		return nil
	}

	openLogfile, err := os.OpenFile(abs, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Error opening file:", err)
	}

	fileLogger := log.New(openLogfile, "Log:\t", log.Ldate|log.Ltime|log.Lshortfile)
	consoleLogger := log.Default()

	return &Logger{
		loggerConsole: consoleLogger,
		loggerFile:    fileLogger,
	}
}
