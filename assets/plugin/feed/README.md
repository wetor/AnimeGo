# Feed订阅解析插件
自定义解析订阅的插件

## 插件设置
### \_\_name\_\_
全局变量`__name__`，str类型，可选  
定时任务的名称，方便日志中区分

### \_\_url\_\_
全局变量`__url__`，str类型，必要  
订阅地址，AnimeGo将会下载此地址的文件，传递给入口函数  

### \_\_cron\_\_
全局变量`__cron__`，str类型，必要  
订阅的定时规则，执行完毕后，将会把返回值传递给启用的**过滤器**([过滤器文档](../filter/README.md))，之后进行下载    
支持 [秒] [分] [小时] [日] [月] [周] 六项的Cron表达式    
[Cron表达式文档（维基百科，推荐）](https://zh.wikipedia.org/wiki/Cron)  
[Cron表达式文档（百度百科）](https://baike.baidu.com/item/cron/10952601)


## 入口函数
### parse

`parse(args) -> dict`
#### 入参
入参args为dict类型  
结构为：
```python
args = {
    "data": "..."
}
```
其中`args['data']`为`__url__`的下载内容，str类型  
如果为mikan订阅地址，则可以使用`core.parse_mikan_rss(args['data'])`解析  
也可以选择自行解析  

#### 返回值
返回值为dict类型，返回订阅解析结果。  
其中各个字段的来源参考 [core.parse_mikan_rss 方法](../README.md)  
```python
return {
    "items": [
        {
            url:      "https://mikanani.me/Home/Episode/134903ffdc03d1e7b2f3440191ac0f18720a9ff0",
            name:     "[OPFans枫雪动漫][one_piece 海贼王][1052][1080p]_mkv[周日版]",
            date:     "2023-02-19",
            type:     "application/x-bittorrent",
            download: "https://mikanani.me/Download/20230219/134903ffdc03d1e7b2f3440191ac0f18720a9ff0.torrent",
            length:   1503238528
        },
        {},
        ...
    ],
    "error": None
}
```
