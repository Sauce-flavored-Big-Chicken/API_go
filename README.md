# 数字社区 API 后端

基于 Go + Gin + Gorm + SQLite 实现的数字社区后端服务，严格按照接口文档实现。

## 技术栈

- **Go 1.21+**
- **Gin** - Web 框架
- **Gorm** - ORM
- **SQLite** - 数据库

## 项目结构

```
.
├── cmd/
│   └── main.go           # 入口
├── internal/
│   ├── config/           # 配置和数据库
│   ├── handlers/        # API 处理器
│   ├── middleware/      # 鉴权、CORS、日志
│   ├── models/          # 数据模型
│   └── router/         # 路由
├── admin.html           # 数据管理后台
├── go.mod / go.sum     # 依赖
└── server              # 编译后的二进制
```

## 快速启动

```bash
# 编译
go build -o server ./cmd/main.go

# 运行（默认端口 8080）
./server
```

## 测试账号

- 用户名: `test01`
- 密码: `123456`

## API 列表

| 模块 | 接口 | 说明 |
|------|------|------|
| 认证 | POST /prod-api/api/login | 用户登录 |
| 认证 | POST /prod-api/api/phone/login | 手机登录 |
| 认证 | GET /prod-api/api/SMSCode | 获取短信验证码 |
| 认证 | POST /prod-api/api/register | 注册 |
| 用户 | GET /prod-api/api/user/getUserInfo | 获取用户信息 |
| 用户 | PUT /prod-api/api/user/updateUserInfo | 更新用户信息 |
| 用户 | PUT /prod-api/api/user/resetPwd | 重置密码 |
| 新闻 | GET /prod-api/api/press/newsList | 新闻列表 |
| 新闻 | GET /prod-api/api/press/news/{id} | 新闻详情 |
| 公告 | GET /prod-api/api/notice/list | 公告列表 |
| 活动 | GET /prod-api/api/activity/List | 活动列表 |
| 活动 | GET /prod-api/api/activity/topList | 热门活动 |
| 活动 | POST /prod-api/api/activity/search | 搜索活动 |
| 友邻圈 | GET /prod-api/api/friendly_neighborhood/list | 友邻圈列表 |
| 上传 | POST /prod-api/api/common/upload | 文件上传 |

## 管理后台

启动服务后，浏览器打开 `admin.html` 即可管理数据。

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| SERVER_PORT | 8080 | 服务端口 |
| DB_PATH | ./data.db | 数据库路径 |
| JWT_SECRET | xxx | JWT 密钥 |

## 许可证

MIT

## 部署指南

### 本地部署

```bash
# 1. 克隆项目
git clone https://github.com/Sauce-flavored-Big-Chicken/API_go.git
cd API_go

# 2. 安装依赖
go mod tidy

# 3. 编译
go build -o server ./cmd/main.go

# 4. 运行
./server
```

### Docker 部署

项目已内置单容器配置：

- `Dockerfile`：一个容器同时运行 Go API + Nginx 前端
- `docker-compose.yml`：一键启动

一键启动：

```bash
docker compose up --build -d
```

访问地址：

- 前端：`http://localhost:3000`
- API（经同域反代）：`http://localhost:3000/prod-api/...`

停止并清理：

```bash
docker compose down
```

若需要清理数据库和上传文件卷：

```bash
docker compose down -v
```

### GitHub Actions 自动构建并推送 Docker Hub

已内置工作流：`.github/workflows/docker-publish.yml`

触发条件：

- push 到 `main` / `master`
- push `v*` 标签（如 `v1.0.0`）
- 手动触发 `workflow_dispatch`

请在 GitHub 仓库配置：

- `Secrets`
  - `DOCKERHUB_USERNAME`: Docker Hub 用户名
  - `DOCKERHUB_TOKEN`: Docker Hub Access Token
- `Variables`（可选）
  - `DOCKERHUB_REPOSITORY`: Docker Hub 仓库名（不填则默认用 `api-go`）

推送后镜像标签包含：

- 分支标签（如 `main`）
- 版本标签（如 `v1.0.0`）
- `sha-<commit>`
- 默认分支额外推送 `latest`

### 生产环境部署

```bash
# 使用 Systemd 服务
sudo nano /etc/systemd/system/digital-community.service
```

```ini
[Unit]
Description=Digital Community API
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/digital-community
ExecStart=/opt/digital-community/server
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
# 启用服务
sudo systemctl daemon-reload
sudo systemctl enable digital-community
sudo systemctl start digital-community
```

### Nginx 反向代理配置

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        root /var/www/digital-community;
        index admin.html;
        try_files $uri $uri/ /admin.html;
    }

    location /prod-api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /profile/ {
        proxy_pass http://127.0.0.1:8080;
    }
}
```

### 生产环境建议

1. **安全**
   - 修改默认 JWT_SECRET
   - 使用 HTTPS
   - 配置防火墙规则

2. **数据库**
   - 定期备份 `data.db`
   - 生产环境可迁移到 MySQL/PostgreSQL（修改 gorm.io/driver）

3. **监控**
   - 配置日志轮转
   - 使用 PM2 或 Supervisor 管理进程

4. **性能**
   - 根据需要调整 SQLite 配置
   - 考虑使用缓存
