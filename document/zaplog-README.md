# zaplog

zaplog 是一个基于 [zap](https://github.com/uber-go/zap) 的高性能日志库，支持灵活的日志配置和分布式追踪功能，适用于高并发应用。

## 功能特色

- **高性能**：基于 Zap 提供极快的日志写入速度。
- **分布式追踪**：支持日志中注入分布式追踪上下文，便于问题排查和全链路分析。
- **灵活配置**：支持自定义日志级别、日志格式和输出目标。
- **模块化封装**：便于扩展和维护，支持快速接入项目。

## 快速开始

### 使用示例

#### 基本日志

```go
package main

import (
	"fmt"
	"os"

	"yourpackage/zaplog"
)

func main() {
	err := zaplog.RegisterGlobalLogger(zaplog.Config{
		Level:        "debug",
		Encoding:     "console",
		Directory:    "/tmp/zaplog",
		MaxAge:       7,
		LogInConsole: true,
	})
	if err != nil {
		fmt.Printf("register global logger error: %s\n", err)
		os.Exit(1)
	}
	defer zaplog.Sync()

	zaplog.Info("应用已启动")
}
```

#### 启用追踪日志

```go
package main

import (
	"fmt"
	"os"

	"github.com/scrawld/zaplog"
)

func main() {
	err := zaplog.RegisterGlobalLogger(zaplog.Config{
		Level:        "debug",
		Encoding:     "console",
		Directory:    "/tmp/zaplog",
		MaxAge:       7,
		LogInConsole: true,
	})
	if err != nil {
		fmt.Printf("register global logger error: %s\n", err)
		os.Exit(1)
	}
	defer zaplog.Sync()

	tracingLogger := zaplog.New()

	tracingLogger.Info("Tracing log example")
	tracingLogger.Error("Tracing error log example")
}
```
