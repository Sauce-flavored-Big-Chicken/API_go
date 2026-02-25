# 数字社区 API 接口说明文档

# 项目说明

## 1. 服务器地址

> **待补充**：文档未给出具体服务器地址（host:port）

部分接口需要先登录获取授权 TOKEN 才能调用，具体接口在每个章节中说明。

## 2. 注销说明

| 项目 | 说明 |
|------|------|
| 接口地址 | `/logout` |
| 请求方法 | POST |
| 请求参数 | 无 |

## 3. 安全认证

需要认证的接口需在请求头设置：

| 请求头参数 | 参数值 |
|------------|--------|
| Authorization | 登录获取的 TOKEN（可拼接 Bearer 前缀） |

> 注：支持两种格式 `TOKEN` 或 `Bearer TOKEN`

## 4. 系统默认用户

| 账号 | 密码 |
|------|------|
| test01 | 123456 |

> 注：以下接口示例数据均基于此账号。

## 5. 表格分页参数

| 参数名 | 说明 | 类型 |
|--------|------|------|
| pageNum | 当前页码 | Integer |
| pageSize | 每页数据条数 | Integer |

## 6. 返回状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 正常 |
| 500 | 系统异常 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 未找到资源 |

## 7. 业务辅助字段（可忽略）

| 字段名 | 含义 |
|--------|------|
| searchValue | 搜索内容 |
| createBy | 创建用户 |
| updateBy | 更新用户 |
| updateTime | 更新时间 |
| remark | 备注 |
| Params | 参数集合 |

## 8. 资源URL规范

| 类型 | 说明 | 示例 |
|------|------|------|
| 绝对URL | 完整网址 | `http://your-domain.com/dev-api/profile/upload/xxx.jpg` |
| 相对路径 | 需要拼接基础地址 | `/dev-api/profile/upload/xxx.jpg`、`/profile/upload/xxx.jpg` |

> 注：文档中示例混用了两种形式，实际开发时需与后端确认统一的资源地址策略。

# 通用接口

## 1. 登录注册

### 1.1 手机登录

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/phone/login` |
| 请求方法 | POST |
| 请求类型 | application/json |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 |
|--------|------|------|------|
| phone | 手机号码 | true | string |
| smsCode | 验证码 | true | string |

**请求示例**
```json
{
  "phone": "13411551115",
  "smsCode": "6767"
}
```

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码，200成功，其他失败 | number |
| msg | 返回消息 | string |
| token | 返回token信息 | string |

**响应示例**
```json
{
  "code": 200,
  "msg": "操作成功",
  "token": "38022e59b215a6735021f9452fbbafe8"
}
```

### 1.2 用户登录

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/login` |
| 请求方法 | POST |
| 请求类型 | application/json |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 |
|--------|------|------|------|
| userName | 用户名 | true | string |
| password | 用户密码 | true | string |

**请求示例**
```json
{
  "userName": "test01",
  "password": "123456"
}
```

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码，200成功，其他失败 | number |
| msg | 返回消息 | string |
| token | 返回token信息 | string |

**响应示例**
```json
{
  "code": 200,
  "msg": "操作成功",
  "token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiI0MGMxNjMzZTI4MGFmMTRiMWMwY2NiNzc2ZjJhYThmZSIsImF1ZCI6IiIsImlhdCI6MTc3MTgyOTk1MSwibmJmIjoxNzcxODI5OTUxLCJleHAiOjE3NzE4MzcxNTEsImRhdGEiOnsidWlkIjoidGVzdDAxIn19.F6tqaY7Uh8k_MfRGdiJ4lU_7EG-mPHDtfJs-GYbSB8w"
}
```

### 1.3 获取验证码

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/smsCode` |
| 请求方法 | GET |
| 请求类型 | application/x-www-form-urlencoded |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| phone | 手机号码 | true | string | query |

**请求示例**
```
/prod-api/api/smsCode?phone=13411551115
```

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码，200成功，其他失败 | number |
| msg | 返回消息 | string |
| data | 验证码 | string |

**响应示例**
```json
{
  "code": 200,
  "msg": "请求成功",
  "data": "6184"
}
```

### 1.4 用户注册

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/register` |
| 请求方法 | POST |
| 请求类型 | application/json |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 |
|--------|------|------|------|
| avatar | 头像 | false | string |
| userName | 用户名 | true | string |
| nickName | 昵称 | false | string |
| password | 密码 | true | string |
| phoneNumber | 电话号码 | true | string |
| sex | 性别（0男1女） | true | string |
| email | 邮箱 | false | string |
| idCard | 身份证 | false | string |
| address | 住址 | false | string |
| introduction | 个人简介 | false | string |

**请求示例**
```json
{
  "avatar": "27e7fd58-0972-4dbf-941c-590624e6a886.png",
  "userName": "David",
  "nickName": "大卫",
  "password": "123456",
  "phoneNumber": "15840669812",
  "sex": "0",
  "email": "David@163.com",
  "idCard": "210113199808242137",
  "address": "XXX",
  "introduction": "XXXX"
}
```

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码，200成功，其他失败 | number |
| msg | 返回消息 | string |

**响应示例**
```json
{
  "code": 200,
  "msg": "操作成功"
}
```

## 2. 用户信息

