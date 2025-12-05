# gobase 文档索引

本目录包含 gobase 项目的所有文档。

## 文档列表

### 更新日志
- [CHANGELOG.md](CHANGELOG.md) - 项目更新日志
  - 功能改进记录
  - 文档重组说明

### 测试相关
- [TESTING.md](TESTING.md) - 快速测试指南
  - 如何运行测试
  - Logger trace ID 功能使用示例
  - 格式化日志 vs 结构化日志
  - 测试覆盖率

- [../tests/README.md](../tests/README.md) - 详细的集成测试文档
  - 完整的测试说明
  - 使用场景示例
  - 测试最佳实践

### Bug 修复
- [BUGFIX.md](BUGFIX.md) - 详细的 bug 修复记录
  - 2025-12-05: 修复格式化日志函数（*fCtx）的实现
  - 问题描述、根本原因、修复方案
  - 测试覆盖和验证结果

- [FIX_SUMMARY.md](FIX_SUMMARY.md) - 修复总结
  - 修复内容概览
  - 测试结果
  - 使用指南

### 开发指南
- [../CLAUDE.md](../CLAUDE.md) - Claude Code 开发指南
  - 项目概述
  - 模块结构
  - 开发命令
  - 常用模式

## 主要功能文档

### Logger 包

Logger 包提供了带有 trace ID 支持的日志功能。

**两种日志方式**：

1. **格式化日志 (`*fCtx` 函数)**
   - `DebugfCtx`, `InfofCtx`, `WarnfCtx`, `ErrorfCtx`
   - trace_id 和 span_id 在消息文本中
   - 示例：`logger.InfofCtx(ctx, "用户 %s 登录, uid=%d", "john", 9527)`
   - 输出：`trace_id=xxx, span_id=yyy, 用户 john 登录, uid=9527`

2. **结构化日志 (`*wCtx` 函数)**
   - `DebugwCtx`, `InfowCtx`, `WarnwCtx`, `ErrorwCtx`
   - trace_id 和 span_id 作为 JSON 字段
   - 示例：`logger.InfowCtx(ctx, "用户登录", "username", "john", "uid", 9527)`
   - 输出：`用户登录 {"trace_id": "xxx", "span_id": "yyy", "username": "john", "uid": 9527}`

**快速开始**：参考 [TESTING.md](TESTING.md)

## 示例代码

查看 `examples/` 目录获取完整的使用示例：

- `examples/logger_trace_example.go` - Logger trace ID 功能完整示例

## 获取帮助

如有问题，请查阅相关文档或查看测试用例了解使用方法。
