# Rename重命名插件
使用番剧信息拼接重命名后路径


## 插件设置

### scrape
全局变量`scrape`，bool类型，可选，默认`True`  
是否重命名完成后进行刮削，如写入`tvshow.nfo`文件


## 入口函数
### rename
`rename(args) -> dict`

#### 入参
入参args为dict类型  
结构为：
```python
args = {
    "anime": {
        "id": 329114, 
        "themoviedb_id": 119495, 
        "mikan_id": 2822, 
        "name": "陰の実力者になりたくて！", 
        "name_cn": "想要成为影之实力者！", 
        "season": 1, 
        "ep": 19, 
        "ep_type": 1,
        "eps": 20, 
        "air_date": "2022-10-05"
    },
    "filename": "[LoliHouse] Kage no Jitsuryokusha ni Naritakute! - 19 [WebRip 1080p HEVC-10bit AAC][简繁内封字幕].mp4",
    "filepath": "想要成为影之实力者！/S01/E19.mp4"
}
```
其中  
`anime`: 动画信息  
`anime.ep_type`: 剧集类型。0: 无法解析, 1: 正常剧集, 2: SP。剧集类型为0时，`anime.ep`值为0  
`filename`: 下载默认文件名（原文件名），不包含路径  

#### 返回值
`error`: 必要，错误信息，为None则没有错误  
`filename`: 必要，重命名后的目标文件路径。最终将会保存到`{save_path}/{filepath}`    
`dir`: 可选，动画根目录。如为空或文件夹不存在，将会使用`filename`路径的最顶层路径      

```python
def rename(args):
    # ...
    return {
        "error": None,
        "filename": "想要成为影之实力者！/S01/E19.mp4",
        "dir": "想要成为影之实力者！"
    }
```
上面的返回值最终将会产生以下文件夹结构：  
```text
|--想要成为影之实力者！
   |--tvshow.nfo
   |--S01
      |--E19.mp4
```