# Feed订阅解析插件
自定义解析订阅的插件

## [特殊]内置插件
内置自动订阅下载插件  
```yaml
 - enable: false
   type: builtin
   file: builtin_mikan_rss.py
   args: {}
   vars:
     cron: 0 0/20 * * * ?
     name: Example
     url: https://example.com/
``` 
对于内置的自动订阅插件，其中`type`和`file`为固定写法  
启用自定订阅`enable`设置为`true`  
设置`cron`设置定时(cron表达式)，设置`url`设置订阅地址  

## 插件设置

### name
全局变量`name`，str类型，可选  
定时任务的名称，方便日志中区分

### url
全局变量`url`，str类型，必要  
订阅地址，AnimeGo将会使用GET方式请求此地址，将结果转为str类型，传递给入口函数。  
其中`args['data']`为请求响应的body  

### header
全局变量`header`，dict类型，可选  
请求订阅地址时所携带的header请求头，其中`user-agent`默认为AnimeGo信息，无法设置

### cron
全局变量`cron`，str类型，必要  
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
    "status": 200,
    "data": "..."
}
```
其中`args['data']`为`url`的请求结果，str类型  
如果为mikan订阅地址，则可以使用`core.parse_mikan_rss(args['data'])`解析  
也可以选择自行解析  

#### 返回值
返回值为dict类型，返回订阅解析结果。  
其中各个字段的来源参考 [core.parse_mikan_rss 方法](../README.md)  
```python
def parse(args):
    # ...
    return {
        "items": [
            {
                "url":      "https://mikanani.me/Home/Episode/134903ffdc03d1e7b2f3440191ac0f18720a9ff0",
                "name":     "[OPFans枫雪动漫][one_piece 海贼王][1052][1080p]_mkv[周日版]",
                "date":     "2023-02-19",
                "type":     "application/x-bittorrent",
                "download": "https://mikanani.me/Download/20230219/134903ffdc03d1e7b2f3440191ac0f18720a9ff0.torrent",
                "length":   1503238528
            },
            {},
            # ...
        ],
        "error": None
    }
```
