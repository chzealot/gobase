# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is `gobase`, a Go library providing common infrastructure capabilities for Go applications. The library is under active development and may have breaking changes.

**Warning**: This library is still in development and should be used cautiously in production environments (as noted in README.md).

## Module Structure

The repository is organized into independent packages:

- `dingtalk/` - DingTalk API client for OAuth, calendar, user management, and todo tasks
- `database/` - MySQL database initialization with GORM
- `logger/` - Structured logging built on zap with log rotation
- `utils/` - HTTP utilities
- `constants/` - Shared constants (primarily HTTP headers)

Each package can be imported independently as needed.

## Development Commands

### Building and Testing

```bash
# Build the module
go build ./...

# Run tests (if any exist)
go test ./...

# Update dependencies
go mod tidy

# Check module for issues
go mod verify
```

### Go Version

This module requires Go 1.19+ (as specified in go.mod). Current development uses Go 1.25.0.

## Architecture

### DingTalk Client (`dingtalk/`)

The DingTalk client is the most complex package, providing integration with DingTalk's APIs:

- **Client initialization**: `NewDingTalkClient(clientId, clientSecret)` creates a new client
- **Token management**: Access tokens are cached with mutex protection and auto-refresh (see `GetAccessToken()` at dingtalk/client.go:207)
- **User authentication**: Supports both user access tokens and app access tokens
- **Calendar operations**: Fetches user calendars and events (currently limited to primary calendar)
- **Todo task creation**: Creates tasks for users
- **API patterns**:
  - Old API endpoints use `oapi.dingtalk.com` with query-based access tokens
  - New API endpoints use `api.dingtalk.com/v1.0` with header-based tokens (`x-acs-dingtalk-access-token`)

Key implementation notes:
- Token expiry includes 60-second buffer to avoid edge-case failures (dingtalk/client.go:213)
- Default HTTP timeout is 60 seconds
- Error responses may use `panic()` for HTTP errors (should be improved for production use)

### Database Package (`database/`)

Provides MySQL connection management using GORM:

- **Environment-based initialization**: `InitDbFromEnv()` reads MySQL config from environment variables:
  - `MYSQL_USERNAME` (required)
  - `MYSQL_PASSWORD` (required)
  - `MYSQL_HOST` (required)
  - `MYSQL_PORT` (optional, defaults to 3306)
  - `MYSQL_DATABASE` (required)
- **Connection pooling**: Configured with 10 max idle, 100 max open connections, 1-hour max lifetime
- **GORM logger**: Set to Info level by default

### Logger Package (`logger/`)

Built on uber-go/zap with file rotation:

- **Initialization**: Call `InitWithConfig(Config{AppName: "yourapp"})` before use
- **Log location**: `~/logs/{AppName}/main.log` (falls back to `{cwd}/logs/{AppName}/main.log`)
- **Log rotation**: Daily rotation, keeps last 3 days
- **Debug mode**: Set via `DEBUG` environment variable (accepts "1", "true", "on", "enable")
  - Enables debug level logging
  - Outputs to both console and file
- **Default log level**: Info (unless debug mode enabled)
- **Log format**: Console format by default (not JSON), uses ISO8601 timestamps
- **Available functions**: `Debugf/Debugw`, `Infof/Infow`, `Warnf/Warnw`, `Errorf/Errorw`
  - `*f` variants use printf-style formatting
  - `*w` variants use structured key-value pairs

### Utils Package (`utils/`)

HTTP helper utilities:

- `DumpHttpRequest()`: Debug function to log complete HTTP request details
- `GetHttpProto()`: Determines HTTP scheme from `X-Forwarded-Proto` header (useful behind proxies)

## Common Patterns

### Logger Initialization

```go
import "github.com/chzealot/gobase/logger"

logger.InitWithConfig(logger.Config{
    AppName: "myapp",
})

logger.Infof("Application started")
logger.Debugw("Debug info", zap.String("key", "value"))
```

### Database Initialization

```go
import "github.com/chzealot/gobase/database"

db := database.NewDatabase()
if err := db.InitDbFromEnv(); err != nil {
    // Handle error
}
// Use db.DB for GORM operations
```

### DingTalk Client Usage

```go
import "github.com/chzealot/gobase/dingtalk"

client := dingtalk.NewDingTalkClient(clientId, clientSecret)
token, err := client.GetAccessToken()
// Use token for API calls
```

## Model Definitions

DingTalk models are extensively defined in `dingtalk/models/`:
- `auth.go`: Authentication request/response types
- `user.go`: User information structures
- `calendar.go`: Calendar metadata
- `event.go`: Calendar event structures with full detail (attendees, recurrence, reminders, etc.)
- `todo.go`: Todo task creation structures
- `api.go`: Common API response wrappers

All models use JSON tags for serialization.
