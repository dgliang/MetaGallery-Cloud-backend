# 项目名称 MetaGallery-Cloud-backend

_A cloud disk web backend with a file gallery based on IPFS_

基于 Gin + GORM 和 IPFS 的 Go 文件库云盘网页后端

## 简介

该项目是一个使用 Go 语言开发的后端应用，采用了 Gin Web 框架来处理 HTTP 请求，并通过 GORM 与数据库进行交互。该项目的目的是构建一个快速、简单的 RESTful API 后端服务。

## 主要技术栈

- [Go](https://golang.org/)：项目的主要编程语言
- [Gin](https://gin-gonic.com/)：轻量级、高性能的 Go Web 框架
- [GORM](https://gorm.io/)：强大且灵活的 ORM 库，用于处理数据库操作
- [MySQL](https://www.mysql.com/cn/)：数据库选择 MySQL

## 功能

- 用户注册与登录
- JWT 身份认证
- 文件的基本的 CRUD 操作（Create、Read、Update、Delete）
- 文件画廊（Files Gallery）

## 项目结构

```
├── config         # 配置文件
├── controllers    # 处理 HTTP 请求的控制器
├── dao            # 数据访问对象与数据库进行交互
├── models         # 数据库模型
├── routes         # 路由定义
├── middlewares    # 中间件
├── services       # 业务逻辑
├── resources      # 静态资源
├── util           # 工具函数
├── main.go        # 项目入口文件
├── .env           # 环境变量配置文件
└── go.mod         # 依赖管理文件
```

## 环境要求

- Go 1.18 及以上版本
- 数据库 MySQL

## 安装步骤

1. 克隆仓库

   ```bash
   git clone https://github.com/dgliang/MetaGallery-Cloud-backend.git
   cd MetaGallery-Cloud-backend
   ```

2. 安装依赖

   ```bash
   go mod tidy
   ```

3. 配置 `.env` 文件

   在项目根目录下创建 `.env` 文件，并添加以下内容：

   ```env
   # 数据库配置
   DB_HOST=<your-database-host>
   DB_PORT=<your-database-port>
   DB_USER=<your-database-user>
   DB_PASSWORD=<your-database-password>
   DB_NAME=<your-database-name>

   # 后端服务器配置
   HOST_URL=<backend-server-url>

   # JWT 配置
   JWT_SECRET=<your_jwt_secret>

   # 文件和缓存文件存储配置
   FILE_DIR_PATH="./resources/files/"
   LOCAL_CACHE_PATH="./resources/.cache"

   # Pinata 和 IPFS 配置
   PINATA_JWT=
   PINATA_HOST_URL=
   PINATA_GATEWAY_KEY=

   # SSL 证书路径
   SSL_CRT_PATH="./config/keys/xxx.crt"
   SSL_KEY_PATH="./config/keys/xxx.key"

   ```

4. 启动项目

   运行以下命令启动项目：

   ```bash
   go run main.go
   ```

   项目将运行在 `http://localhost:8080`

## API 文档

项目的 API 文档可以通过以下方式查看：

https://apifox.com/apidoc/shared-3065fd37-8877-4aa1-bc69-7d5e2aa013e8
