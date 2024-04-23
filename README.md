# goctl-swagger
一个根据 `gozero` API描述文件生成 `swagger` 文档的插件。
基于:
```bash
go:     1.19  # 及以上
gozero: 1.6.0 # 及以上
```

支持功能:
- [x] 自定义 `tag` 前缀 
- [x] 自定义外层响应
- [x] 在外层响应中指定响应数据的key

暂不支持:
1. 根据 `tag` 中的 `validate` 标签生成相应的文档 

## 使用指南

### 1. 编译goctl-swagger插件

```
GOPROXY=https://goproxy.cn/,direct go install github.com/henryjhenry/goctl-swagger@latest
```

### 2. 配置环境

将$GOPATH/bin中的goctl-swagger添加到环境变量

### 3. 使用姿势

#### 生成 swagger.json 文件
```shell script
# 在goctl中使用
goctl api plugin -plugin goctl-swagger="swagger -target swagger.json" -api your.api -dir .

# 在本地使用
go run main.go swagger -target swagger.json 0<~/tmp.json
# 或
goctl-swagger swagger -target swagger.json 0<~/tmp.json
```
tmp.json:
```json
{
    "ApiFilePath": "your.api",
    "Dir": "working dir"
}
```

#### 指定Host，basePath，schemes [api-host-and-base-path](https://swagger.io/docs/specification/2-0/api-host-and-base-path/)
```shell script
goctl api plugin -plugin goctl-swagger="swagger -target swagger.json -host 127.0.0.2 -basepath /api -schemes https,wss" -api your.api -dir .
```

#### 指定外层响应
```shell script
goctl api plugin -plugin goctl-swagger="swagger -target swagger.json -outsideSchema ./outsideSchema.api" -api your.api -dir .
```
外层响应示例见 `testdata/api/outside_schema.api`

`goctl-swagger` 会解析并使用指定的 `outsideSchema` 文件并使用第一个声明作为外层响应，并默认使用 `data` 作为内层响应的key，如果需要自定义，请加上 `-responseKey` 参数，例如：`-responseKey response`


### 4. run test
```bash
go test ./render --count=1 -v
``` 
运行完毕后，会在 `testdata` 目录生成 `swagger.json`。可使用IDE对应的插件或自建服务查看文档。(JetBrains已原生支持，VS Code推荐: OpenAPI (Swagger) Editor)

**如有feature或bug，欢迎提交issue 👏🏻**