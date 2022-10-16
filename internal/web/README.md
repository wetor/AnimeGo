## AnimeGo WEB Api接口文档

### 响应代码
code值所代码含义
```text
200 成功
300 失败
403 鉴权失败
500 服务端未知错误
```

### 鉴权
HTPP请求Header中，`Access-Key`为SHA256后的access_key

### 测试连通
`GET /ping`
#### 响应
```text
pong bool 默认为true
time int64 AnimeGo所在主机当前时间
```

### 下载项目
`POST /api/rss`
#### 请求
```text
source string 来源
rss object 订阅信息
    url string 订阅链接
is_select_ep bool 是否选中部分ep
ep_links []string 选中的ep链接数组
```


### 上传插件配置
`POST /api/plugin/config`
#### 请求
```text
name string 插件名，如`filter/default.js`
data string 配置文件数据，json格式转成的base64字符串
```

### 获取插件配置
`GET /api/plugin/config`
#### 请求
```text
name string 插件名，如`filter/default.js`
```
#### 响应
```text
name string 插件名，如`filter/default.js`
data string 配置文件数据，json格式转成的base64字符串
```


