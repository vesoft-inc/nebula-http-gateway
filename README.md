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

- Go 1.13+
- [beego](https://beego.me/)

## Version
| Nebula Graph version | Http-gateway tag | 
|----------------------|---------------------------|
| 1.x                  | v1.0                      |
| 2.0.x & 2.5.x        | v2.0                      |
| 2.6.x                | v2.1.x                    |
| 3.0.x                | 2.2.x                     |
| 3.1.x                | 3.1.x                     |

## User Guide

### API Definition

| Name       | Path               | Method |
|------------|--------------------|--------|
| connect    | /api/db/connect    | POST   |
| exec       | /api/db/exec       | POST   |
| disconnect | /api/db/disconnect | POST   |

#### Connect API ####

The request json body

```json
{
  "username": "user",
  "password": "password",
  "address": "192.168.8.26",
  "port": 9669
}
```

The description of the parameters is as follows:

| Field    | Description                                                                                                                 |
|----------|-----------------------------------------------------------------------------------------------------------------------------|
| username | Sets the username of your Nebula Graph account. Before enabling authentication, you can use any characters as the username. |
| password | Sets the password of your Nebula Graph account. Before enabling authentication, you can use any characters as the password. |
| address  | Sets the IP address of the graphd service.                                                                                  |
| port     | Sets the port number of the graphd service. The default port number is 9669.                                                |

```bash
$ curl -i -X POST \
    -d '{"username":"user","password":"password","address":"192.168.8.26","port":9669}' \
    http://127.0.0.1:8080/api/db/connect
```

response:

```
HTTP/1.1 200 OK
Content-Length: 100
Content-Type: application/json; charset=utf-8
Server: beegoServer:1.12.3
Set-Cookie: common-nsid=bec2e665ba62a13554b617d70de8b9b9; Path=/; HttpOnly
Set-Cookie: Secure=true; Path=/
Set-Cookie: SameSite=None; Path=/
Date: Fri, 02 Apr 2021 08:49:18 GMT

{
  "code": 0,
  "data": "5e18fa40-5343-422f-84e3-e7f9cad6b735",
  "message": "Login successfully"
}
```

Notice:

The response data nsid `5e18fa40-5343-422f-84e3-e7f9cad6b735` is encoded by HMAC-SH256 encryption algorithm, so it's not the same as what you get from a cookie.
If you connect to the graphd service successfully, remember to save the *NSID* locally, which is important for the *exec* api to execute nGQL.
If you restart the gateway server, all authenticated session will be lost, please be aware of this.

#### Exec API ####

The requested json body

```json
{
  "gql": "show spaces;"
}
```

**Cookie** is required in the request header to request `exec` api.


```bash
$ curl -X POST \
    -H "Cookie: SameSite=None; common-nsid=bec2e665ba62a13554b617d70de8b9b9" \
    -d '{"gql": "show spaces;"}' \
    http://127.0.0.1:8080/api/db/exec
```

response:

```json
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
}
```

#### Disconnect API ####

```bash
$ curl -X POST http://127.0.0.1:8080/api/db/disconnect
```

response:

```json
{
  "code": 0,
  "data": null,
  "message": "Disconnect successfully"
}
```

#### Import API #### 

The requested json body

```json
{
  "configPath": "examples/v2/example.yaml"
}
```

The description of the parameters is as follows.

| Field      | Description                                                  |
| ---------- | ------------------------------------------------------------ |
| configPath | `configPath` is a relative path that under the `uploadspath` in `app.conf`. |
| configBody | `configBody` is the detail configuration with JSON format (instead of YAML format).|

If you choose to use `configPath`, you need to make sure that the config file has been uploaded to `uploadspath`.

```bash
$ curl -X POST -d '{"configPath": "./examples/v2/example.yaml","configBody": {}}' http://127.0.0.1:8080/api/task/import
```

If you choose to use `configBody`, you need to set the `configPath` value to `""` and set the `configBody` as JSON format.

response:

```json
{
    "code": 0,
    "data": [
        "1"
    ],
    "message": "Import task 1 submit successfully"
}
```

#### Action API ####

The requested json body

```json
{
  "taskID": "1",
  "taskAction": "actionQuery"
}
```

The description of the parameters is as follows.

| Field      | Description                                          |
| ---------- | ---------------------------------------------------- |
| taskID     | Set the task id to do task action                    |
| taskAction | The task action enums, include: actionQuery, actionQueryAll, actionStop, actionStopAll, etc. |

```bash
$ curl -X POST -d '{"taskID": "1", "taskAction": "actionQuery"}' http://127.0.0.1:8080/api/task/import/action
```

response:

```json
{
  "code": 0,
  "data": {
    "results": [
      {
        "taskID": "1",
        "taskStatus": "statusProcessing"
      }
    ],
    "msg": "Task is processing"
  },
  "message": "Processing a task action successfully"
}
```

