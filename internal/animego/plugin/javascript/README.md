## 数据结构
### main
js中main函数的参数：
```javascript
// 参数
argv = [
    {
        Url: "https://mikanani.me/Home/Episode/0c0a30b9b7ee437e33fdea6522eb223377dd1d48", // Link，详情页连接，用于下一步解析番剧信息
        Name: "", // 种子名
        Date: "", // 发布日期
        Torrent: "", // 种子连接
        Hash: "", // 种子hash，唯一ID
    },
    // ...
]
// 返回值
main = {
    index: [0, 1, 2], // 筛选结果
    error: null,
}
```
### 日志
```javascript
// 日志输出
log.debug(...params)
log.info(...params)
log.error(...params)
```
### animeGo.parseName
```javascript
// 初步解析资源名
// 参数
name = "" // 种子名
// 返回值
animeGo.parseName = {
    TitleRaw: "", // 种子名
    Name: "", // 番剧名
    Season: 0, // 季度
    Ep: 0, // ep
    Group: "", // 字幕组
    Definition: "", // 分辨率
    Sub: "", // 字幕语言
    Source: "", // 资源平台
}
```
### animeGo.getMikanInfo
```javascript
// 获取Mikan信息
// 参数
url = "https://mikanani.me/Home/Episode/0c0a30b9b7ee437e33fdea6522eb223377dd1d48" // mikanUrl
// 返回值
animeGo.getMikanInfo = {
    ID: 0,
    SubGroupID: 0,
    PubGroupID: 0,
    GroupName: ""
}
```

### 其他函数
#### sleep
```javascript
sleep(ms) // ms 毫秒，1000ms=1s
```

#### print
```javascript
print(...params)
```