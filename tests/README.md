# 集成测试说明

本目录包含 gobase 项目的集成测试相关文档和工具。

## 概览

gobase 项目现已支持完整的集成测试能力，包括：

1. **单元测试**: 各个包内的 `*_test.go` 文件
2. **集成测试示例**: `examples/` 目录下的可执行示例程序
3. **测试工具**: 用于测试的辅助函数和工具

## Logger 包集成测试

### 测试文件位置

- 单元测试: `logger/init_test.go`
- 示例程序: `examples/logger_trace_example.go`

### 运行测试

#### 1. 运行 Logger 单元测试

```bash
# 运行所有 logger 测试
go test -v ./logger/...

# 只运行 trace ID 相关测试
go test -v ./logger/... -run TestInfowCtxWithTraceID

# 运行 context 辅助函数测试
go test -v ./logger/... -run TestContextHelpers

# 运行性能基准测试
go test -bench=. ./logger/...
```

#### 2. 运行示例程序

```bash
# 运行 logger trace ID 示例
go run examples/logger_trace_example.go

# 查看生成的日志文件
tail -f ~/logs/logger-trace-example/main.log

# 或者使用 cat 查看完整内容
cat ~/logs/logger-trace-example/main.log
```

### 测试覆盖的功能

#### InfowCtx 方法测试

`logger.InfowCtx()` 方法支持在日志中自动添加 trace_id 和 span_id，测试覆盖：

1. **完整 trace 信息**: 同时包含 trace_id 和 span_id
2. **仅 trace_id**: 只设置 trace_id，span_id 为空
3. **无 trace 信息**: 两者都为空字符串
4. **多种日志级别**: DebugwCtx, InfowCtx, WarnwCtx, ErrorwCtx

#### Context 辅助函数测试

测试以下辅助函数：

- `WithTraceID(ctx, traceID)`: 向 context 添加 trace_id
- `WithSpanID(ctx, spanID)`: 向 context 添加 span_id
- `WithTrace(ctx, traceID, spanID)`: 同时添加两者
- `GetTraceID(ctx)`: 从 context 获取 trace_id
- `GetSpanID(ctx)`: 从 context 获取 span_id

### 示例日志输出

运行示例程序后，日志输出格式如下：

```
2025-12-05T14:57:05.949+0800	INFO	examples/logger_trace_example.go:23	用户登录请求	{"trace_id": "trace-abc-123-def-456", "span_id": "span-001", "user_id": 1001, "username": "john_doe", "ip_address": "192.168.1.100", "action": "login", "status": "success"}
```

**日志字段说明**:

- 时间戳: `2025-12-05T14:57:05.949+0800` (ISO8601 格式)
- 日志级别: `INFO`
- 位置: `examples/logger_trace_example.go:23` (文件名:行号)
- 消息: `用户登录请求`
- 结构化字段: JSON 格式，包含 `trace_id`, `span_id` 和其他业务字段

### 使用场景

#### 1. 分布式追踪

在微服务架构中，使用 trace_id 追踪请求在多个服务间的调用链路：

```go
// HTTP 请求处理
func handleRequest(w http.ResponseWriter, r *http.Request) {
    // 从请求头获取或生成 trace_id
    traceID := r.Header.Get("X-Trace-ID")
    if traceID == "" {
        traceID = generateTraceID()
    }

    ctx := logger.WithTraceID(r.Context(), traceID)

    // 记录请求日志
    logger.InfowCtx(ctx, "收到请求",
        "method", r.Method,
        "path", r.URL.Path)

    // 调用服务，传递 context
    result, err := service.Process(ctx, data)

    // 记录响应日志（自动包含 trace_id）
    logger.InfowCtx(ctx, "请求完成",
        "status", "success")
}
```

#### 2. 调用链不同阶段

使用不同的 span_id 标识同一请求的不同处理阶段：

```go
func ProcessOrder(ctx context.Context, order Order) error {
    // HTTP 层
    httpCtx := logger.WithSpanID(ctx, "span-http")
    logger.InfowCtx(httpCtx, "HTTP 层处理")

    // 服务层
    serviceCtx := logger.WithSpanID(ctx, "span-service")
    logger.InfowCtx(serviceCtx, "服务层处理")

    // 数据库层
    dbCtx := logger.WithSpanID(ctx, "span-database")
    logger.InfowCtx(dbCtx, "数据库操作")

    return nil
}
```

#### 3. 错误追踪

使用 trace_id 关联同一请求的所有日志，便于错误排查：

```go
func handleAPICall(ctx context.Context) error {
    ctx = logger.WithTrace(ctx, generateTraceID(), "span-api")

    logger.InfowCtx(ctx, "开始处理 API 请求")

    if err := validateInput(ctx); err != nil {
        logger.ErrorwCtx(ctx, "输入验证失败", "error", err)
        return err
    }

    if err := processData(ctx); err != nil {
        logger.ErrorwCtx(ctx, "数据处理失败", "error", err)
        return err
    }

    logger.InfowCtx(ctx, "API 请求处理成功")
    return nil
}
```

## 添加新的集成测试

### 1. 添加包级测试

在目标包目录下创建 `*_test.go` 文件：

```go
package yourpackage

import "testing"

func TestYourFeature(t *testing.T) {
    // 测试代码
}
```

### 2. 添加集成测试示例

在 `examples/` 目录创建可执行示例：

```go
package main

import "github.com/chzealot/gobase/yourpackage"

func main() {
    // 示例代码
}
```

### 3. 运行测试

```bash
# 运行特定包的测试
go test -v ./yourpackage/...

# 运行所有测试
go test -v ./...

# 生成测试覆盖率报告
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 测试最佳实践

1. **测试命名**: 使用 `Test` 前缀，清晰描述测试内容
2. **子测试**: 使用 `t.Run()` 组织相关测试用例
3. **测试隔离**: 每个测试应该独立，不依赖其他测试的状态
4. **清理资源**: 使用 `defer` 或 `t.Cleanup()` 清理测试资源
5. **表驱动测试**: 对于多个类似场景，使用表驱动测试方式
6. **错误检查**: 使用 `t.Errorf()` 而不是 `panic()` 报告错误
7. **基准测试**: 使用 `Benchmark` 前缀测试性能

## 持续集成

可以在 CI/CD 流程中添加以下命令：

```bash
# 运行所有测试
go test -v ./...

# 检查测试覆盖率
go test -cover ./... -coverprofile=coverage.out

# 运行性能测试
go test -bench=. ./...

# 检查竞态条件
go test -race ./...
```

## 相关资源

- [Go Testing 官方文档](https://golang.org/pkg/testing/)
- [Go Test 示例](https://go.dev/blog/examples)
- [Table Driven Tests](https://go.dev/wiki/TableDrivenTests)
