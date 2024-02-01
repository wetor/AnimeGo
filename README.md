# AnimeGo

使用Golang编写的全自动追番工具，简单的部署和使用，方便的模块化扩展

## 使用帮助
```text
  -backup
        配置文件升级前是否备份 (default true)
  -config string
        配置文件路径；配置文件中的相对路径均是相对与程序的位置 (default "data/animego.yaml")
  -debug
        Debug模式，将会显示更多的日志
  -web
        启用Web API (default true)
```
### [可选]0 安装和使用MikanTool插件
- 配置文件中`setting/filter/plugin`新增如下内容后即可开启插件。
  ```
  plugin:
    filter:
      - enable: true
        type: py
        file: filter/mikan_tool.py
        args: {}
        vars: {}
  ```
- 然后 [安装MikanTool Tampermonkey插件](https://greasyfork.org/zh-CN/scripts/449596) ，需要浏览器中已安装Tampermonkey（油猴插件）  
- 具体过滤设置根据油猴插件面板的要求进行  

### 1 首次启动：释放资源、升级配置
```shell
./AnimeGo
```

会在程序所在目录输出`data`文件夹，其中`data/animego.yaml`为配置文件。
### 2 修改配置
打开并编辑`data/animego.yaml`

其中主要需要修改的配置项为：
- `setting.client.*` : 必选，客户端参数，根据实际情况填写
- `setting.download_path`: 下载客户端的下载文件夹
- `setting.save_path`: 重命名后移动到位置。此时将会改名
- `plugin.feed`中的`builtin_mikan_rss.py`插件 : 可选，内置自动订阅插件
  - `vars.url`: 订阅地址，如Mikan的Rss订阅地址
  - `vars.cron`: 订阅时间，Cron格式，参考[Feed订阅插件帮助](assets/plugin/feed/README.md)
- 其余配置项根据需求修改

### 3 启动程序
```shell
./AnimeGo
```
> 可选`-debug`，启用后将输出更详细的日志

可以使用如screen等工具放至后台执行，也可以创建服务并启动
### 4 高级使用
[插件函数文档](assets/plugin/README.md)  
## 文档
1. 配置文件，参考注释
2. [插件函数文档](assets/plugin/README.md)
3. [webapi(Swagger)接口文档](internal/web/README.md)

## 目的
- 简化部署和使用，以及模块化扩展
- 学习

## 目前进度
- 可使用配置、筛选和下载等功能
- 比较高度自由的插件配置  
- python编写插件 [帮助文档](assets/plugin/README.md) 
  - 定时订阅，[帮助文档](assets/plugin/feed/README.md)
  - 筛选器，[帮助文档](assets/plugin/filter/README.md)
  - 重命名规则，[帮助文档](assets/plugin/rename/README.md)
  - 定时任务，[帮助文档](assets/plugin/schedule/README.md)
- 支持Tampermonkey(油猴)插件 [AnimeGo\[Mikan快速订阅\]](https://greasyfork.org/zh-CN/scripts/449596) 快速订阅下载
- Jellyfin支持
- qBittorrent支持
- Transmission支持

## 开发计划
- [x] 增加读取网站离线Archive的缓存功能 降低网站请求
  - [x] [Bangumi数据](https://github.com/bangumi/Archive)
  - [ ] [Mikan数据](https://github.com/MikanProject/bangumi-data/blob/master/dist/data.json)
- [x] [Mikan Project](https://mikanani.me) 订阅支持
- [x] [Jellyfin](https://jellyfin.org/) 媒体库软件识别 会写入bgmid到tvshow.nfo 可以配合[jellyfin-plugin-bangumi](https://github.com/kookxiang/jellyfin-plugin-bangumi)使用
- [ ] 多种下载器支持
  - [x] [qBittorrent](https://qbittorrent.org) 支持
  - [x] [Transmission](https://transmissionbt.com/) 支持
  - [ ] [Aria2](https://aria2.github.io/) 支持
- [ ] Web界面支持
- [x] 模块化与高级自定义功能支持
  - [x] 独立的订阅支持
  - [x] 独立下载控制
  - [x] 自定义订阅、过滤、筛选、解析和重命名插件
  - [x] 自定义定时任务

## 开发日志
 [Changelog.md](Changelog.md) 
 