### 2.1 查询个人基本信息

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/user/getUserInfo` |
| 请求方法 | GET |
| 请求类型 | application/x-www-form-urlencoded |
| 认证 | 需要Token |

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码，200成功，其他失败 | number |
| msg | 返回消息 | string |
| data | 用户信息对象 | object |
| ├userId | 用户ID | number |
| ├userName | 用户名 | string |
| ├nickName | 用户昵称 | string |
| ├avatar | 用户头像 | string |
| ├phoneNumber | 手机号 | string |
| ├email | 邮箱 | string |
| ├sex | 用户性别（0男1女） | string |
| ├idCard | 身份证号 | string |
| ├address | 住址 | string |
| ├introduction | 个人简介 | string |
| ├balance | 账户余额 | number |
| ├score | 用户积分 | number |

**响应示例**
```json
{
  "code": 200,
  "msg": "获取数据成功",
  "data": {
    "userId": 1,
    "userName": "David",
    "avatar": "/dev-api/profile/upload/image/client-img-2.jpg",
    "nickName": "大卫",
    "phoneNumber": "15898125461",
    "sex": "1",
    "email": "lixl@163.com",
    "idCard": "210882199807251656",
    "address": "XXX",
    "introduction": "XXXX",
    "balance": 0,
    "score": 0
  }
}
```

### 2.2 修改个人基本信息

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/user/updateUserInfo` |
| 请求方法 | PUT |
| 请求类型 | application/json |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 |
|--------|------|------|------|
| nickName | 用户昵称 | true | string |
| avatar | 用户头像 | false | string |
| email | 邮箱 | false | string |
| idCard | 身份证号 | false | string |
| phoneNumber | 手机号 | false | string |
| sex | 用户性别（0男1女） | false | string |
| address | 住址 | false | string |
| introduction | 个人简介 | false | string |

**请求示例**
```json
{
  "email": "lixl@163.com",
  "idCard": "210882199807251656",
  "nickName": "大卫王",
  "phoneNumber": "15898125461",
  "sex": "0",
  "address": "XXX",
  "introduction": "XXXX"
}
```

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码，200成功，其他失败 | number |
| msg | 返回消息 | string |

**响应示例**
```json
{
  "code": 200,
  "msg": "操作成功"
}
```

### 2.3 修改用户密码

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/user/resetPwd` |
| 请求方法 | PUT |
| 请求类型 | application/json |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 |
|--------|------|------|------|
| oldPassword | 用户旧密码 | true | string |
| newPassword | 用户新密码 | true | string |

**请求示例**
```json
{
  "oldPassword": "123456",
  "newPassword": "123789"
}
```

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码，200成功，其他失败 | number |
| msg | 返回消息 | string |

**响应示例**
```json
{
  "code": 200,
  "msg": "操作成功"
}
```

### 2.4 用户列表

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/user/list` |
| 请求方法 | GET |
| 请求类型 | application/x-www-form-urlencoded |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| pageNum | 页码 | false | int | query |
| pageSize | 每页条数 | false | int | query |

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码 | number |
| msg | 返回消息 | string |
| data | 用户列表 | array |
| total | 总数 | number |

**响应示例**
```json
{
  "code": 200,
  "msg": "请求成功",
  "data": [
    {
      "userId": 1,
      "userName": "test01",
      "nickName": "测试用户",
      "avatar": "/dev-api/profile/upload/image/avatar.jpg",
      "phoneNumber": "13800138000",
      "email": "test@example.com",
      "sex": "0",
      "status": "0"
    }
  ],
  "total": 1
}
```

## 3. 广告轮播

### 3.1 查询引导页及主页轮播

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/rotation/list` |
| 请求方法 | GET |
| 请求类型 | application/x-www-form-urlencoded |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| type | 广告类型（1引导页轮播，2主页轮播） | true | string | query |
| pageNum | 页码 | false | string | query |
| pageSize | 每页数量 | false | string | query |

**请求示例**
```
/prod-api/api/rotation/list?pageNum=1&pageSize=8&type=2
```

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码，200成功，其他失败 | number |
| msg | 消息内容 | string |
| data | 广告轮播列表 | array |
| ├id | 广告ID | number |
| ├advTitle | 广告标题 | string |
| ├advImg | 广告图片 | string |
| ├type | 广告类型 | string |
| total | 总记录数 | number |

**响应示例**
```json
{
  "code": 200,
  "msg": "获取数据成功",
  "data": [
    {
      "advTitle": "低代码技术",
      "advImg": "/dev-api/profile/upload/image/484d77ef-b2ce-4861-9d32-727f6be3079e.png",
      "type": "2"
    },
    {
      "advTitle": "全球创客路演",
      "advImg": "/dev-api/profile/upload/image/b9d9f081-8a76-41dc-8199-23bcb3a64fcc.png",
      "type": "2"
    }
  ],
  "total": 2
}
```

## 4. 新闻资讯

### 4.1 获取新闻分类

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/press/category/list` |
| 请求方法 | GET |
| 请求类型 | application/x-www-form-urlencoded |
| 认证 | 需要Token |

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码，200成功，其他失败 | string |
| msg | 消息内容 | string |
| data | 新闻分类列表 | array |
| ├id | 分类编号 | number |
| ├name | 分类名称 | string |
| ├sort | 分类序号 | number |
| ├appType | app类型 | string |
| total | 总记录数 | number |

