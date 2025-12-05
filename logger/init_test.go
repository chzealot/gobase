package logger

import (
	"context"
	"testing"
)

// TestInfowCtxWithTraceID 测试带有 trace_id 和 span_id 的日志输出
func TestInfowCtxWithTraceID(t *testing.T) {
	// 初始化 logger，使用测试应用名
	err := InitWithConfig(Config{
		AppName:   "gobase-test",
		DebugMode: DebugModeOn, // 开启调试模式以便在控制台看到输出
	})
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// 创建一个基础 context
	ctx := context.Background()

	// 测试用例 1: 添加 trace_id 和 span_id
	t.Run("WithTraceIDAndSpanID", func(t *testing.T) {
		testTraceID := "test-trace-id-12345"
		testSpanID := "test-span-id-67890"

		// 使用 WithTrace 同时设置 trace_id 和 span_id
		ctxWithTrace := WithTrace(ctx, testTraceID, testSpanID)

		// 验证可以正确获取 trace_id 和 span_id
		if gotTraceID := GetTraceID(ctxWithTrace); gotTraceID != testTraceID {
			t.Errorf("GetTraceID() = %v, want %v", gotTraceID, testTraceID)
		}
		if gotSpanID := GetSpanID(ctxWithTrace); gotSpanID != testSpanID {
			t.Errorf("GetSpanID() = %v, want %v", gotSpanID, testSpanID)
		}

		// 使用 InfowCtx 记录日志，应该包含 trace_id 和 span_id
		InfowCtx(ctxWithTrace, "测试日志：带有 trace_id 和 span_id",
			"user_id", 1001,
			"action", "login",
			"status", "success")

		t.Log("日志已输出，请检查控制台或日志文件确认包含 trace_id 和 span_id")
	})

	// 测试用例 2: 只添加 trace_id
	t.Run("WithTraceIDOnly", func(t *testing.T) {
		testTraceID := "test-trace-id-only-99999"

		ctxWithTrace := WithTraceID(ctx, testTraceID)

		if gotTraceID := GetTraceID(ctxWithTrace); gotTraceID != testTraceID {
			t.Errorf("GetTraceID() = %v, want %v", gotTraceID, testTraceID)
		}

		InfowCtx(ctxWithTrace, "测试日志：只有 trace_id",
			"user_id", 1002,
			"action", "logout")

		t.Log("日志已输出（只有 trace_id）")
	})

	// 测试用例 3: 不带任何 trace 信息
	t.Run("WithoutTrace", func(t *testing.T) {
		InfowCtx(ctx, "测试日志：不带 trace 信息",
			"user_id", 1003,
			"action", "view_page")

		t.Log("日志已输出（不带 trace 信息，trace_id 和 span_id 应为空字符串）")
	})

	// 测试用例 4: 多种日志级别的 Ctx 方法
	t.Run("DifferentLogLevels", func(t *testing.T) {
		testTraceID := "test-trace-id-levels"
		testSpanID := "test-span-id-levels"
		ctxWithTrace := WithTrace(ctx, testTraceID, testSpanID)

		DebugwCtx(ctxWithTrace, "Debug 级别日志", "level", "debug")
		InfowCtx(ctxWithTrace, "Info 级别日志", "level", "info")
		WarnwCtx(ctxWithTrace, "Warn 级别日志", "level", "warn")
		ErrorwCtx(ctxWithTrace, "Error 级别日志", "level", "error")

		t.Log("不同级别的日志已输出，都应该包含 trace_id 和 span_id")
	})
}

