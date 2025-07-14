package logger

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// import 中的 testing 包是 Go 语言的测试包，用于编写测试代码。
// testing 包提供了一些测试相关的函数和类型，例如：
// - *testing.T：测试对象，用于记录测试结果和执行测试操作。
// - *testing.M：测试主函数，用于管理测试的生命周期。
// - *testing.B：性能测试对象，用于执行性能测试。
// - *testing.Verbosity：测试的详细程度，用于控制测试输出的详细程度。
// - *testing.Verbose：测试是否详细，用于控制测试输出的详细程度。

func TestNewLogger(t *testing.T) {
	config := &Config{
		Level:      "debug",
		Format:     "json",
		Output:     "stdout",
		Filename:   "",
		MaxSize:    100,
		MaxAge:     28,
		MaxBackups: 3,
		Compress:   true,
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	if logger == nil {
		t.Fatal("logger is nil")
	}

	//test logger console output
	logger.Info("test logger Info")
	logger.Error("test logger Error")
	logger.Warn("test logger Warn")
	logger.Debug("test logger Debug")
}

func TestLoggerWithField(t *testing.T) {
	tmpDir := t.TempDir()

	logFile := filepath.Join(tmpDir, "test.log")

	config := &Config{
		Level:      "debug",
		Format:     "json",
		Output:     "file",
		Filename:   logFile,
		MaxSize:    100,
		MaxAge:     28,
		MaxBackups: 3,
		Compress:   true,
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	if logger == nil {
		t.Fatal("logger is nil")
	}

	//test logger file output
	logger.Info("test logger Info")

	//with fields
	logger.WithField("key", "value").Info("test logger Info with fields")

	logger.Close()

	//check log file if existss
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Fatalf("log file %s does not exist", logFile)
	}

}

func TestLoggerWithFields(t *testing.T) {
	config := &Config{
		Level:  "debug",
		Format: "json",
		Output: "stdout",
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	if logger == nil {
		t.Fatal("logger is nil")
	}

	logger.WithFields(map[string]interface{}{
		"user_id":   "123",
		"username":  "test",
		"email":     "test@example.com",
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}).Info("User Login")

	logger.WithField("request_id", "123").Info("Request ID")
}