**响应示例**
```json
{
  "code": 200,
  "msg": "查询成功",
  "data": [
    {
      "id": 9,
      "appType": "smart_city",
      "name": "今日要闻"
    }
  ],
  "total": "1"
}
```

### 4.2 获取所有新闻列表

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/press/newsList` |
| 请求方法 | GET |
| 请求类型 | application/x-www-form-urlencoded |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| pageNum | 页码 | false | string | query |
| pageSize | 每页数量 | false | string | query |

**请求示例**
```
/prod-api/api/press/newsList?pageNum=1&pageSize=8
```

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码 | string |
| msg | 消息内容 | string |
| data | 新闻列表 | array |
| ├id | 新闻ID | number |
| ├title | 新闻标题 | string |
| ├subTitle | 新闻副标题 | string |
| ├content | 新闻内容 | string |
| ├cover | 新闻封面图片地址 | string |
| ├publishDate | 发布日期 | string(date-time) |
| ├readNum | 阅读数 | number |
| ├likeNum | 点赞数 | number |
| ├commentNum | 评论数 | number |
| ├status | 状态 | string |
| ├tags | 标签 | string |
| ├type | 类型 | string |
| ├top | 是否推荐（Y/N） | string |
| ├hot | 是否热点（Y/N） | string |
| total | 总记录数 | number |

**响应示例**
```json
{
  "code": 200,
  "msg": "查询成功",
  "data": [
    {
      "id": 5,
      "cover": "/dev-api/profile/upload/image/2021/04/01/c1eb74b2-e964-4388-830a-1b606fc9699f.png",
      "title": "测试新闻标题",
      "subTitle": "测试新闻子标题",
      "content": "<p>内容</p>",
      "status": "Y",
      "publishDate": "2021-04-01",
      "tags": null,
      "commentNum": 1,
      "likeNum": 2,
      "readNum": 10,
      "type": "2",
      "top": "Y",
      "hot": "N"
    }
  ],
  "total": 1
}
```

### 4.3 获取分类新闻列表

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/press/category/newsList` |
| 请求方法 | GET |
| 请求类型 | application/x-www-form-urlencoded |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| pageNum | 页码 | false | string | query |
| pageSize | 每页数量 | false | string | query |
| id | 新闻分类id | true | string | query |

**请求示例**
```
/prod-api/api/press/category/newsList?pageNum=1&pageSize=8&id=2
```

**响应参数**

（同4.2获取所有新闻列表）

### 4.4 获取新闻详细信息

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/press/news/{id}` |
| 请求方法 | GET |
| 请求类型 | application/x-www-form-urlencoded |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| id | 新闻id | true | number | path |

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码 | string |
| msg | 消息内容 | string |
| data | 新闻实体 | object |
| ├id | 新闻ID | number |
| ├appType | app类型 | string |
| ├title | 新闻标题 | string |
| ├subTitle | 新闻副标题 | string |
| ├content | 新闻内容 | string |
| ├cover | 新闻封面图片地址 | string |
| ├publishDate | 发布日期 | string(date-time) |
| ├readNum | 阅读数 | number |
| ├likeNum | 点赞数 | number |
| ├commentNum | 评论数 | number |
| ├status | 状态 | string |
| ├tags | 标签 | string |
| ├type | 类型 | string |
| ├top | 是否推荐（Y/N） | string |
| ├hot | 是否热点（Y/N） | string |

**响应示例**
```json
{
  "code": 200,
  "data": {
    "id": 5,
    "appType": "movie",
    "cover": "/dev-api/profile/upload/image/2021/04/01/c1eb74b2-e964-388-830a-1b606fc9699f.png",
    "title": "驱蚊器无去",
    "subTitle": "123123123",
    "content": "<p>企鹅王请问</p>",
    "status": "Y",
    "publishDate": "2021-04-01",
    "tags": null,
    "commentNum": null,
    "likeNum": 3,
    "readNum": null,
    "type": "2",
    "top": "Y",
    "hot": "N"
  },
  "msg": "操作成功"
}
```

### 4.5 新闻点赞

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/press/like/{id}` |
| 请求方法 | PUT |
| 请求类型 | application/json |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| id | 新闻id | true | number | path |

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码，200成功，其他失败 | number |
| msg | 返回消息 | string |

**响应示例**
```json
{
  "code": 200,
  "msg": "操作成功"
}
```

### 4.6 发表新闻评论

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/comment/pressComment` |
| 请求方法 | POST |
| 请求类型 | application/json |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| content | 评论内容 | true | string | body |
| newsId | 新闻id | true | string | body |
| userName | 用户名 | true | string | body |

**请求示例**
```json
{
  "content": "漫步在长江国家文化公园九江城区段，可以感受到长江的自然风光和文化底蕴交融碰撞，让人忍不住去探寻那些藏在历史背后的辉煌与灿烂。",
  "newsId": "1",
  "userName": "test01"
}
```

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码，200成功，其他失败 | number |
| msg | 返回消息 | string |

**响应示例**
```json
{
  "code": 200,
  "msg": "操作成功"
}
```

### 4.7 获取新闻评论列表

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/comment/comment/{id}` |
| 请求方法 | GET |
| 请求类型 | application/x-www-form-urlencoded |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| id | 新闻ID | true | number | path |

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码 | string |
| msg | 消息内容 | string |
| data | 评论列表 | array |
| ├id | 评论ID | number |
| ├content | 评论内容 | string |
| ├commentDate | 评论时间 | string |
| ├userId | 评论人ID | number |
| ├newsId | 新闻ID | number |
| ├likeNum | 点赞数 | number |
| total | 总记录数 | number |

