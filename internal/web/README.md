## AnimeGo WEB Api接口文档

默认开启Swagger：http://localhost:7991/swagger/index.html  
Swagger文档：[swagger.json](../../docs/swagger.json)

### 鉴权
HTPP请求Header中，`Access-Key`为SHA256后的access_key

#### Swagger验证方式
通过 `/sha256`接口，传入配置文件中设置的 `access_key`，返回加密后的`access_key`

### 响应代码
code值所代码含义
```text
200 成功
300 失败
403 鉴权失败
500 服务端未知错误
```