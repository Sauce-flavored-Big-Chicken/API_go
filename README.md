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
