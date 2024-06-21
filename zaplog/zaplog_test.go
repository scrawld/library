package zaplog

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := RegisterGlobalLogger(Config{
		Level:        "debug",
		Encoding:     "console",
		Directory:    "/tmp/zaplog",
		MaxAge:       7,
		LogInConsole: true,
		//FileNames: map[string]string{
		//	"main":  "myname.log",
		//	"warn":  "myname.warn.log",
		//	"error": "myname.err.log",
		//},
	})
	if err != nil {
		fmt.Printf("register global logger error: %s\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestTracingLogger(t *testing.T) {
	// 创建一个新的 TracingLogger 实例
	logger := New()

	// 测试 Debug 方法
	logger.Debug("This is a debug message.")

	// 测试 Info 方法
	logger.Info("This is an info message.")

	// 测试 Warn 方法
	logger.Warn("This is a warning message.")

	// 测试 Error 方法
	logger.Error("This is an error message.")

	// 测试 Debugf 方法
	logger.Debugf("This is a formatted debug message with %s.", "arguments")

	// 测试 Infof 方法
	logger.Infof("This is a formatted info message with %s.", "arguments")

	// 测试 Warnf 方法
	logger.Warnf("This is a formatted warning message with %s.", "arguments")

	// 测试 Errorf 方法
	logger.Errorf("This is a formatted error message with %s.", "arguments")
}
