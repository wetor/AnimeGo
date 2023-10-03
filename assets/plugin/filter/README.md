# Filter过滤插件
在添加下载项前的过滤操作
## 入口函数
### filter_all

`filter_all(args) -> dict`
#### 入参
入参args为dict类型，筛选前的列表。
结构为：
```python
args = {
    "items": [
        {
            "date": "2023-02-15",
            "torrent_url": "https://mikanani.me/Download/20230215/301289d8fcde1751a26f55b4dc156da7f9ca1a6f.torrent",
            "length": 237083040,
            "name": "[OPFans枫雪动漫][ONE PIECE 海贼王][第1051话][周日版][720p][MP4]",
            "type": "application/x-bittorrent",
            "mikan_url": "https://mikanani.me/Home/Episode/301289d8fcde1751a26f55b4dc156da7f9ca1a6f"
        },
        {...},
    ]
}
```
#### 返回值
返回值为dict类型，返回筛选后的结果。 

```python
def filter_all(args):
    # ...
    return {
        "index": [0, 1, 4],
        "error": None
    }
```

