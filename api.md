# Doc
**Base URL**: `http://localhost:8080`

**认证方式**: JWT Token (Bearer Authentication)

## 目录
1. [认证接口](#认证接口)
2. [用户管理接口](#用户管理接口)
3. [农产品追溯接口](#农产品追溯接口)
4. [区块链代理接口](#区块链代理接口)
5. [文件上传接口](#文件上传接口)
6. [错误码说明](#错误码说明)

---

## 认证接口

### 用户登录
**POST** `/login`
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "user_name": "admin",
    "password": "agri_chain"
  }'
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Njg4MzgwMTgsInVzZXJuYW1lIjoiYWRtaW4ifQ.eHJalLIr1HOwdQ_1edoty2s0vqxqXLZlXZ6wDTgSZKQ"}

用户登录获取JWT Token

#### 请求参数
```json
{
  "user_name": "string",    // 用户名
  "password": "string"      // 密码
}
```

#### 响应示例
**成功响应 (200)**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**失败响应 (401)**:
```json
{
  "error": "user do not exist"
}
```
```json
{
  "error": "password not correct"
}
```

#### 默认账户
- 用户名: `admin`
- 密码: `agri_chain`

---

## 用户管理接口

### 获取用户档案
**GET** `/profile`

获取当前登录用户的详细信息

#### 请求头
```
Authorization: {JWT_TOKEN}
```

#### 响应示例
**成功响应 (200)**:
```json
{
  "status": "ok",
  "data": {
    "id": 1,
    "user_name": "admin",
    "nick_name": "农业管理员",
    "email": "admin@agri-chain.com",
    "invite_code": "",
    "agree_terms": false,
    "avatar_url": "/public/avatars/profile.jpg",
    "profile_public": true,
    "email_notifications": false,
    "last_login_at": null,
    "created_at": "2026-01-18T22:47:31.484875+08:00",
    "updated_at": "2026-01-18T22:47:31.484875+08:00"
  }
}
```

### 更新用户档案
**POST** `/profile`

更新当前用户的档案信息

#### 请求头
```
Authorization: {JWT_TOKEN}
Content-Type: application/json
```

#### 请求参数
```json
{
  "nick_name": "string",           // 可选，昵称
  "email": "string",              // 可选，邮箱
  "profile_public": boolean,      // 可选，是否公开档案
  "email_notifications": boolean // 可选，是否接收邮件通知
}
```

#### 响应示例
**成功响应 (200)**:
```json
{
  "status": "ok",
  "msg": "Profile updated successfully"
}
```

### 上传用户头像
**POST** `/upload/avatar`

上传用户头像图片

#### 请求头
```
Authorization: {JWT_TOKEN}
Content-Type: multipart/form-data
```

#### 请求参数
- `avatar`: 图片文件 (支持image/*类型，最大2MB)

#### 响应示例
**成功响应 (200)**:
```json
{
  "status": "ok",
  "msg": "Avatar uploaded successfully",
  "data": {
    "avatar_url": "/public/avatars/admin_1768834104.jpg"
  }
}
```

### 访问主菜单
**GET** `/menu`

验证用户登录状态并返回欢迎信息

#### 请求头
```
Authorization: {JWT_TOKEN}
```

#### 响应示例
**成功响应 (200)**:
```json
{
  "status": "ok",
  "msg": "欢迎登陆农产追溯通客户端"
}
```

---

## 农产品追溯接口

### 上传农产品数据
**POST** `/upload`

将农产品追溯数据上传到区块链网络

#### 请求头
```
Authorization: {JWT_TOKEN}
Content-Type: application/json
```

#### 请求参数
```json
{
  "eggplant_id": 1001,                    // 农产品ID
  "product_height": 100,                  // 生产环节高度
  "product_hash": "a1b2c3d4...",          // 生产环节哈希
  "transport_height": 101,                // 运输环节高度
  "transport_hash": "b2c3d4e5...",        // 运输环节哈希
  "process_height": 102,                  // 加工环节高度
  "process_hash": "c3d4e5f6...",          // 加工环节哈希
  "storage_height": 103,                  // 存储环节高度
  "storage_hash": "d4e5f678...",          // 存储环节哈希
  "sell_height": 104,                     // 销售环节高度
  "sell_hash": "e5f67890..."              // 销售环节哈希
}
```

#### 响应示例
**成功响应 (200)**:
```json
{
  "status": "ok",
  "msg": "上传成功"
}
```

**失败响应 (503)**:
```json
{
  "status": "error",
  "msg": "所有区块链节点连接失败"
}
```

### 查询农产品追溯信息
**GET** `/message`

根据农产品ID查询其在区块链上的追溯信息

#### 请求头
```
Authorization: {JWT_TOKEN}
```

#### URL参数
- `id`: 农产品ID
- `node`: 区块链节点地址 (可选)

#### 请求示例
```
GET /message?id=1001&node=blockchain_node1:8081
```

#### 响应示例
**成功响应 (200)**:
```json
{
  "status": "ok",
  "msg": "{农产品追溯数据JSON字符串}"
}
```

**未找到 (400)**:
```json
{
  "error": "no such eggplant"
}
```

### 获取区块链节点列表
**GET** `/nodes`

获取系统中注册的区块链节点列表

#### 请求头
```
Authorization: {JWT_TOKEN}
```

#### 响应示例
**成功响应 (200)**:
```json
[
  {
    "id": 1411106969,
    "addr": "0.0.0.0:8081",
    "pub_key": "A9FLe1gP/41TAqlq/6xPET+lA+9S8Lmfow4xseW0RL6t",
    "create_time": 1768747651686333063,
    "verify_time": 0
  }
]
```

---

## 区块链代理接口

### 获取区块链节点信息
**GET** `/blockchain/nodes`

curl -X GET http://localhost:8080/blockchain/nodes \
  -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Njg4MzgwMTgsInVzZXJuYW1lIjoiYWRtaW4ifQ.eHJalLIr1HOwdQ_1edoty2s0vqxqXLZlXZ6wDTgSZKQ" \
  -H "Accept: application/json"
{"data":[{"addr":"1d8e9dc010e7:8081","create_time":1768748398753541380,"http_addr":"http://1d8e9dc010e7:8080","id":9817828935,"pub_key":"AprYI3RkrxEVxfVki0hLa5+GnqN7+WrmvO57zHiI79Hc","verify_time":0},{"addr":"8bb8411e82bf:8081","create_time":1768748398819238068,"http_addr":"http://8bb8411e82bf:8080","id":3399646069,"pub_key":"AuXlQKjuMwV9R0uE50VXTTPUp8zLhDTj1e1fyOtez+dR","verify_time":0},{"addr":"4b6ad63b9238:8081","create_time":1768748398900239595,"http_addr":"http://4b6ad63b9238:8080","id":9031847448,"pub_key":"Aq+w5nbR8QGf9hPHzdgTjJjIS7/fL511dhtvy48mKhYa","verify_time":0},{"addr":"0610f2270c1f:8081","create_time":1768748398955672352,"http_addr":"http://0610f2270c1f:8080","id":1399741211,"pub_key":"A67rnx9ywS7MC+18cgGAm5y4iso2ocRncqVp8EhrrnXU","verify_time":0}],"status":"ok"}

curl -X GET http://localhost:8080/blockchain/nodes \
-H "Accept: application/json"

获取可用的区块链节点列表及其HTTP访问地址

#### 请求头
```
Authorization: {JWT_TOKEN}
```

#### 响应示例
**成功响应 (200)**:
```json
{
  "status": "ok",
  "data": [
    {
      "id": 1411106969,
      "addr": "0.0.0.0:8081",
      "pub_key": "A9FLe1gP/41TAqlq/6xPET+lA+9S8Lmfow4xseW0RL6t",
      "create_time": 1768747651686333063,
      "verify_time": 0,
      "http_addr": "http://0.0.0.0:8080"
    }
  ]
}
```

### 获取区块链高度
**GET** `/blockchain/height`

获取区块链的最新高度信息
curl -X GET http://localhost:8080/blockchain/height \
  -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Njg4MzgwMTgsInVzZXJuYW1lIjoiYWRtaW4ifQ.eHJalLIr1HOwdQ_1edoty2s0vqxqXLZlXZ6wDTgSZKQ" \
  -H "Accept: application/json"

curl -X GET http://localhost:8080/blockchain/height \
  -H "Accept: application/json"

{"height":6,"success":true}

#### 请求头
```
Authorization: {JWT_TOKEN}
```

#### 响应示例
**成功响应 (200)**:
```json
{
  "height": 150,
  "success": true
}
```

**失败响应 (500)**:
```json
{
  "status": "error",
  "msg": "未找到可用的区块链节点"
}
```

### 获取区块范围数据
**GET** `/blockchain/blocks`

curl -X GET http://localhost:8080/blockchain/blocks \
  -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Njg4MzgwMTgsInVzZXJuYW1lIjoiYWRtaW4ifQ.eHJalLIr1HOwdQ_1edoty2s0vqxqXLZlXZ6wDTgSZKQ" \
  -H "Accept: application/json"

curl -X GET http://localhost:8080/blockchain/blocks \
  -H "Accept: application/json"

{"blocks":[{"dataHash":"30303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030","height":0,"leader":9817828935,"nonce":0,"prevHash":"30303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030","scores":{},"timestamp":0,"version":0},{"dataHash":"65336230633434323938666331633134396166626634633839393666623932343237616534316534363439623933346361343935393931623738353262383535","height":1,"leader":3399646069,"nonce":3511354979366533000,"prevHash":"32663634626362616339373166303835363937326162353033343635383632383766393739353335643666643136386566643339316230323439333338303864","scores":{"9817828935":0},"timestamp":1768748578,"version":0},{"dataHash":"65336230633434323938666331633134396166626634633839393666623932343237616534316534363439623933346361343935393931623738353262383535","height":2,"leader":9817828935,"nonce":2025425436650788400,"prevHash":"38353037383536313838613235303662613933653162623566323666613331376162313139663663383337323838393162613961656334623430626266313062","scores":{"3399646069":0},"timestamp":1768748758,"version":0},{"dataHash":"65336230633434323938666331633134396166626634633839393666623932343237616534316534363439623933346361343935393931623738353262383535","height":3,"leader":3399646069,"nonce":4288239474819126000,"prevHash":"39656430366266333266626330353639633335303733633038373562623035393765386237346564636565316139643335323366376165346536343266383835","scores":{"9817828935":0},"timestamp":1768751179,"version":0},{"dataHash":"65336230633434323938666331633134396166626634633839393666623932343237616534316534363439623933346361343935393931623738353262383535","height":4,"leader":9817828935,"nonce":7786840486826589000,"prevHash":"31383337616134396562343462393638626463663237346134613235333061656337663932306564383537303035366564653463653936636332366333303364","scores":{"3399646069":0},"timestamp":1768751359,"version":0},{"dataHash":"65336230633434323938666331633134396166626634633839393666623932343237616534316534363439623933346361343935393931623738353262383535","height":5,"leader":3399646069,"nonce":1777681661085323000,"prevHash":"34316234313562383065663038643838366465323438346666346338646434656564623630626162633831633230393266363231656232396438663437363065","scores":{"9817828935":0},"timestamp":1768751571,"version":0},{"dataHash":"65336230633434323938666331633134396166626634633839393666623932343237616534316534363439623933346361343935393931623738353262383535","height":6,"leader":9817828935,"nonce":141357664291397820,"prevHash":"38626666326232393839663361646264643932623330373937366436353538616537386264353839383365363134396237386263313532633163386635383163","scores":{"3399646069":0},"timestamp":1768751751,"version":0},{"dataHash":"65336230633434323938666331633134396166626634633839393666623932343237616534316534363439623933346361343935393931623738353262383535","height":7,"leader":3399646069,"nonce":4677321512122869000,"prevHash":"62366338373535646133613330346139393237363133383136303932393735343239313839306434636263383065316135383062336466313531346531643338","scores":{"9817828935":0},"timestamp":1768751931,"version":0}],"count":8,"success":true}

根据指定范围获取区块数据

#### 返回值含义解释
{
    "dataHash": "653...855",     // 区块数据哈希
    "height": 7,                 // 区块高度
    "leader": 3399646069,        // 出块Leader节点ID
    "nonce": 4677321512122869000, // 随机数/工作量证明
    "prevHash": "623...d38",     // 前一个区块的哈希
    "scores": {"9817828935": 0}, // 信用评分系统
    "timestamp": 1768751931,     // 区块时间戳
    "version": 0                 // 区块版本号
  }


#### 请求头
```
Authorization: {JWT_TOKEN}
```



#### URL参数
0索引开始，左闭右闭
- `start`: 起始区块高度
- `end`: 结束区块高度

#### 请求示例
```
GET /blockchain/blocks?start=0&end=5
```

curl -X GET "http://localhost:8080/blockchain/blocks?start=0&end=5" \
  -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Njg4MzgwMTgsInVzZXJuYW1lIjoiYWRtaW4ifQ.eHJalLIr1HOwdQ_1edoty2s0vqxqXLZlXZ6wDTgSZKQ" \
  -H "Accept: application/json"

  curl -X GET "http://localhost:8080/blockchain/blocks?start=0&end=5" \
>   -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Njg4MzgwMTgsInVzZXJuYW1lIjoiYWRtaW4ifQ.eHJalLIr1HOwdQ_1edoty2s0vqxqXLZlXZ6wDTgSZKQ" \
>   -H "Accept: application/json"
{"blocks":[{"dataHash":"30303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030","height":0,"leader":9817828935,"nonce":0,"prevHash":"30303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030","scores":{},"timestamp":0,"version":0},{"dataHash":"65336230633434323938666331633134396166626634633839393666623932343237616534316534363439623933346361343935393931623738353262383535","height":1,"leader":3399646069,"nonce":3511354979366533000,"prevHash":"32663634626362616339373166303835363937326162353033343635383632383766393739353335643666643136386566643339316230323439333338303864","scores":{"9817828935":0},"timestamp":1768748578,"version":0},{"dataHash":"65336230633434323938666331633134396166626634633839393666623932343237616534316534363439623933346361343935393931623738353262383535","height":2,"leader":9817828935,"nonce":2025425436650788400,"prevHash":"38353037383536313838613235303662613933653162623566323666613331376162313139663663383337323838393162613961656334623430626266313062","scores":{"3399646069":0},"timestamp":1768748758,"version":0},{"dataHash":"65336230633434323938666331633134396166626634633839393666623932343237616534316534363439623933346361343935393931623738353262383535","height":3,"leader":3399646069,"nonce":4288239474819126000,"prevHash":"39656430366266333266626330353639633335303733633038373562623035393765386237346564636565316139643335323366376165346536343266383835","scores":{"9817828935":0},"timestamp":1768751179,"version":0},{"dataHash":"65336230633434323938666331633134396166626634633839393666623932343237616534316534363439623933346361343935393931623738353262383535","height":4,"leader":9817828935,"nonce":7786840486826589000,"prevHash":"31383337616134396562343462393638626463663237346134613235333061656337663932306564383537303035366564653463653936636332366333303364","scores":{"3399646069":0},"timestamp":1768751359,"version":0},{"dataHash":"65336230633434323938666331633134396166626634633839393666623932343237616534316534363439623933346361343935393931623738353262383535","height":5,"leader":3399646069,"nonce":1777681661085323000,"prevHash":"34316234313562383065663038643838366465323438346666346338646434656564623630626162633831633230393266363231656232396438663437363065","scores":{"9817828935":0},"timestamp":1768751571,"version":0}],"count":6,"success":true}

#### 响应示例
**成功响应 (200)**:
```json
{
  "blocks": [
    {
      "height": 0,
      "leader": "node1",
      "scores": {...},
      "timestamp": 1768834200,
      "dataHash": "abc123...",
      "prevHash": "def456...",
      "nonce": 12345,
      "version": "1.0"
    }
  ],
  "count": 1,
  "success": true
}
```

### 获取区块链节点状态
**GET** `/blockchain/status`

获取区块链节点的运行状态信息

#### 请求头
```
Authorization: {JWT_TOKEN}
```

#### 响应示例
**成功响应 (200)**:
```json
{
  "nodeId": 1411106969,
  "address": "0.0.0.0:8081",
  "isLeader": true,
  "chainHeight": 150,
  "sequenceId": 5,
  "poolSize": 10,
  "poolCap": 100,
  "success": true
}
```

---

## 文件上传接口

### 静态文件访问
**GET** `/public/{filepath}`

访问系统中的静态文件，如上传的头像

#### 请求示例
```
GET /public/avatars/admin_1768834104.jpg
```

#### 响应
返回对应的静态文件内容

---

## 系统接口

### 区块链响应处理
**POST** `/meta_data`

接收来自区块链节点的农产品数据响应 (系统内部使用)

#### 请求参数
```json
{
  "eggplant_id": 1001,
  "product_height": 100,
  "product_hash": "a1b2c3d4...",
  "transport_height": 101,
  "transport_hash": "b2c3d4e5...",
  "process_height": 102,
  "process_hash": "c3d4e5f6...",
  "storage_height": 103,
  "storage_hash": "d4e5f678...",
  "sell_height": 104,
  "sell_hash": "e5f67890...",
  "first_seen": "2026-01-18T22:40:00Z"
}
```

---

## 错误码说明

### HTTP状态码
- **200 OK**: 请求成功
- **400 Bad Request**: 请求参数错误
- **401 Unauthorized**: 未授权或Token无效
- **404 Not Found**: 资源不存在
- **500 Internal Server Error**: 服务器内部错误
- **503 Service Unavailable**: 服务不可用

### 错误响应格式
```json
{
  "status": "error",
  "msg": "错误描述信息"
}
```

或

```json
{
  "error": "错误描述信息"
}
```

### 常见错误信息
- `请完成登陆`: 需要提供有效的JWT Token
- `user do not exist`: 用户不存在
- `password not correct`: 密码错误
- `Invalid token`: Token格式错误或已过期
- `Missing authorization token`: 缺少认证Token
- `所有区块链节点连接失败`: 无法连接到区块链网络
- `未找到可用的区块链节点`: 系统中没有配置区块链节点

---

## 请求示例

### 使用curl进行API调用

1. **登录获取Token**:
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"user_name": "admin", "password": "agri_chain"}'
```

2. **获取用户档案**:
```bash
curl -X GET http://localhost:8080/profile \
  -H "Authorization: YOUR_JWT_TOKEN"
```

3. **上传农产品数据**:
```bash
curl -X POST http://localhost:8080/upload \
  -H "Authorization: YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "eggplant_id": 1001,
    "product_height": 100,
    "product_hash": "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456",
    "transport_height": 101,
    "transport_hash": "b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef1234567a",
    "process_height": 102,
    "process_hash": "c3d4e5f6789012345678901234567890abcdef1234567890abcdef1234567ab2",
    "storage_height": 103,
    "storage_hash": "d4e5f6789012345678901234567890abcdef1234567890abcdef1234567ab2c3",
    "sell_height": 104,
    "sell_hash": "e5f6789012345678901234567890abcdef1234567890abcdef1234567ab2c3d4"
  }'
```

4. **查询农产品信息**:
```bash
curl -X GET "http://localhost:8080/message?id=1001&node=blockchain_node1:8081" \
  -H "Authorization: YOUR_JWT_TOKEN"
```
