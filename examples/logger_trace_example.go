package main

import (
	"context"

	"github.com/chzealot/gobase/logger"
)

func main() {
	// 初始化 logger，开启调试模式以便在控制台查看输出
	err := logger.InitWithConfig(logger.Config{
		AppName:   "logger-trace-example",
		DebugMode: logger.DebugModeOn,
	})
	if err != nil {
		panic(err)
	}

	// 示例 1: 带有完整 trace_id 和 span_id 的日志
	println("\n=== 示例 1: 带有完整 trace_id 和 span_id ===")
	ctx1 := context.Background()
	ctx1 = logger.WithTrace(ctx1, "trace-abc-123-def-456", "span-001")
	logger.InfofCtx(ctx1, "用户登录请求, uid=%s", "9527")

	logger.InfowCtx(ctx1, "用户登录请求",
		"user_id", 1001,
		"username", "john_doe",
		"ip_address", "192.168.1.100",
		"action", "login",
		"status", "success")

	logger.InfowCtx(ctx1, "用户查询数据库",
		"user_id", 1001,
		"table", "users",
		"query_time_ms", 45,
		"rows_returned", 1)

	// 示例 2: 只有 trace_id
	println("\n=== 示例 2: 只有 trace_id ===")
	ctx2 := context.Background()
	ctx2 = logger.WithTraceID(ctx2, "trace-xyz-789")

	logger.InfowCtx(ctx2, "API 请求处理",
		"endpoint", "/api/v1/products",
		"method", "GET",
		"response_code", 200)

	// 示例 3: 不带 trace 信息（trace_id 和 span_id 为空）
	println("\n=== 示例 3: 不带 trace 信息 ===")
	ctx3 := context.Background()

	logger.InfowCtx(ctx3, "系统启动事件",
		"component", "main",
		"version", "1.0.0",
		"environment", "production")

	// 示例 4: 不同级别的日志
	println("\n=== 示例 4: 不同级别的日志（都带 trace）===")
	ctx4 := logger.WithTrace(context.Background(), "trace-multi-level", "span-multi")

	logger.DebugwCtx(ctx4, "调试信息",
		"debug_flag", true,
		"memory_usage_mb", 128)

	logger.InfowCtx(ctx4, "信息日志",
		"operation", "data_sync",
		"records_processed", 1000)

	logger.WarnwCtx(ctx4, "警告信息",
		"warning_type", "rate_limit_approaching",
		"current_rate", 95,
		"limit", 100)

	logger.ErrorwCtx(ctx4, "错误信息",
		"error_type", "database_connection",
		"retry_count", 3,
		"last_error", "connection timeout")

	// 示例 5: 格式化日志（*fCtx 函数）
	println("\n=== 示例 5: 格式化日志（*fCtx 系列函数）===")
	ctx5 := logger.WithTrace(context.Background(), "trace-formatted", "span-fmt")

	// InfofCtx - 格式化信息日志
	logger.InfofCtx(ctx5, "用户 %s 执行了 %s 操作", "john_doe", "login")
	logger.InfofCtx(ctx5, "处理了 %d 条记录，耗时 %d ms", 100, 250)

	// DebugfCtx - 格式化调试日志
	logger.DebugfCtx(ctx5, "调试信息: 变量值=%d, 状态=%s, 完成度=%.2f%%", 42, "active", 87.5)

	// WarnfCtx - 格式化警告日志
	logger.WarnfCtx(ctx5, "警告: CPU 使用率 %.2f%%, 超过阈值 %d%%", 95.5, 90)

	// ErrorfCtx - 格式化错误日志
	logger.ErrorfCtx(ctx5, "错误: 连接 %s:%d 失败，重试次数 %d", "database.example.com", 3306, 3)

	// 示例 6: 模拟 HTTP 请求链路追踪
	println("\n=== 示例 6: 模拟 HTTP 请求处理流程 ===")
	requestCtx := logger.WithTrace(context.Background(),
		"req-2024-12-05-abc123",
		"span-http-handler")

	logger.InfowCtx(requestCtx, "收到 HTTP 请求",
		"method", "POST",
		"path", "/api/orders",
		"client_ip", "10.0.1.50")

	// 模拟调用服务层（继承 trace）
	serviceCtx := logger.WithSpanID(requestCtx, "span-order-service")
	logger.InfowCtx(serviceCtx, "订单服务处理",
		"order_id", "ORD-20241205-001",
		"user_id", 2001,
		"amount", 299.99)

	// 模拟数据库操作（继承 trace）
	dbCtx := logger.WithSpanID(requestCtx, "span-database")
	logger.InfowCtx(dbCtx, "数据库插入操作",
		"table", "orders",
		"operation", "INSERT",
		"execution_time_ms", 23)

	logger.InfowCtx(requestCtx, "HTTP 请求处理完成",
		"status_code", 201,
		"response_time_ms", 156)

	println("\n=== 所有示例运行完成 ===")
	println("提示: 日志也会输出到文件 ~/logs/logger-trace-example/main.log")
}
