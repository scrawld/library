package loops

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestLoops_AddFuncAndStart(t *testing.T) {
	var counter int32

	// 创建 Loops 实例
	loop := New()

	// 添加一个每 100ms 运行一次的任务
	loop.AddFunc(100*time.Millisecond, func() {
		atomic.AddInt32(&counter, 1)
	})

	// 启动任务
	loop.Start()

	// 运行一段时间后停止
	time.Sleep(500 * time.Millisecond)

	// 等待任务完全停止
	ctx := loop.Stop()
	select {
	case <-ctx.Done():
		// 任务应已停止
	case <-time.After(1 * time.Second):
		t.Fatal("tasks did not stop in time")
	}

	// 检查任务执行次数是否在合理范围内
	if atomic.LoadInt32(&counter) < 4 || atomic.LoadInt32(&counter) > 6 {
		t.Fatalf("expected counter to be between 4 and 6, got %d", counter)
	}
}

func TestLoops_MultipleStart(t *testing.T) {
	var counter int32

	// 创建 Loops 实例
	loop := New()

	// 添加一个每 100ms 运行一次的任务
	loop.AddFunc(100*time.Millisecond, func() {
		atomic.AddInt32(&counter, 1)
	})

	// 多次调用 Start，确保任务只被执行一次
	loop.Start()
	loop.Start()

	// 运行一段时间后停止
	time.Sleep(500 * time.Millisecond)

	// 等待任务完全停止
	ctx := loop.Stop()
	select {
	case <-ctx.Done():
		// 任务应已停止
	case <-time.After(1 * time.Second):
		t.Fatal("tasks did not stop in time")
	}

	// 检查任务执行次数是否在合理范围内
	if atomic.LoadInt32(&counter) < 4 || atomic.LoadInt32(&counter) > 6 {
		t.Fatalf("expected counter to be between 4 and 6, got %d", counter)
	}
}

func TestLoops_StopWithoutStart(t *testing.T) {
	// 创建 Loops 实例
	loop := New()

	// 直接调用 Stop，不应该产生错误或阻塞
	ctx := loop.Stop()
	select {
	case <-ctx.Done():
		// 成功停止，无需任何操作
	case <-time.After(1 * time.Second):
		t.Fatal("stop without start caused a delay")
	}
}