**响应示例**
```json
{
  "code": 200,
  "msg": "获取数据成功",
  "data": [
    {
      "id": 2,
      "content": "漫步在长江国家文化公园九江城区段，可以感受到长江的自然风光和文化底蕴交融碰撞，让人忍不住去探寻那些藏在历史背后的辉煌与灿烂。",
      "commentDate": "2023-10-13 15:46:33",
      "userId": 2,
      "newsId": 1,
      "likeNum": 0
    }
  ],
  "total": 1
}
```

### 4.8 评论点赞

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/comment/like/{id}` |
| 请求方法 | PUT |
| 请求类型 | application/json |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| id | 评论id | true | number | path |

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码，200成功，其他失败 | number |
| msg | 返回消息 | string |

**响应示例**
```json
{
  "code": 200,
  "msg": "操作成功"
}
```

## 5. 文件上传

### 5.1 通用上传接口

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/common/upload` |
| 请求方法 | POST |
| 请求类型 | multipart/form-data |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| file | 上传的文件对象 | true | text | formData |

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码，200成功，其他失败 | number |
| msg | 返回消息 | string |
| fileName | 文件名 | string |
| url | 文件访问地址 | string |

**响应示例**
```json
{
  "code": 200,
  "fileName": "test.txt",
  "url": "/profile/upload/file/test.txt",
  "msg": "操作成功"
}
```

## 6. 公告通知

### 6.1 通知列表

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/notice/list` |
| 请求方法 | GET |
| 请求类型 | application/x-www-form-urlencoded |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| pageNum | 页码 | false | string | query |
| pageSize | 每页数量 | false | string | query |
| noticeStatus | 通知状态（1已读，0未读） | true | string | query |

**请求示例**
```
/prod-api/api/notice/list?pageNum=1&pageSize=8&noticeStatus=1
```

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码，200成功，其他失败 | string |
| msg | 消息内容 | string |
| data | 通知列表 | array |
| ├id | 通知ID | number |
| ├noticeTitle | 标题 | string |
| ├noticeStatus | 状态（1已读，0未读） | string |
| ├contentNotice | 内容 | string |
| ├releaseUnit | 发布单位 | string |
| ├phone | 手机 | string |
| ├createTime | 时间 | string |
| ├expressId | 通知类型id | string |
| ├noticeName | 通知类型 | string |
| total | 总记录数 | number |

**响应示例**
```json
{
  "code": 200,
  "msg": "请求成功",
  "data": [
    {
      "id": 1,
      "noticeTitle": "缴费通知",
      "noticeStatus": 1,
      "contentNotice": "现 XXXX 物业服务中心自 20xx 年 X 月 X 日起开始收取 20xx年 7 月的物业管理费及代收水、电费。",
      "releaseUnit": "物业管理",
      "phone": "13999999999",
      "createTime": "2023-04-22 15:18:25",
      "expressId": 1,
      "noticeName": "重要通知"
    }
  ],
  "total": 2
}
```

### 6.2 通知列表详情

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/notice/{id}` |
| 请求方法 | GET |
| 请求类型 | application/x-www-form-urlencoded |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| id | 通知ID | true | number | path |

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码 | string |
| msg | 消息内容 | string |
| data | 通知详情 | object |
| ├id | 通知ID | number |
| ├noticeTitle | 标题 | string |
| ├noticeStatus | 状态（1已读，0未读） | string |
| ├contentNotice | 内容 | string |
| ├releaseUnit | 发布单位 | string |
| ├phone | 手机 | string |
| ├createTime | 时间 | string |
| ├expressId | 通知类型id | string |
| ├noticeName | 通知类型 | string |

**响应示例**
```json
{
  "code": 200,
  "msg": "获取数据成功",
  "data": {
    "id": 1,
    "noticeTitle": "缴费通知",
    "noticeStatus": 1,
    "contentNotice": "现 XXXX 物业服务中心自 20xx 年 X 月 X 日起开始收取 20xx 年 7 月的物业管理费及代收水、电费。",
    "releaseUnit": "物业管理",
    "phone": "13999999999",
    "createTime": "2023-04-22 15:18:25",
    "expressId": 1,
    "noticeName": "重要通知"
  }
}
```

### 6.3 社区通知变已读

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/readNotice/{id}` |
| 请求方法 | PUT |
| 请求类型 | application/json |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| id | 通知id | true | number | path |

**请求示例**
```
/prod-api/api/readNotice/1
```

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码，200成功，其他失败 | number |
| msg | 返回消息 | string |

**响应示例**
```json
{
  "code": 200,
  "msg": "操作成功"
}
```

# 数字社区

## 1. 友邻帖子

### 1.1 友邻帖子列表

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/friendly_neighborhood/list` |
| 请求方法 | GET |
| 请求类型 | application/x-www-form-urlencoded |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| pageNum | 页码 | false | string | query |
| pageSize | 每页数量 | false | string | query |

