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
            "download": "https://mikanani.me/Download/20230215/301289d8fcde1751a26f55b4dc156da7f9ca1a6f.torrent",
            "length": 237083040,
            "name": "[OPFans枫雪动漫][ONE PIECE 海贼王][第1051话][周日版][720p][MP4]",
            "type": "application/x-bittorrent",
            "url": "https://mikanani.me/Home/Episode/301289d8fcde1751a26f55b4dc156da7f9ca1a6f"
        },
        {...},
    ]
}
```
#### 返回值
返回值为dict类型，返回筛选后的结果。  
结构分为两种：  
1.  只返回筛选后入参的索引
    ```python
    return {
        "index": [0, 1, 4],
        "error": None
    }
    ```

2.  返回筛选后入参的索引的同时，携带`item.name`标题解析结果  
    标题`item.name`解析可以使用 [Auto_Bangumi.raw_parser.py](Auto_Bangumi/raw_parser.py) 中提供的`analyse`方法。  
    也可以自行解析，其中`parsed`字段中的`episode`字段为**必要**项，其余字段均**不必要**。
    ```python
    from Auto_Bangumi.raw_parser import analyse
    parsed = analyse("[猎户不鸽压制] 万事屋斋藤先生转生异世界 / 斋藤先生无所不能 Benriya Saitou-san, Isekai ni Iku [03] [1080p] [繁中内嵌] [2023年1月番]")
    # parsed = {
    #     "episode": 3,
    #     "group": "猎户不鸽压制",
    #     "resolution": "1080p",
    #     "season": 1,
    #     "season_raw": "",
    #     "source": "",
    #     "sub": "繁中内嵌",
    #     "title_en": "斋藤先生无所不能 Benriya Saitou-san, Isekai ni Iku",
    #     "title_jp": "",
    #     "title_zh": "万事屋斋藤先生转生异世界"
    # }
    ```
    返回结构如下：
    ```python
    return {
        "data": [
            {
                "index": 0,
                "parsed": {
                    "episode": 3,
                    "group": "猎户不鸽压制",
                    "resolution": "1080p",
                    "season": 1,
                    "season_raw": "",
                    "source": "",
                    "sub": "繁中内嵌",
                    "title_en": "斋藤先生无所不能 Benriya Saitou-san, Isekai ni Iku",
                    "title_jp": "",
                    "title_zh": "万事屋斋藤先生转生异世界"
                }
            },
            {...},
        ]
        "error": None
    }
    ```
    
