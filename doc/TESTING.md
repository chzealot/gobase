# 测试说明

## 快速开始

### 运行所有测试

```bash
go test -v ./...
```

### 运行 Logger Trace ID 测试

```bash
# 运行测试
go test -v ./logger/... -run TestInfowCtxWithTraceID

# 运行示例程序查看实际日志输出
go run examples/logger_trace_example.go

# 查看生成的日志文件
tail -f ~/logs/logger-trace-example/main.log
```

### 查看测试覆盖率

```bash
go test -cover ./...
```

## 测试覆盖情况

### Logger 包 (80.8% 覆盖率)

**完全覆盖的功能 (100%)**:

- `InfowCtx` - 带 context 的结构化日志（支持 trace_id 和 span_id）
- `DebugwCtx`, `WarnwCtx`, `ErrorwCtx` - 其他日志级别的结构化日志 Ctx 方法
- `InfofCtx` - 带 context 的格式化日志（支持 trace_id 和 span_id）
- `DebugfCtx`, `WarnfCtx`, `ErrorfCtx` - 其他日志级别的格式化日志 Ctx 方法
- `WithTraceID`, `WithSpanID`, `WithTrace` - Context 辅助函数
- `GetTraceID`, `GetSpanID` - 从 context 获取 trace 信息

### 测试文件

- `logger/init_test.go` - Logger 包的单元测试和集成测试
- `examples/logger_trace_example.go` - 实际使用示例

## 示例日志输出

```
2025-12-05T14:57:05.949+0800	INFO	examples/logger_trace_example.go:23	用户登录请求	{"trace_id": "trace-abc-123-def-456", "span_id": "span-001", "user_id": 1001, "username": "john_doe", "ip_address": "192.168.1.100", "action": "login", "status": "success"}
```

## 使用场景

### 基本用法

```go
import (
    "context"
    "github.com/chzealot/gobase/logger"
)

func main() {
    // 初始化 logger
    logger.InitWithConfig(logger.Config{
        AppName: "myapp",
    })

    // 创建带有 trace_id 的 context
    ctx := logger.WithTrace(context.Background(), "trace-123", "span-001")

    // 方式 1: 使用结构化日志（*wCtx 函数）
    logger.InfowCtx(ctx, "用户操作",
        "user_id", 1001,
        "action", "login")

    // 方式 2: 使用格式化日志（*fCtx 函数）
    logger.InfofCtx(ctx, "用户 %d 执行了 %s 操作", 1001, "login")
}
```

### 格式化日志 vs 结构化日志

Logger 包提供了两种日志记录方式：

**1. 格式化日志（`*fCtx` 函数）**
- 使用 printf 风格的格式化字符串
- trace_id 和 span_id 以文本形式拼接在消息前面
- 适合简单的日志消息和人类可读的输出
- 函数：`DebugfCtx`, `InfofCtx`, `WarnfCtx`, `ErrorfCtx`

```go
logger.InfofCtx(ctx, "用户 %s 登录成功, uid=%d", "john_doe", 9527)
// 输出: trace_id=xxx, span_id=yyy, 用户 john_doe 登录成功, uid=9527
```

**2. 结构化日志（`*wCtx` 函数）**
- 使用键值对记录结构化数据
- trace_id 和 span_id 作为 JSON 字段
- 适合需要查询和分析的日志
- 函数：`DebugwCtx`, `InfowCtx`, `WarnwCtx`, `ErrorwCtx`

```go
logger.InfowCtx(ctx, "用户登录成功",
    "username", "john_doe",
    "user_id", 9527)
// 输出: 用户登录成功 {"trace_id": "...", "span_id": "...", "username": "john_doe", "user_id": 9527}
```

**两者的关键区别**：
- **格式化日志**: trace_id 和 span_id 在消息文本中 → `trace_id=xxx, span_id=yyy, 消息内容`
- **结构化日志**: trace_id 和 span_id 在 JSON 字段中 → `消息内容 {"trace_id": "xxx", "span_id": "yyy", ...}`

**选择建议**:
- 简单通知消息、人类可读的日志：使用格式化日志 (`*fCtx`)
- 需要后续查询分析的数据、机器解析：使用结构化日志 (`*wCtx`)

### 分布式追踪

```go
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    // 从请求头获取或生成 trace_id
    traceID := r.Header.Get("X-Trace-ID")
    if traceID == "" {
        traceID = generateTraceID()
    }

    ctx := logger.WithTraceID(r.Context(), traceID)

    logger.InfowCtx(ctx, "收到请求", "method", r.Method, "path", r.URL.Path)

    // 处理请求...
    processRequest(ctx, r)

    logger.InfowCtx(ctx, "请求完成", "status", 200)
}

func processRequest(ctx context.Context, r *http.Request) {
    // 使用不同的 span_id 标识不同处理阶段
    serviceCtx := logger.WithSpanID(ctx, "span-service")
    logger.InfowCtx(serviceCtx, "服务层处理")

    dbCtx := logger.WithSpanID(ctx, "span-database")
    logger.InfowCtx(dbCtx, "数据库操作")
}
```

## 更多文档

详细的测试文档请参考 [tests/README.md](tests/README.md)