**请求示例**
```
/prod-api/api/friendly_neighborhood/list?pageNum=1&pageSize=8
```

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码 | string |
| msg | 消息内容 | string |
| data | 帖子列表 | array |
| ├id | 帖子ID | number |
| ├publishName | 发布人 | string |
| ├likeNum | 点赞数 | number |
| ├title | 标题 | string |
| ├publishTime | 发布时间 | string |
| ├publishContent | 发布内容 | string |
| ├commentNum | 评论数 | number |
| ├imgUrl | 图片 | string |
| ├userImgUrl | 用户头像图片 | string |
| total | 总记录数 | number |

**响应示例**
```json
{
  "code": 200,
  "msg": "请求成功",
  "data": [
    {
      "id": 1,
      "publishName": "小黄",
      "likeNum": 2,
      "title": "今天天气真好啊",
      "publishTime": "2023-04-22 15:01:11",
      "publishContent": "夏天的风我永远记得，清清楚楚的说要热死我……",
      "commentNum": 15,
      "imgUrl": "/dev-api/profile/upload/image/f23f9d02-ae1e-4065-9730-42df2e539e20.jpg",
      "userImgUrl": "/dev-api/profile/upload/image/10e57ee1-c19e-407d-bc0e-f13abab96bcb.jpg"
    },
    {
      "id": 2,
      "publishName": "刘詩涵",
      "likeNum": 2,
      "title": "今天天气真差呀",
      "publishTime": "2023-04-22 15:01:11",
      "publishContent": "自打这入夏以来啊，太阳就偏偏独宠我一人，于是啊我就劝太阳一定要雨露均沾，这太阳非是不听哪，就晒我！就晒我！就晒我！把我晒黢黑黢黑呀！",
      "commentNum": 2,
      "imgUrl": "/dev-api/profile/upload/image/98ee36c3-c2bc-4608-bf13-c9b3281f3cc4.jpg",
      "userImgUrl": "/dev-api/profile/upload/image/10e57ee1-c19e-407d-bc0e-f13abab96bcb.jpg"
    }
  ],
  "total": 2
}
```

### 1.2 友邻帖子发布评论

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/friendly_neighborhood/add/comment` |
| 请求方法 | POST |
| 请求类型 | application/json |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| content | 回帖内容 | true | string | body |
| neighborhoodId | 友邻帖子ID | true | int | body |

**请求示例**
```json
{
  "content": "在中国式现代化新征程上策马扬鞭",
  "neighborhoodId": 1
}
```

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码，200成功，其他失败 | number |
| msg | 返回消息 | string |

**响应示例**
```json
{
  "code": 200,
  "msg": "操作成功"
}
```

### 1.3 友邻帖子详情

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/friendly_neighborhood/{id}` |
| 请求方法 | GET |
| 请求类型 | application/json |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| id | 友邻帖子id | true | number | path |

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码 | string |
| msg | 消息内容 | string |
| data | 帖子详情 | object |
| ├id | 帖子ID | number |
| ├publishName | 发布人 | string |
| ├likeNum | 点赞数 | number |
| ├title | 标题 | string |
| ├publishTime | 发布时间 | string |
| ├publishContent | 发布内容 | string |
| ├commentNum | 评论数 | number |
| ├imgUrl | 图片 | string |
| ├userImgUrl | 用户头像图片 | string |
| ├userComment | 评论列表 | array |
| ├├id | 评论id | string |
| ├├userName | 用户名 | string |
| ├├userId | 用户ID | string |
| ├├avatar | 用户头像 | string |
| ├├content | 评论内容 | string |
| ├├likeNum | 点赞数 | string |
| ├├publishTime | 评论时间 | string |
| ├├neighborhoodId | 友邻帖子id | string |

**响应示例**
```json
{
  "code": 200,
  "msg": "请求成功",
  "data": {
    "id": 1,
    "publishName": "小黄",
    "likeNum": 2,
    "title": "今天天气真好啊",
    "publishTime": "2023-04-22 15:01:11",
    "publishContent": "夏天的风我永远记得，清清楚楚的说要热死我……",
    "commentNum": 15,
    "imgUrl": "/dev-api/profile/upload/image/f23f9d02-ae1e-4065-9730-42df2e539e20.jpg",
    "userImgUrl": "/dev-api/profile/upload/image/10e57ee1-c19e-407d-bc0e-f13abab96bcb.jpg",
    "userComment": [
      {
        "id": 1,
        "userName": "test01",
        "userId": 1,
        "avatar": "/dev-api/profile/upload/image/27e7fd58-0972-4dbf-941c-590624e6a886.png",
        "content": "在中国式现代化新征程上策马扬鞭",
        "likeNum": 12,
        "publishTime": "2026-02-15 01:44:32",
        "neighborhoodId": 1
      }
    ]
  }
}
```

## 2. 社区活动

