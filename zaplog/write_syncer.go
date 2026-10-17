package zaplog

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap/zapcore"
)

// GetWriteSyncer 获取日志输出目标,按日分割
func GetWriteSyncer(filename string, maxAge int, logInConsole bool) (zapcore.WriteSyncer, error) {
	if maxAge < 0 {
		return nil, errors.New("maxAge must be a non-negative value")
	}
	fileWriter, err := rotatelogs.New(
		AddDateToFilename(filename, "%Y-%m-%d"),                   // 日志文件名格式，支持时间格式化
		rotatelogs.WithLinkName(filename),                         // 生成软连接
		rotatelogs.WithMaxAge(time.Duration(maxAge)*24*time.Hour), // 日志留存时间
		rotatelogs.WithRotationTime(time.Hour*24),                 // 日志切割时间间隔
	)
	if err != nil {
		return nil, fmt.Errorf("rotatelogs new error: %s", err)
	}
	if logInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), nil
	}
	return zapcore.AddSync(fileWriter), nil
}

// AddDateToFilename creates a new filename from the given name, inserting a date
func AddDateToFilename(name string, date string) string {
	var (
		dir      = filepath.Dir(name)
		filename = filepath.Base(name)
		ext      = filepath.Ext(filename)
		prefix   = strings.TrimSuffix(filename, ext)
	)
	return filepath.Join(dir, fmt.Sprintf("%s.%s%s", prefix, date, ext))
}

/* 使用lumberjack按文件大小做日志分割

"gopkg.in/natefinch/lumberjack.v2"

fileWriter := &lumberjack.Logger{
	Filename:   filename, // 日志文件名
	MaxSize:    500,      // MB
	MaxAge:     maxAge,   // 控制日志文件的最大保留天数
	MaxBackups: 0,        // 控制日志文件的最大保留个数,0表示保留所有文件
	LocalTime:  true,     // 使用本地时间,false则使用UTC时间
}
go rotateLogsAtMidnight(fileWriter)

// rotateLogsAtMidnight 定时在每天零点进行日志切割操作
func rotateLogsAtMidnight(fileWriter *lumberjack.Logger) {
	for {
		now := time.Now()
		nextDayZero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1)
		after := nextDayZero.UnixNano() - now.UnixNano()
		<-time.After(time.Duration(after) * time.Nanosecond)

		fi, err := os.Stat(fileWriter.Filename)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			fmt.Printf("os stat error: %s\n", err)
			continue
		}
		if fi.IsDir() || fi.Size() == 0 {
			continue
		}
		fileWriter.Rotate()
	}
}
*/
