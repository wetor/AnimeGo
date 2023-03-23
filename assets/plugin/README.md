# AnimeGo插件扩展
插件使用python编写，语法版本为python3.4。python解析器实现使用 [gpython](https://github.com/go-python/gpython) 

与标准的python3.4存在**较大差异**，主要体现为：
- 缺少大量的内置方法
- 缺少内置模块

## 插件编写
所有的插件都至少有一个入口函数，并且有且仅有一个入参，为dict类型

插件至多有一个返回值，为dict类型

不同类型的插件对于入口函数、全局变量和返回值有不同要求，参照各个类型插件文档

## 内置模块

### core 核心模块

提供AnimeGo接口的模块


#### core.loads

`loads(s, type='json') -> dict`  
支持type='yaml'，解析yaml格式字符串为dict  
示例：
```python
import core

d = core.loads('{"a": 123, "b": "ccc"}')
# d = {'a': 123, 'b': 'ccc'}
```
#### core.dumps

`dumps(obj, type='json') -> str`  
支持type='yaml'，dict转成yaml格式字符串  
示例：  
```python
import core

s = core.dumps({'a': 123, 'b': 'ccc'})
# s = '{"a": 123, "b": "ccc"}'
```

#### core.parse_mikan

`parse_mikan(url) -> dict`  
解析mikan单集url，url格式为`https://mikanani.me/Home/Episode/****`  

```python
import core

result = core.parse_mikan('https://mikanani.me/Home/Episode/a6f48155e7648a945e9bf85949c6cf8d8eb7ad61')
# result = {
#     "group_name": "OPFans枫雪动漫",
#     "id": 228,
#     "pub_group_id": 1,
#     "sub_group_id": 1
# }
```


#### core.parse_mikan_rss

`parse_mikan_rss(raw_data) -> dict`  
解析mikan的rss xml内容，其他网站的rss的xml中需要含有以下必要字段  
xml格式为：
```xml
<rss version="2.0">
    <channel>
        <title>Mikan Project - 海贼王</title>
        <link>http://mikanani.me/RSS/Bangumi?bangumiId=228&amp;subgroupid=1</link>
        <description>Mikan Project - 海贼王</description>
        <item>
            <guid isPermaLink="false">[OPFans枫雪动漫][ONE PIECE 海贼王][第1051话][周日版][720p][MP4]</guid>
            <link>https://mikanani.me/Home/Episode/301289d8fcde1751a26f55b4dc156da7f9ca1a6f</link>
            <title>[OPFans枫雪动漫][ONE PIECE 海贼王][第1051话][周日版][720p][MP4]</title>
            <description>[OPFans枫雪动漫][ONE PIECE 海贼王][第1051话][周日版][720p][MP4][226.1MB]</description>
            <torrent xmlns="https://mikanani.me/0.1/">
                <link>https://mikanani.me/Home/Episode/301289d8fcde1751a26f55b4dc156da7f9ca1a6f</link>
                <contentLength>237083040</contentLength>
                <pubDate>2023-02-15T18:19:00</pubDate>
            </torrent>
            <enclosure type="application/x-bittorrent" length="237083040"
                       url="https://mikanani.me/Download/20230215/301289d8fcde1751a26f55b4dc156da7f9ca1a6f.torrent"/>
        </item>
        <item>
            <guid isPermaLink="false"><!--无用--></guid>
            <link><!--url，必要，mikan信息也，用于获取mikan信息--></link>
            <title><!--name，必要，torrent名，用于解析番剧各种信息用--></title>
            <description><!--无用--></description>
            <torrent xmlns="无用">
                <link><!--无用--></link>
                <contentLength><!--无用--></contentLength>
                <pubDate><!--date，可选，torrent放出时间--></pubDate>
            </torrent>
            <enclosure type="可选" length="length，可选，种子大小"
                       url="download，必要，torrent下载链接"/>
            <!--enclosure，download，必要，url为torrent下载链接-->
        </item>
        <item>
            <!-- ... -->
        </item>
    </channel>
</rss>
```
示例：
```python
import core

# rss内容
xml_data = '<rss version="2.0"><channel><title>Mikan Project - 海贼王</title>...'

result = core.parse_mikan_rss(xml_data)
# result = [
#     {
#         "date": "2023-02-15",
#         "download": "https://mikanani.me/Download/20230215/301289d8fcde1751a26f55b4dc156da7f9ca1a6f.torrent",
#         "length": 237083040,
#         "name": "[OPFans枫雪动漫][ONE PIECE 海贼王][第1051话][周日版][720p][MP4]",
#         "type": "application/x-bittorrent",
#         "url": "https://mikanani.me/Home/Episode/301289d8fcde1751a26f55b4dc156da7f9ca1a6f"
#     },
#     {...},
# ]
```


### log 日志模块

log模块与print不同，使用的时AnimeGo的日志输出，将会写入日志文件

```python
import log

# 不定长参数
log.info(...)
log.debug(...)
log.warn(...)
log.error(...)

# 其中format为go的格式化方式
log.infof(format, ...)
log.debugf(format, ...)
log.warnf(format, ...)
log.errorf(format, ...)
```
其中go格式化文档参照：[https://pkg.go.dev/fmt](https://pkg.go.dev/fmt)


## 插件配置
可以在脚本中直接声明变量并赋值，也可以配置文件中为vars设置值，格式如下：
```yaml
 - enable: false
   type: builtin
   file: builtin_mikan_rss.py
   args: 
     test: input_test
   vars:
     __cron__: 0 0/20 * * * ?
     __name__: Example
     __url__: https://example.com/
``` 
`args`将会追加并覆盖 `builtin_mikan_rss.py` 中**入口函数**的参数`args['test']`  
`vars`将会追加并覆盖 `builtin_mikan_rss.py` 中的全局变量`__cron__`,`__name__`,`__url__`变量  
其他字段参考配置文件中的注释  

## 通用的内置变量和函数
所有被AnimeGo载入的插件，都将拥有以下变量或方法
### \_\_debug\_\_
全局变量`__debug__`，bool类型，当前是否为debug模式   

### \_\_plugin_name\_\_
全局变量`__plugin_name__`，插件文件名，不含扩展名

### \_\_plugin_dir\_\_
全局变量`__plugin_dir__`，插件所在目录，绝对路径

### \_\_animego_version\_\_
全局变量`__animego_version__`，AnimeGo版本号，`vx.x.x`格式，如`v1.0.0`

### \_get_config()
内置函数`_get_config`，获取插件配置  
`_get_config() -> dict`  
读取插件所在目录下，同名的**yaml**或**json**文件并解析为dict。优先读取yaml，其次json，若都不存在返回空  

```python
conf = _get_config()
```

## Feed订阅插件
解析订阅内容  
[Feed订阅插件帮助](feed/README.md)

## Filter过滤插件
在添加下载项前的过滤操作  
[Filter过滤插件帮助](filter/README.md)

## Schedule定时任务插件
使用Cron表达式定时执行任务的插件
[Schedule定时任务插件帮助](schedule/README.md)