### 2.1 推荐活动列表

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/activity/topList` |
| 请求方法 | GET |
| 请求类型 | application/x-www-form-urlencoded |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| pageNum | 页码 | false | string | query |
| pageSize | 每页数量 | false | string | query |

**请求示例**
```
/prod-api/api/activity/topList?pageNum=1&pageSize=8
```

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码 | string |
| msg | 消息内容 | string |
| data | 活动列表 | array |
| ├id | 活动ID | number |
| ├category | 分类（1文化，2体育，3公益，4亲子） | string |
| ├title | 标题 | string |
| ├picPath | 图片地址 | string |
| ├startDate | 活动开始时间 | string |
| ├endDate | 活动结束时间 | string |
| ├sponsor | 主办方 | string |
| ├content | 详情 | string |
| ├signUpNum | 已报名人数 | number |
| ├maxNum | 最大人数 | number |
| ├signUpEndDate | 报名截止时间 | string |
| ├isTop | 是否推荐（1推荐） | string |
| total | 总记录数 | number |

**响应示例**
```json
{
  "code": 200,
  "msg": "请求成功",
  "data": [
    {
      "id": 15,
      "category": "2",
      "title": "植绿护绿",
      "picPath": "/dev-api/profile/upload/image/665014142627b085436522.jpg",
      "startDate": "2028/5/15 11:11",
      "endDate": "2028/5/25 11:11",
      "sponsor": "常德市西洞庭管理区农业农村局志愿服务队",
      "content": "西洞庭农业农村局开展植绿护绿志愿活动",
      "signUpNum": 6,
      "maxNum": 100,
      "signUpEndDate": null,
      "isTop": "1"
    }
  ],
  "total": 1
}
```

### 2.2 社区活动列表

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/activity/list` |
| 请求方法 | GET |
| 请求类型 | application/x-www-form-urlencoded |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| pageNum | 页码 | false | string | query |
| pageSize | 每页数量 | false | string | query |

**请求示例**
```
/prod-api/api/activity/list?pageNum=1&pageSize=8
```

**响应参数**

（同2.1推荐活动列表）

### 2.3 搜索社区活动

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/activity/search` |
| 请求方法 | POST |
| 请求类型 | application/x-www-form-urlencoded |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| words | 搜索关键字 | true | string | query |

**请求示例**
```
/prod-api/api/activity/search?words=活动
```

**响应参数**

（同2.1推荐活动列表）

### 2.4 社区活动详情

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/activity/{id}` |
| 请求方法 | GET |
| 请求类型 | application/json |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| id | 社区活动id | true | number | path |

**响应参数**

（同2.1推荐活动列表）

### 2.5 社区分类活动列表

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/activity/category/list/:id` |
| 请求方法 | GET |
| 请求类型 | application/x-www-form-urlencoded |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| id | 分类id | true | int | path |
| pageNum | 页码 | false | int | query |
| pageSize | 每页数量 | false | int | query |

**请求示例**
```
/prod-api/api/activity/category/list/1?pageNum=1&pageSize=8
```

**响应参数**

（同2.1推荐活动列表）

### 2.6 活动报名

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/registration` |
| 请求方法 | POST |
| 请求类型 | application/json |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 |
|--------|------|------|------|
| activityId | 活动id | true | int |

**请求示例**
```json
{
  "activityId": 2
}
```

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码，200成功，其他失败 | number |
| msg | 返回消息 | string |

**响应示例**
```json
{
  "code": 200,
  "msg": "操作成功"
}
```

### 2.7 活动签到

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/checkin/{id}` |
| 请求方法 | PUT |
| 请求类型 | application/json |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| id | 活动id | true | int | path |

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码，200成功，其他失败 | number |
| msg | 返回消息 | string |

**响应示例**
```json
{
  "code": 200,
  "msg": "操作成功"
}
```

### 2.8 活动评论

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/registration/comment/{id}` |
| 请求方法 | PUT |
| 请求类型 | application/json |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| id | 活动id | true | int | path |
| evaluate | 学习感言 | true | string | body |
| star | 评分 | true | int | body |

**请求示例**
```json
{
  "evaluate": "sssssssssssssssss",
  "star": 4
}
```

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码，200成功，其他失败 | number |
| msg | 返回消息 | string |

**响应示例**
```json
{
  "code": 200,
  "msg": "操作成功"
}
```

# 绿动未来

## 1. 数据卡片

### 1.1 获取数据卡片列表

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/common/datacard` |
| 请求方法 | GET |
| 请求类型 | application/x-www-form-urlencoded |

**请求参数**  
无

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码，200成功，其他失败 | number |
| msg | 返回消息 | string |
| data | 数据列表 | array |

**响应示例**
```json
{
  "code": 200,
  "msg": "操作成功",
  "data": [
    {
      "id": 1,
      "icon": "/static/image/icon/news_hot.png",
      "title": "AQI 指数",
      "num": "45",
      "unit": "优",
      "trend": "/static/image/icon/down_arrow.png",
      "sort": 1
    }
  ]
}
```

### 1.2 新增数据卡片

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/common/datacard` |
| 请求方法 | POST |
| 请求类型 | application/json |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 |
|--------|------|------|------|
| icon | 图标URL | false | string |
| title | 标题 | true | string |
| num | 数值 | true | string |
| unit | 单位 | true | string |
| trend | 趋势图标URL | false | string |
| sort | 排序 | false | int |

