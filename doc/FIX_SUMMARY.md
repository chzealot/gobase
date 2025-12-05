# 修复总结

## 问题

你发现 `logger.InfofCtx(ctx1, "用户登录请求, uid=%s", "9527")` 的输出不正确：

**实际输出**：
```
用户登录请求, uid=%s    {"trace_id": "trace-abc-123-def-456", "span_id": "span-001", "args": ["9527"]}
```

**期望输出**：
```
用户登录请求, uid=9527    {"trace_id": "trace-abc-123-def-456", "span_id": "span-001"}
```

## 修复内容

### 1. 修复了 4 个格式化日志函数的实现

修改文件：`logger/init.go`

修复的函数：
- `DebugfCtx` (行 206-211)
- `InfofCtx` (行 222-227)
- `WarnfCtx` (行 238-243)
- `ErrorfCtx` (行 254-259)

**修复方法**：使用 `fmt.Sprintf(format, args...)` 先格式化消息，然后再调用 `*w` 函数输出。

### 2. 新增完整的测试用例

新增测试：`TestInfofCtx` (logger/init_test.go:140-198)

测试覆盖：
- InfofCtx 格式化功能
- DebugfCtx 格式化功能
- WarnfCtx 格式化功能
- ErrorfCtx 格式化功能
- 不带 trace 信息的格式化日志

### 3. 更新示例程序

文件：`examples/logger_trace_example.go`

新增示例 5：演示所有 `*fCtx` 函数的使用

### 4. 更新文档

更新的文档：
- `TESTING.md` - 更新测试覆盖率和使用说明
- `tests/README.md` - 添加格式化日志的使用场景
- `BUGFIX.md` - 详细的 bug 修复记录

## 测试结果

### 测试通过率：100%

```bash
$ go test -v ./logger/...
=== RUN   TestInfowCtxWithTraceID
--- PASS: TestInfowCtxWithTraceID (0.00s)
=== RUN   TestContextHelpers
--- PASS: TestContextHelpers (0.00s)
=== RUN   TestInfofCtx
--- PASS: TestInfofCtx (0.00s)
PASS
```

### 测试覆盖率提升

- **修复前**: 64.1%
- **修复后**: 80.8%
- **所有 `*fCtx` 函数**: 100% 覆盖

```bash
$ go test -cover ./logger/...
ok  	github.com/chzealot/gobase/logger	0.486s	coverage: 80.8% of statements
```

### 实际输出验证

```bash
$ go run examples/logger_trace_example.go
...
2025-12-05T15:07:21.831+0800	INFO	examples/logger_trace_example.go:23	用户登录请求, uid=9527	{"trace_id": "trace-abc-123-def-456", "span_id": "span-001"}
...
2025-12-05T15:07:28.361+0800	INFO	examples/logger_trace_example.go:84	用户 john_doe 执行了 login 操作	{"trace_id": "trace-formatted", "span_id": "span-fmt"}
2025-12-05T15:07:28.361+0800	INFO	examples/logger_trace_example.go:85	处理了 100 条记录，耗时 250 ms	{"trace_id": "trace-formatted", "span_id": "span-fmt"}
```

✅ 所有格式化都正确！

## 使用指南

### 格式化日志（`*fCtx` 函数）

适用于简单的日志消息，使用 printf 风格：

```go
logger.InfofCtx(ctx, "用户 %s 登录成功, uid=%d", "john_doe", 9527)
logger.WarnfCtx(ctx, "CPU 使用率 %.2f%%, 超过阈值 %d%%", 95.5, 90)
logger.ErrorfCtx(ctx, "连接 %s:%d 失败，重试 %d 次", "db.com", 3306, 3)
```

### 结构化日志（`*wCtx` 函数）

适用于需要查询和分析的数据：

```go
logger.InfowCtx(ctx, "用户登录成功",
    "username", "john_doe",
    "user_id", 9527,
    "ip", "192.168.1.1")
```

### 选择建议

- **简单通知消息** → 使用格式化日志 (`*fCtx`)
- **需要后续查询分析的数据** → 使用结构化日志 (`*wCtx`)

## 相关文件

### 修改的文件
- `logger/init.go` - 修复 4 个函数的实现
- `logger/init_test.go` - 新增测试用例
- `examples/logger_trace_example.go` - 新增示例 5

### 新增的文档
- `BUGFIX.md` - 详细的 bug 修复记录
- `FIX_SUMMARY.md` - 本文件

### 更新的文档
- `TESTING.md` - 更新测试覆盖率和使用说明
- `tests/README.md` - 更新测试指南

## 下一步

建议在其他使用了这些函数的代码中：

1. **检查使用方式**：确保使用的是正确的函数（`*fCtx` 用于格式化，`*wCtx` 用于结构化）
2. **运行测试**：验证现有代码在修复后工作正常
3. **更新日志**：如果有需要，调整日志输出格式

---

修复完成时间：2025-12-05
