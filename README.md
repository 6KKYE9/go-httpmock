# go-httpmock

HTTP Mock 服务器。JSON 配置定义接口规则，用于前端调试或测试。

## 用法

先创建 `mock.json`：
```json
{
  "name": "test-api",
  "port": 8888,
  "rules": [
    {"method": "GET", "path": "/api/health", "status": 200, "body": "{\"ok\":true}"},
    {"method": "POST", "path": "/api/users", "status": 201, "body": "{\"id\":1}"}
  ]
}
```

然后启动：
```
go-httpmock -c mock.json
go-httpmock -c mock.json -port 9999
```
