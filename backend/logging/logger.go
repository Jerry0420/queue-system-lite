package logging

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type LoggerTool interface {
	parseContextValue(context context.Context) []interface{}
	generateLogContent(level string, message []interface{}) []interface{}
	INFOf(message ...interface{})
	WARNf(message ...interface{})
	ERRORf(message ...interface{})
	FATALf(message ...interface{})
}

type loggerTool struct {
	sync.Mutex
	logger         *log.Logger
	logFile        *os.File
	contextKeys    []string
	receiveMessage chan bool
}

func getLogFilePathOfToday() string {
	// save log file inside "logs" dir.
	currentPath, _ := filepath.Abs("./")
	logBaseDir := filepath.Join(currentPath, "logs")

	err := os.MkdirAll(logBaseDir, 0777)
	if err != nil {
		log.Fatalf("Fail to getLogFilePath :%v", err)
	}
	return fmt.Sprintf("%s/%s.log", logBaseDir, time.Now().Format("20060102"))
}

func rotateLogFile(lt *loggerTool) {
	for {
		<-lt.receiveMessage
		// daily rotate the log file.
		today := fmt.Sprint(time.Now().Format("20060102"))
		if !strings.Contains(lt.logFile.Name(), today) {
			lt.Lock()
			logFilePath := getLogFilePathOfToday()
			newlogFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatalf("Fail to open log file :%v", err)
			}

			lt.logger.SetOutput(io.MultiWriter(os.Stdout, newlogFile))

			err = lt.logFile.Close()
			if err != nil {
				log.Fatalf("Fail to close old log file :%v", err)
			}
			lt.logFile = newlogFile
			defer lt.Unlock()
		}
	}
}

func NewLogger(contextKeys []string, disable bool) *loggerTool {
	var lt *loggerTool

	if disable == true {
		lt = &loggerTool{
			logger:         log.New(os.Stdout, "", 0),
			receiveMessage: make(chan bool, 10000),
		}
		return lt
	}

	logFilePath := getLogFilePathOfToday()
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to open log file :%v", err)
	}

	lt = &loggerTool{
		logger:         log.New(io.MultiWriter(os.Stdout, logFile), "", log.Ldate|log.Ltime),
		logFile:        logFile,
		contextKeys:    contextKeys,
		receiveMessage: make(chan bool, 10000),
	}

	go rotateLogFile(lt)

	return lt
}

func (lt *loggerTool) parseContextValue(context context.Context) []interface{} {
	var values []interface{}
	for _, key := range lt.contextKeys {
		// add pipe symbol for seperation.
		if key == "sep" {
			values = append(values, "|")
			continue
		}
		if value := context.Value(key); value != nil {
			values = append(values, value)
		}
	}
	return values
}

func (lt *loggerTool) generateLogContent(level string, message []interface{}) []interface{} {
	var content []interface{}
	_, file, line, _ := runtime.Caller(2)
	content = append(content, fmt.Sprint(level, file, ":", line, " |"))

	if len(message) == 0 {
		return content
	}

	if firstMessage, ok := message[0].(context.Context); ok {
		contextValues := lt.parseContextValue(firstMessage)
		content = append(content, contextValues...)
		if len(message) > 1 {
			message = message[1:]
		} else {
			return content
		}
	}

	content = append(content, "|")

	if formatString, ok := message[0].(string); ok && len(message) > 1 {
		content = append(content, fmt.Sprintf(formatString, message[1:]...))
	} else {
		content = append(content, message...)
	}

	return content
}

func (lt *loggerTool) INFOf(message ...interface{}) {
	lt.receiveMessage <- true
	content := lt.generateLogContent("[INFO] ", message)
	lt.logger.Println(content...)
}

func (lt *loggerTool) WARNf(message ...interface{}) {
	lt.receiveMessage <- true
	content := lt.generateLogContent("[WARN] ", message)
	lt.logger.Println(content...)
}

func (lt *loggerTool) ERRORf(message ...interface{}) {
	lt.receiveMessage <- true
	content := lt.generateLogContent("[ERROR] ", message)
	lt.logger.Println(content...)
}

func (lt *loggerTool) FATALf(message ...interface{}) {
	lt.receiveMessage <- true
	content := lt.generateLogContent("[FATAL] ", message)
	lt.logger.Fatalln(content...)
	// TODO: (Prod env) Send an email to the admin user.
}
