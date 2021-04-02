# nebula-http-gateway

Gateway to provide a http interface for the Nebula Graph service.

## Build

```bash
$ cd /path/to/nebula-http-gateway
$ make
```

## Run

```bash
$ ./nebula-httpd
```

## Required

- Go 11+
- [beego](https://beego.me/)

## 使用说明

#### API 定义
| Name | Path | Method |
| --- | --- | --- |
| connect | /api/db/connect | POST |
| exec | /api/db/exec | POST |
| disconnect | /api/db/disconnect | POST |

#### Connect接口
- 接口请求

请求的 json body 为
```json
{
  "username": "user",
  "password": "password",
  "address": "192.168.8.26",
  "port": 32172
}
```
```bash
$ curl -i  -X POST -d '{"username":"user","password":"password","address":"192.168.8.26","port":9669}' http://127.0.0.1:8080/api/db/connect
HTTP/1.1 200 OK
Content-Length: 100
Content-Type: application/json; charset=utf-8
Server: beegoServer:1.12.3
Set-Cookie: nsid=bec2e665ba62a13554b617d70de8b9b9; Path=/; HttpOnly
Set-Cookie: Secure=true; Path=/
Set-Cookie: SameSite=None; Path=/
Date: Fri, 02 Apr 2021 08:49:18 GMT

{
  "code": 0,
  "data": "5e18fa40-5343-422f-84e3-e7f9cad6b735",
  "message": "Login successfully"
}
```

- 实现细节
```go
// controllers/db.go
// 根据请求 body 里的 graph address、port，认证的 username, password 去建立连接会话
func (this *DatabaseController) Connect() {
	...
	nsid, err := dao.Connect(params.Address, params.Port, params.Username, params.Password)
	if err == nil {
		res.Code = 0
		m := make(map[string]common.Any)
		m["nsid"] = nsid
		res.Data = nsid 
	    // 在 response header 里设置 cookie 
	    // nsid 会被 HMAC-SH256 加密算法改写
	    // 下次再请求的时候 http header 里会自动带入这些信息
		this.Ctx.SetCookie("Secure","true")  
		this.Ctx.SetCookie("SameSite","None")
		this.SetSession("nsid", nsid)
	...
}

// service/pool/pool.go
func NewConnection(address string, port int, username string, password string) (nsid string, err error) {
    // 初始化连接池
	pool, err := nebula.NewConnectionPool(hostList, poolConfig, nebulaLog)
	if err != nil {
		return "", errors.New("Fail to initialize the connection pool")
	}
    // 检查目标 host 是否可以连接
	err = pool.Ping(hostList[0], 5000*time.Millisecond)
	if err != nil {
		return "", err
	}

	...

    // 根据认证信息获取 session，这里会使用 uuid v4 的格式生成 nsid
    // connectionPool 是 map 类型，作为会话连接池保存 nsid 到 Connection 的映射关系
	session, err := pool.GetSession(username, password)
	if err == nil {
		nsid = uuid.NewV4().String()
		connectionPool[nsid] = &Connection{
			RequestChannel: make(chan ChannelRequest),
			CloseChannel:   make(chan bool),
			updateTime:     time.Now().Unix(),
			session:        session,
			account: &Account{
				username: username,
				password: password,
			},
		}
        // 当前用户并发连接计数
		currentConnectionNum++

        // 开启 goroutine 并发执行查询
		go func() {
			connection := connectionPool[nsid]
			for {
				select {
				case request := <-connection.RequestChannel:
                    // 从 RequestChannel 获取查询语句，执行后再将结果写入 ResponseChannel
					response, err := connection.session.Execute(request.Gql)
					request.ResponseChannel <- ChannelResponse{
						Result: response,
						Error:  err,
					}
				case <-connection.CloseChannel:
                    // 从 CloseChannel 收到 close 信号，释放 session，同时从
                    // 资源池 connectionPool 移除 connection，由于 map 是非线程安全的，
                    // 因此在递减 currentConnectionNum 时需要加锁
					connection.session.Release()
					connectLock.Lock()
					delete(connectionPool, nsid)
					currentConnectionNum--
					connectLock.Unlock()
					// Exit loop
					return
				}
			}
		}()
		return nsid, err
	}
}
```

#### Exec接口
- 请求接口
```json
{
"gql": "show spaces;"
}
```
```bash
$ curl -H "Cookie: SameSite=None; nsid=56aa2463c850cfc5e2199e943cf4807b" -H "nsid: 56aa2463c850cfc5e2199e943cf4807b" -X POST -d '{"gql": "show spaces;"}' http://127.0.0.1:8080/api/db/exec
  {
    "code": 0,
    "data": {
      "headers": [
        "Name"
      ],
      "tables": [
        {
          "Name": "nba"
        }
      ],
      "timeCost": 4232
    },
    "message": ""
```

- 实现细节
```go
// controllers/db.go
func Execute(nsid string, gql string) (result ExecuteResult, err error) {
	result = ExecuteResult{
		Headers: make([]string, 0),
		Tables:  make([]map[string]common.Any, 0),
	}
    // 根据 nsid 从回话连接池里获取 connection
	connection, err := pool.GetConnection(nsid)
	if err != nil {
		return result, err
	}

	responseChannel := make(chan pool.ChannelResponse)
    // 向 RequestChannel 中写入包含 gql 的 request
	connection.RequestChannel <- pool.ChannelRequest{
		Gql:             gql,
		ResponseChannel: responseChannel,
	}
    // 等待查询结果返回
	response := <-responseChannel
    // 后面为查询结果类型解析
	resp := response.Result
    ...
```

#### Disconnect接口
- 接口请求
```bash
$ curl -X POST http://127.0.0.1:8080/api/db/disconnect
{
  "code": 0,
  "data": null,
  "message": "Disconnect successfully"
```

- 实现细节
```go
// controllers/db.go
func (this *DatabaseController) Disconnect() {
	var res Response
	nsid := this.GetSession("nsid")
	if nsid != nil {
        // 根据 nsid 从连接池释放连接
		dao.Disconnect(nsid.(string))
	}
	res.Code = 0
	res.Message = "Disconnect successfully"
	this.Data["json"] = &res
	this.ServeJSON()
}

// serivce/pool/pool.go
func Disconnect(nsid string) {
	connection := connectionPool[nsid]
	if connection != nil {
        // 调用 nebula.Session 的 Release 方法释放登录后产生的会话及连接信息
		connection.session.Release()
        // 从 connectionPool 删除 connection
		delete(connectionPool, nsid)
	}
	return
}
```
