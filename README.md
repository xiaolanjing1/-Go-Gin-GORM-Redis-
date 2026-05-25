# Go User System

一个基于 Go + Gin + Gorm + MySQL + Redis + Docker 的用户系统。

## Tech Stack

- Go
- Gin
- Gorm
- MySQL
- Redis
- Docker

## Features

- 用户注册
- 用户登录
- MySQL 数据持久化
- Docker 容器化部署

## Run With Docker

### Build Image

```bash
docker build -t user_system .
```

### Run Container

```bash
docker run -p 8080:8080 user_system
```

## Project Structure

```text
.
├── config
├── router
├── model
├── api
├── Dockerfile
├── go.mod
└── main.go
```

## Database

请提前准备 MySQL 数据库：

```sql
CREATE DATABASE mygorm;
```

## Author

学习 Go 云原生与 Docker 工程化实践项目。