**请求示例**
```json
{
  "icon": "/static/image/icon/news_hot.png",
  "title": "AQI 指数",
  "num": "45",
  "unit": "优",
  "trend": "/static/image/icon/down_arrow.png",
  "sort": 1
}
```

**响应示例**
```json
{
  "code": 200,
  "msg": "创建成功"
}
```

### 1.3 更新数据卡片

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/common/datacard/{id}` |
| 请求方法 | PUT |
| 请求类型 | application/json |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 |
|--------|------|------|------|
| id | 卡片ID | true | int |
| icon | 图标URL | false | string |
| title | 标题 | false | string |
| num | 数值 | false | string |
| unit | 单位 | false | string |
| trend | 趋势图标URL | false | string |
| sort | 排序 | false | int |

**响应示例**
```json
{
  "code": 200,
  "msg": "更新成功"
}
```

### 1.4 删除数据卡片

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/common/datacard/{id}` |
| 请求方法 | DELETE |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 |
|--------|------|------|------|
| id | 卡片ID | true | int |

**响应示例**
```json
{
  "code": 200,
  "msg": "删除成功"
}
```

## 2. 数据系列

### 2.1 获取数据系列列表

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/data/list` |
| 请求方法 | GET |
| 请求类型 | application/x-www-form-urlencoded |
| 认证 | 需要Token |

**请求参数**  
无

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码 | number |
| msg | 返回消息 | string |
| data | 数据列表 | array |
| total | 总数 | number |

**响应示例**
```json
{
  "code": 200,
  "msg": "请求成功",
  "data": [
    {
      "id": 1,
      "list_key": "list_1",
      "name": "list_1",
      "data": "[{\"name\":\"aqi\",\"data\":[65,82,95]}]",
      "sort": 1
    }
  ],
  "total": 1
}
```

### 2.2 新增数据系列

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/data/list` |
| 请求方法 | POST |
| 请求类型 | application/json |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 |
|--------|------|------|------|
| data | JSON数据数组字符串 | true | string |

**请求示例**
```json
{
  "data": "[{\"name\":\"aqi\",\"data\":[65,82,95,78,105,88,70]},{\"name\":\"pm2.5\",\"data\":[22,31,38,27,45,33,24]}]"
}
```

**响应示例**
```json
{
  "code": 200,
  "msg": "创建成功",
  "data": {
    "listKey": "list_1"
  }
}
```

> 注：系统会自动生成下一个 list_N 的 Key（N 为现有最大数字+1），如 list_1、list_2 等。

### 2.3 更新数据系列

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/data/list/{id}` |
| 请求方法 | PUT |
| 请求类型 | application/json |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 |
|--------|------|------|------|
| id | 数据ID | true | int |
| data | JSON数据数组字符串 | true | string |

**响应示例**
```json
{
  "code": 200,
  "msg": "更新成功"
}
```

### 2.4 删除数据系列

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/data/list/{id}` |
| 请求方法 | DELETE |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 |
|--------|------|------|------|
| id | 数据ID | true | int |

**响应示例**
```json
{
  "code": 200,
  "msg": "删除成功"
}
```

### 2.5 按Key查询数据

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/data/{listKey}` |
| 请求方法 | GET |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 |
|--------|------|------|------|
| listKey | 列表Key，如 list_1 | true | string |

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码 | number |
| msg | 返回消息 | string |
| data | 数据数组 | array |

**响应示例**
```json
{
  "code": 200,
  "msg": "请求成功",
  "data": [
    {
      "name": "aqi",
      "data": [65, 82, 95, 78, 105, 88, 70]
    },
    {
      "name": "pm2.5",
      "data": [22, 31, 38, 27, 45, 33, 24]
    }
  ]
}
```

## 3. 题库管理

### 3.1 获取题库列表

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/question/list` |
| 请求方法 | GET |
| 请求类型 | application/x-www-form-urlencoded |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 |
|--------|------|------|------|
| pageNum | 页码 | false | int |
| pageSize | 每页条数 | false | int |

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码 | number |
| msg | 返回消息 | string |
| data | 题目列表 | array |
| total | 总数 | number |

**响应示例**
```json
{
  "code": 200,
  "msg": "请求成功",
  "data": [
    {
      "id": 1,
      "questionType": "1",
      "level": "1",
      "question": "以下哪项属于可再生能源？",
      "optionA": "煤炭",
      "optionB": "太阳能",
      "optionC": "石油",
      "optionD": "天然气",
      "answer": "B",
      "score": 2,
      "status": "0"
    }
  ],
  "total": 10
}
```

### 3.2 新增题目

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/question` |
| 请求方法 | POST |
| 请求类型 | application/json |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 |
|--------|------|------|------|
| questionType | 题型(1=选择题,4=判断题) | true | string |
| level | 难度(1=简单,2=中等,3=困难) | true | string |
| question | 题目内容 | true | string |
| optionA | 选项A | false | string |
| optionB | 选项B | false | string |
| optionC | 选项C | false | string |
| optionD | 选项D | false | string |
| optionE | 选项E | false | string |
| optionF | 选项F | false | string |
| answer | 答案 | true | string |
| score | 分值 | false | int |
| status | 状态(0=启用,1=禁用) | false | string |

**响应示例**
```json
{
  "code": 200,
  "msg": "创建成功"
}
```

### 3.3 更新题目

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/question/{id}` |
| 请求方法 | PUT |
| 请求类型 | application/json |
| 认证 | 需要Token |