// TestContextHelpers 测试 context 辅助函数
func TestContextHelpers(t *testing.T) {
	ctx := context.Background()

	t.Run("WithTraceID", func(t *testing.T) {
		traceID := "trace-123"
		newCtx := WithTraceID(ctx, traceID)

		if got := GetTraceID(newCtx); got != traceID {
			t.Errorf("GetTraceID() = %v, want %v", got, traceID)
		}

		// 验证原始 context 未被修改
		if got := GetTraceID(ctx); got != "" {
			t.Errorf("Original context should not have trace_id, got %v", got)
		}
	})

	t.Run("WithSpanID", func(t *testing.T) {
		spanID := "span-456"
		newCtx := WithSpanID(ctx, spanID)

		if got := GetSpanID(newCtx); got != spanID {
			t.Errorf("GetSpanID() = %v, want %v", got, spanID)
		}
	})

	t.Run("WithTrace", func(t *testing.T) {
		traceID := "trace-789"
		spanID := "span-012"
		newCtx := WithTrace(ctx, traceID, spanID)

		if got := GetTraceID(newCtx); got != traceID {
			t.Errorf("GetTraceID() = %v, want %v", got, traceID)
		}
		if got := GetSpanID(newCtx); got != spanID {
			t.Errorf("GetSpanID() = %v, want %v", got, spanID)
		}
	})

	t.Run("NilContext", func(t *testing.T) {
		// 测试 nil context 的边界情况
		if got := GetTraceID(nil); got != "" {
			t.Errorf("GetTraceID(nil) should return empty string, got %v", got)
		}
		if got := GetSpanID(nil); got != "" {
			t.Errorf("GetSpanID(nil) should return empty string, got %v", got)
		}
	})
}

// TestInfofCtx 测试格式化日志函数
func TestInfofCtx(t *testing.T) {
	err := InitWithConfig(Config{
		AppName:   "gobase-test-f",
		DebugMode: DebugModeOn,
	})
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	ctx := context.Background()

	// 测试用例 1: InfofCtx 格式化字符串
	t.Run("InfofCtxFormatting", func(t *testing.T) {
		testTraceID := "test-trace-infof"
		testSpanID := "test-span-infof"
		ctxWithTrace := WithTrace(ctx, testTraceID, testSpanID)

		// 测试格式化，应该输出 "用户登录请求, uid=9527" 而不是 "用户登录请求, uid=%s"
		InfofCtx(ctxWithTrace, "用户登录请求, uid=%s", "9527")
		InfofCtx(ctxWithTrace, "用户 %s 执行了 %s 操作", "john_doe", "login")
		InfofCtx(ctxWithTrace, "处理了 %d 条记录，耗时 %d ms", 100, 250)

		t.Log("InfofCtx 格式化日志已输出，请检查控制台确认格式化正确")
	})

	// 测试用例 2: DebugfCtx 格式化
	t.Run("DebugfCtxFormatting", func(t *testing.T) {
		ctxWithTrace := WithTrace(ctx, "debug-trace", "debug-span")

		DebugfCtx(ctxWithTrace, "调试信息: 变量值=%d, 状态=%s", 42, "active")

		t.Log("DebugfCtx 格式化日志已输出")
	})

	// 测试用例 3: WarnfCtx 格式化
	t.Run("WarnfCtxFormatting", func(t *testing.T) {
		ctxWithTrace := WithTrace(ctx, "warn-trace", "warn-span")

		WarnfCtx(ctxWithTrace, "警告: CPU 使用率 %.2f%%, 超过阈值 %d%%", 95.5, 90)

		t.Log("WarnfCtx 格式化日志已输出")
	})

	// 测试用例 4: ErrorfCtx 格式化
	t.Run("ErrorfCtxFormatting", func(t *testing.T) {
		ctxWithTrace := WithTrace(ctx, "error-trace", "error-span")

		ErrorfCtx(ctxWithTrace, "错误: 连接 %s 失败，重试次数 %d", "database", 3)

		t.Log("ErrorfCtx 格式化日志已输出")
	})

	// 测试用例 5: 不带 trace 信息的格式化日志
	t.Run("FormatWithoutTrace", func(t *testing.T) {
		InfofCtx(ctx, "无 trace 信息: 用户=%s, 操作=%s", "admin", "logout")

		t.Log("无 trace 信息的格式化日志已输出")
	})
}

// BenchmarkInfowCtx 性能测试
func BenchmarkInfowCtx(b *testing.B) {
	InitWithConfig(Config{
		AppName:   "gobase-bench",
		DebugMode: DebugModeOff,
	})

	ctx := WithTrace(context.Background(), "bench-trace-id", "bench-span-id")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		InfowCtx(ctx, "benchmark log message", "iteration", i)
	}
}

// BenchmarkInfofCtx 性能测试 - 格式化日志
func BenchmarkInfofCtx(b *testing.B) {
	InitWithConfig(Config{
		AppName:   "gobase-bench-f",
		DebugMode: DebugModeOff,
	})

	ctx := WithTrace(context.Background(), "bench-trace-id", "bench-span-id")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		InfofCtx(ctx, "benchmark log message %d", i)
	}
}
