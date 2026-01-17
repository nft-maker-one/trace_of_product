package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

var (
	logBuffer     []LogEntry
	bufferMutex   sync.RWMutex
	maxBufferSize = 1000
	logFile       *os.File
)

func InitLogger(logDir string) error {
	// 创建日志目录
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// 打开日志文件
	logPath := filepath.Join(logDir, "app.log")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	logFile = file

	// 设置logrus输出到文件和控制台
	logrus.SetOutput(file)
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	return nil
}

func addToBuffer(entry LogEntry) {
	bufferMutex.Lock()
	defer bufferMutex.Unlock()

	logBuffer = append(logBuffer, entry)
	if len(logBuffer) > maxBufferSize {
		logBuffer = logBuffer[1:]
	}
}

func GetLogs(limit int) []LogEntry {
	bufferMutex.RLock()
	defer bufferMutex.RUnlock()

	if limit <= 0 || limit > len(logBuffer) {
		limit = len(logBuffer)
	}

	start := len(logBuffer) - limit
	if start < 0 {
		start = 0
	}

	result := make([]LogEntry, limit)
	copy(result, logBuffer[start:])
	return result
}

func LogMsg(keys, values []string) error {
	fields := logrus.Fields{}
	if len(keys) != len(values) {
		logrus.Errorf("输入不符合规范，keys 长度为 %d values 长度为 %d，两者不相同", len(keys), len(values))
		return fmt.Errorf("输入不符合规范，keys 长度为 %d values 长度为 %d，两者不相同", len(keys), len(values))
	}
	for i := 0; i < len(keys); i++ {
		fields[keys[i]] = values[i]
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   "Log message",
		Fields:    fields,
	}
	addToBuffer(entry)

	logrus.WithFields(fields).Println()
	return nil
}

func LogWarn(keys, values []string) error {
	fields := logrus.Fields{}
	if len(keys) != len(values) {
		logrus.Errorf("输入不符合规范，keys 长度为 %d values 长度为 %d，两者不相同", len(keys), len(values))
		return fmt.Errorf("输入不符合规范，keys 长度为 %d values 长度为 %d，两者不相同", len(keys), len(values))
	}
	for i := 0; i < len(keys); i++ {
		fields[keys[i]] = values[i]
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     "warning",
		Message:   "Warning message",
		Fields:    fields,
	}
	addToBuffer(entry)

	logrus.WithFields(fields).Warnln()
	return nil
}

func LogError(keys, values []string) error {
	fields := logrus.Fields{}
	if len(keys) != len(values) {
		logrus.Errorf("输入不符合规范，keys 长度为 %d values 长度为 %d，两者不相同", len(keys), len(values))
		return fmt.Errorf("输入不符合规范，keys 长度为 %d values 长度为 %d，两者不相同", len(keys), len(values))
	}
	for i := 0; i < len(keys); i++ {
		fields[keys[i]] = values[i]
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     "error",
		Message:   "Error message",
		Fields:    fields,
	}
	addToBuffer(entry)

	logrus.WithFields(fields).Errorln()
	return nil
}

func LogDebug(keys, values []string) error {
	fields := logrus.Fields{}
	if len(keys) != len(values) {
		logrus.Errorf("输入不符合规范，keys 长度为 %d values 长度为 %d，两者不相同", len(keys), len(values))
		return fmt.Errorf("输入不符合规范，keys 长度为 %d values 长度为 %d，两者不相同", len(keys), len(values))
	}
	for i := 0; i < len(keys); i++ {
		fields[keys[i]] = values[i]
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     "debug",
		Message:   "Debug message",
		Fields:    fields,
	}
	addToBuffer(entry)

	logrus.WithFields(fields).Debugln()
	return nil
}
