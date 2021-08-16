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
$ curl -i \
    -X POST \
    -d '{"username":"user","password":"password","address":"192.168.8.26","port":9669}' \
    http://127.0.0.1:8080/api/db/connect
```

response:

```
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
    -H "Cookie: SameSite=None; nsid=bec2e665ba62a13554b617d70de8b9b9" \
    -H "nsid: bec2e665ba62a13554b617d70de8b9b9" \
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
| configBody | `configBody` is the detail configuration.|

If you choose to use `configPath`, you need to make sure that the config file has been uploaded to `uploadspath`.

```bash
$ curl -X POST -d '{"configPath": "./examples/v2/example.yaml","configBody": {}}' http://127.0.0.1:8080/api/task/import
```

If you choose to use `configBody`, you need to set the `configPath` value to `""` and set the configBody as JSON format:

<details><summary><strong>Please to see the request body example</strong></summary>

```json
{
  "configPath": "",
  "configBody": {
    "version": "v2",
    "description": "example",
    "removeTempFiles": false,
    "clientSettings": {
      "retry": 3,
      "concurrency": 2,
      "channelBufferSize": 1,
      "space": "importer_test",
      "connection": {
        "user": "root",
        "password": "nebula",
        "address": "graphd1:9669,graphd2:9669"
      },
      "postStart": {
        "commands": "UPDATE CONFIGS storage:wal_ttl=3600;\nUPDATE CONFIGS storage:rocksdb_column_family_options = { disable_auto_compactions = true };\nDROP SPACE IF EXISTS importer_test;\nCREATE SPACE IF NOT EXISTS importer_test(partition_num=5, replica_factor=1, vid_type=FIXED_STRING(10));\nUSE importer_test;\nCREATE TAG course(name string, credits int);\nCREATE TAG building(name string);\nCREATE TAG student(name string, age int, gender string);\nCREATE EDGE follow(likeness double);\nCREATE EDGE choose(grade int);\nCREATE TAG course_no_props();\nCREATE TAG building_no_props();\nCREATE EDGE follow_no_props();\n",
        "afterPeriod": "8s"
      },
      "preStop": {
        "commands": "UPDATE CONFIGS storage:rocksdb_column_family_options = { disable_auto_compactions = false };\nUPDATE CONFIGS storage:wal_ttl=86400;\n"
      }
    },
    "logPath": "./uploads/examples/v2/err/test.log",
    "files": [
      {
        "path": "./uploads/examples/v2/choose.csv",
        "batchSize": 2,
        "inOrder": false,
        "type": "csv",
        "csv": {
          "withHeader": false,
          "withLabel": false
        },
        "schema": {
          "type": "edge",
          "edge": {
            "name": "choose",
            "withRanking": false,
            "props": [
              {
                "name": "grade",
                "type": "int"
              }
            ]
          }
        }
      },
      {
        "path": "./uploads/examples/v2/course.csv",
        "failDataPath": "./uploads/examples/v2/err/course.csv",
        "batchSize": 2,
        "inOrder": true,
        "type": "csv",
        "csv": {
          "withHeader": false,
          "withLabel": false
        },
        "schema": {
          "type": "vertex",
          "vertex": {
            "tags": [
              {
                "name": "course",
                "props": [
                  {
                    "name": "name",
                    "type": "string"
                  },
                  {
                    "name": "credits",
                    "type": "int"
                  }
                ]
              },
              {
                "name": "building",
                "props": [
                  {
                    "name": "name",
                    "type": "string"
                  }
                ]
              }
            ]
          }
        }
      },
      {
        "path": "./uploads/examples/v2/course-with-header.csv",
        "failDataPath": "./uploads/examples/v2/err/course-with-header.csv",
        "batchSize": 2,
        "inOrder": true,
        "type": "csv",
        "csv": {
          "withHeader": true,
          "withLabel": true
        },
        "schema": {
          "type": "vertex"
        }
      },
      {
        "path": "./uploads/examples/v2/follow-with-label.csv",
        "failDataPath": "./uploads/examples/v2/err/follow-with-label.csv",
        "batchSize": 2,
        "inOrder": true,
        "type": "csv",
        "csv": {
          "withHeader": false,
          "withLabel": true
        },
        "schema": {
          "type": "edge",
          "edge": {
            "name": "follow",
            "withRanking": true,
            "srcVID": {
              "index": 0
            },
            "dstVID": {
              "index": 2
            },
            "rank": {
              "index": 3
            },
            "props": [
              {
                "name": "likeness",
                "type": "double",
                "index": 1
              }
            ]
          }
        }
      },
      {
        "path": "./uploads/examples/v2/follow-with-label-and-str-vid.csv",
        "failDataPath": "./uploads/examples/v2/err/follow-with-label-and-str-vid.csv",
        "batchSize": 2,
        "inOrder": true,
        "type": "csv",
        "csv": {
          "withHeader": false,
          "withLabel": true
        },
        "schema": {
          "type": "edge",
          "edge": {
            "name": "follow",
            "withRanking": true,
            "srcVID": {
              "index": 0
            },
            "dstVID": {
              "index": 2
            },
            "rank": {
              "index": 3
            },
            "props": [
              {
                "name": "likeness",
                "type": "double",
                "index": 1
              }
            ]
          }
        }
      },
      {
        "path": "./uploads/examples/v2/follow.csv",
        "failDataPath": "./uploads/examples/v2/err/follow.csv",
        "batchSize": 2,
        "type": "csv",
        "csv": {
          "withHeader": false,
          "withLabel": false
        },
        "schema": {
          "type": "edge",
          "edge": {
            "name": "follow",
            "withRanking": true,
            "props": [
              {
                "name": "likeness",
                "type": "double"
              }
            ]
          }
        }
      },
      {
        "path": "./uploads/examples/v2/follow-with-header.csv",
        "failDataPath": "./uploads/examples/v2/err/follow-with-header.csv",
        "batchSize": 2,
        "type": "csv",
        "csv": {
          "withHeader": true,
          "withLabel": false
        },
        "schema": {
          "type": "edge",
          "edge": {
            "name": "follow",
            "withRanking": true
          }
        }
      },
      {
        "path": "./uploads/examples/v2/student.csv",
        "failDataPath": "./uploads/examples/v2/err/student.csv",
        "batchSize": 2,
        "type": "csv",
        "csv": {
          "withHeader": false,
          "withLabel": false
        },
        "schema": {
          "type": "vertex",
          "vertex": {
            "tags": [
              {
                "name": "student",
                "props": [
                  {
                    "name": "name",
                    "type": "string"
                  },
                  {
                    "name": "age",
                    "type": "int"
                  },
                  {
                    "name": "gender",
                    "type": "string"
                  }
                ]
              }
            ]
          }
        }
      },
      {
        "path": "./uploads/examples/v2/student.csv",
        "failDataPath": "./uploads/examples/v2/err/student_index.csv",
        "batchSize": 2,
        "type": "csv",
        "csv": {
          "withHeader": false,
          "withLabel": false
        },
        "schema": {
          "type": "vertex",
          "vertex": {
            "vid": {
              "index": 1
            },
            "tags": [
              {
                "name": "student",
                "props": [
                  {
                    "name": "age",
                    "type": "int",
                    "index": 2
                  },
                  {
                    "name": "name",
                    "type": "string",
                    "index": 1
                  },
                  {
                    "name": "gender",
                    "type": "string"
                  }
                ]
              }
            ]
          }
        }
      },
      {
        "path": "./uploads/examples/v2/student-with-label-and-str-vid.csv",
        "failDataPath": "./uploads/examples/v2/err/student_label_str_vid.csv",
        "batchSize": 2,
        "type": "csv",
        "csv": {
          "withHeader": false,
          "withLabel": true
        },
        "schema": {
          "type": "vertex",
          "vertex": {
            "vid": {
              "index": 1
            },
            "tags": [
              {
                "name": "student",
                "props": [
                  {
                    "name": "age",
                    "type": "int",
                    "index": 2
                  },
                  {
                    "name": "name",
                    "type": "string",
                    "index": 1
                  },
                  {
                    "name": "gender",
                    "type": "string"
                  }
                ]
              }
            ]
          }
        }
      },
      {
        "path": "./uploads/examples/v2/follow.csv",
        "failDataPath": "./uploads/examples/v2/err/follow_index.csv",
        "batchSize": 2,
        "limit": 3,
        "type": "csv",
        "csv": {
          "withHeader": false,
          "withLabel": false
        },
        "schema": {
          "type": "edge",
          "edge": {
            "name": "follow",
            "srcVID": {
              "index": 0
            },
            "dstVID": {
              "index": 1
            },
            "rank": {
              "index": 2
            },
            "props": [
              {
                "name": "likeness",
                "type": "double",
                "index": 3
              }
            ]
          }
        }
      },
      {
        "path": "./uploads/examples/v2/follow-delimiter.csv",
        "failDataPath": "./uploads/examples/v2/err/follow-delimiter.csv",
        "batchSize": 2,
        "type": "csv",
        "csv": {
          "withHeader": true,
          "withLabel": false,
          "delimiter": "|"
        },
        "schema": {
          "type": "edge",
          "edge": {
            "name": "follow",
            "withRanking": true
          }
        }
      },
      {
        "path": "./uploads/examples/v2/follow.csv",
        "failDataPath": "./uploads/examples/v2/err/follow_http.csv",
        "batchSize": 2,
        "limit": 3,
        "type": "csv",
        "csv": {
          "withHeader": false,
          "withLabel": false
        },
        "schema": {
          "type": "edge",
          "edge": {
            "name": "follow",
            "srcVID": {
              "index": 0
            },
            "dstVID": {
              "index": 1
            },
            "rank": {
              "index": 2
            },
            "props": [
              {
                "name": "likeness",
                "type": "double",
                "index": 3
              }
            ]
          }
        }
      },
      {
        "path": "./uploads/examples/v2/course.csv",
        "failDataPath": "./uploads/examples/v2/err/course-empty-props.csv",
        "batchSize": 2,
        "inOrder": true,
        "type": "csv",
        "csv": {
          "withHeader": false,
          "withLabel": false,
          "delimiter": ","
        },
        "schema": {
          "type": "vertex",
          "vertex": {
            "vid": {
              "index": 0
            },
            "tags": [
              {
                "name": "course_no_props"
              }
            ]
          }
        }
      },
      {
        "path": "./uploads/examples/v2/course.csv",
        "failDataPath": "./uploads/examples/v2/err/course-multi-empty-props.csv",
        "batchSize": 2,
        "inOrder": true,
        "type": "csv",
        "csv": {
          "withHeader": false,
          "withLabel": false,
          "delimiter": ","
        },
        "schema": {
          "type": "vertex",
          "vertex": {
            "vid": {
              "index": 0
            },
            "tags": [
              {
                "name": "course_no_props"
              },
              {
                "name": "building_no_props"
              }
            ]
          }
        }
      },
      {
        "path": "./uploads/examples/v2/course.csv",
        "failDataPath": "./uploads/examples/v2/err/course-mix-empty-props.csv",
        "batchSize": 2,
        "inOrder": true,
        "type": "csv",
        "csv": {
          "withHeader": false,
          "withLabel": false,
          "delimiter": ","
        },
        "schema": {
          "type": "vertex",
          "vertex": {
            "vid": {
              "index": 0
            },
            "tags": [
              {
                "name": "course_no_props"
              },
              {
                "name": "building",
                "props": [
                  {
                    "name": "name",
                    "type": "string",
                    "index": 3
                  }
                ]
              }
            ]
          }
        }
      },
      {
        "path": "./uploads/examples/v2/course.csv",
        "failDataPath": "./uploads/examples/v2/err/course-mix-empty-props-2.csv",
        "batchSize": 2,
        "inOrder": true,
        "type": "csv",
        "csv": {
          "withHeader": false,
          "withLabel": false,
          "delimiter": ","
        },
        "schema": {
          "type": "vertex",
          "vertex": {
            "vid": {
              "index": 0
            },
            "tags": [
              {
                "name": "building",
                "props": [
                  {
                    "name": "name",
                    "type": "string",
                    "index": 3
                  }
                ]
              },
              {
                "name": "course_no_props"
              }
            ]
          }
        }
      },
      {
        "path": "./uploads/examples/v2/follow.csv",
        "failDataPath": "./uploads/examples/v2/err/follow-empty-props.csv",
        "batchSize": 2,
        "type": "csv",
        "csv": {
          "withHeader": false,
          "withLabel": false,
          "delimiter": ","
        },
        "schema": {
          "type": "edge",
          "edge": {
            "name": "follow_no_props",
            "withRanking": false,
            "dstVID": {
              "index": 1
            },
            "srcVID": {
              "index": 0
            }
          }
        }
      }
    ]
  }
}
```
</details>


response:

```json
{
  "code": 0,
  "data": null,
  "message": "Import task 0 submit successfully"
}
```

#### Action API ####

The requested json body

```json
{
  "taskID": 0,
  "taskAction": "stopAll"
}
```

The description of the parameters is as follows.

| Field      | Description                                          |
| ---------- | ---------------------------------------------------- |
| taskID     | Set the task id to do task action                    |
| taskAction | Enums, include: stop, stopAll, query, queryAll, etc. |

```bash
$ curl -X POST -d '{"taskID": 0, "taskAction": "stopAll"}' http://127.0.0.1:8080/api/task/import/action
```

response:

```json
{
    "code": 0,
    "data": {
        "taskIDs": [
            "0"
        ],
        "taskStatus": "Task stop successfully"
    },
    "message": "Processing task action successfully"
}
```

