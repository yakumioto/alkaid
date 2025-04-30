# Alkaid

Alkaid 是一个基于 Gin 框架实现的轻量级 RESTful API 服务框架，专注于简单、灵活和高效的 API 开发。

## 主要特性

- **基于 Gin**: 底层使用 Gin Web 框架，保持高性能和灵活性
- **API 版本控制**: 支持通过 URL 前缀和请求头进行 API 版本管理
- **服务组织**: 简单直观的服务和控制器注册机制
- **认证机制**: 可自定义的认证系统接口，支持多种认证策略
- **中间件支持**: 支持自定义中间件处理流程
- **响应渲染**: 根据请求 Accept 头自动选择合适的响应格式（JSON、XML、YAML等）

## 快速开始

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yakumioto/alkaid"
)

func main() {
	// 创建 Alkaid 实例
	alk := alkaid.NewAlkaid(
		alkaid.WithEngineMode(gin.DebugMode),
		alkaid.WithEnableUrlVersionPrefix(),
		alkaid.WithVersionHeader("X-Version"),
	)

	// 注册服务和控制器
	alk.RegisterService(&alkaid.Service{
		Name: "users",
		Controllers: []alkaid.Controller{
			&UserController{},
		},
	})

	// 启动服务器
	alk.Run(":8080")
}
```

## 社区交流

QQ：567880968

Telegram：<https://t.me/fab_alkaid>

