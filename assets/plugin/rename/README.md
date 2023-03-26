# Rename重命名插件
使用番剧信息拼接重命名后路径


## 插件设置

### \_\_write_tvshow\_\_
全局变量`__write_tvshow__`，bool类型，可选，默认`True`  
是否写入Jellyfin的`tvshow.nfo`文件  
启用后将会写入返回值中`tvshow_dir`所设置路径。如未返回`tvshow_dir`，将写入目标文件路径`filepath`的**上上**层目录


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
        "eps: 20, 
        "air_date": "2022-10-05"
    },
    "filename": "[LoliHouse] Kage no Jitsuryokusha ni Naritakute! - 19 [WebRip 1080p HEVC-10bit AAC][简繁内封字幕].mp4",
    "default_filepath": "想要成为影之实力者！/S01/E19.mp4"
}
```
其中  
`anime`: 动画信息  
`filename`: 下载默认文件名  
`default_filepath`: 默认的保存文件名，可直接返回

#### 返回值
`error`: 必要，错误信息，为None则没有错误  
`filepath`: 必要，重命名后的文件路径。最终将会保存到`save_path/{filepath}`    
`tvshow_dir`: 可选，全局变量`__write_tvshow__`开启后，将会把`tvshow.nfo`文件写到此文件夹中。如为空或文件夹不存在，将会使用`filepath`的**上上**层目录      

```python
def rename(args):
    # ...
    return {
        "error": None,
        "filepath": "想要成为影之实力者！/S01/E19.mp4",
        "tvshow_dir": "想要成为影之实力者！"
    }
```
上面的返回值最终将会产生以下文件夹结构：  
```text
|--想要成为影之实力者！
   |--tvshow.nfo
   |--S01
      |--E19.mp4
```