**请求参数**  
同新增题目

**响应示例**
```json
{
  "code": 200,
  "msg": "更新成功"
}
```

### 3.4 删除题目

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/question/{id}` |
| 请求方法 | DELETE |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 |
|--------|------|------|------|
| id | 题目ID | true | int |

**响应示例**
```json
{
  "code": 200,
  "msg": "删除成功"
}
```

### 3.5 随机抽题

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/question/questionList/{questionType}/{level}` |
| 请求方法 | GET |
| 请求类型 | application/x-www-form-urlencoded |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 |
|--------|------|------|------|
| questionType | 题型(1=选择题,4=判断题) | true | string |
| level | 难度(1=简单,2=中等,3=困难) | true | string |

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码 | number |
| msg | 返回消息 | string |
| data | 题目列表 | array |
| total | 总数 | number |

**响应示例**
```json
{
  "code": 200,
  "msg": "请求成功",
  "data": [
    {
      "id": 5,
      "questionType": "1",
      "question": "日常减碳更推荐的出行方式是？",
      "optionA": "开车",
      "optionB": "骑自行车",
      "optionC": "打车",
      "answer": "B",
      "score": 2
    }
  ],
  "total": 1
}
```

### 3.6 提交答卷

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/question/savePaper` |
| 请求方法 | POST |
| 请求类型 | application/json |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 |
|--------|------|------|------|
| score | 分数 | true | number/string |
| answer | 答案JSON数组 | true | array |

**请求示例**
```json
{
  "score": 80,
  "answer": [
    {"qid": 101, "answer": "B"},
    {"qid": 102, "answer": "A"}
  ]
}
```

**响应示例**
```json
{
  "code": 200,
  "msg": "操作成功"
}
```

# 图片管理

## 1. 图片列表

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/common/images` |
| 请求方法 | GET |
| 请求类型 | application/x-www-form-urlencoded |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 |
|--------|------|------|------|
| pageNum | 页码 | false | int |
| pageSize | 每页条数 | false | int |

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码 | number |
| msg | 返回消息 | string |
| data | 图片列表 | array |
| total | 总数 | number |

**响应示例**
```json
{
  "code": 200,
  "msg": "获取成功",
  "data": [
    {
      "name": "image.jpg",
      "url": "/profile/upload/image/20260101120000_test.jpg",
      "thumbUrl": "/profile/upload/thumb/image/20260101120000_test.jpg",
      "size": 123456,
      "created": "2026-01-01 12:00:00"
    }
  ],
  "total": 1
}
```

> 注：返回的 thumbUrl 为缩略图路径（140px内最长边，28KB以内），可优先使用提升加载速度。

## 2. 删除图片

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/common/images` |
| 请求方法 | DELETE |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 |
|--------|------|------|------|
| url | 图片URL路径 | true | string |

**响应示例**
```json
{
  "code": 200,
  "msg": "删除成功"
}
```

# 文件管理

## 1. 文件列表

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/common/files` |
| 请求方法 | GET |
| 请求类型 | application/x-www-form-urlencoded |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 |
|--------|------|------|------|
| pageNum | 页码 | false | int |
| pageSize | 每页条数 | false | int |

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码 | number |
| msg | 返回消息 | string |
| data | 文件列表 | array |
| total | 总数 | number |

**响应示例**
```json
{
  "code": 200,
  "msg": "获取成功",
  "data": [
    {
      "name": "document.pdf",
      "url": "/profile/upload/file/20260101120000_document.pdf",
      "size": 1024000,
      "created": "2026-01-01 12:00:00"
    }
  ],
  "total": 1
}
```

## 2. 删除文件

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/common/files` |
| 请求方法 | DELETE |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 |
|--------|------|------|------|
| url | 文件URL路径 | true | string |

**响应示例**
```json
{
  "code": 200,
  "msg": "删除成功"
}
```

## 2. 报名记录

### 2.1 报名记录列表

| 项目 | 说明 |
|------|------|
| 接口地址 | `/prod-api/api/registration/list` |
| 请求方法 | GET |
| 请求类型 | application/x-www-form-urlencoded |
| 认证 | 需要Token |

**请求参数**

| 参数名 | 说明 | 必须 | 类型 | 请求类型 |
|--------|------|------|------|----------|
| pageNum | 页码 | false | int | query |
| pageSize | 每页条数 | false | int | query |
| activityId | 活动ID筛选 | false | string | query |
| userId | 用户ID筛选 | false | string | query |

**响应参数**

| 参数名 | 说明 | 类型 |
|--------|------|------|
| code | 状态码 | number |
| msg | 返回消息 | string |
| data | 报名列表 | array |
| total | 总数 | number |

**响应示例**
```json
{
  "code": 200,
  "msg": "请求成功",
  "data": [
    {
      "id": 1,
      "userId": 1,
      "nickName": "张三",
      "activityId": 1,
      "status": "0",
      "checkinStatus": "0",
      "comment": "活动很有意义",
      "star": 5,
      "createTime": "2026-01-01 12:00:00"
    }
  ],
  "total": 1
}
```

