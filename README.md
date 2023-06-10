# AnimeGo

使用Golang编写的全自动追番工具，简单的部署和使用，方便的模块化扩展

## 使用帮助
```text
  -config string
        配置文件路径；配置文件中的相对路径均是相对与程序的位置 (default "./data/animego.yaml")
  -debug
        Debug模式，将会显示更多的日志
  -web
        启用Web API，默认启用 (default true)
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

### 1 首次启动：释放资源
```shell
./AnimeGo
```

### 1.1 升级后首次启动：升级配置、释放资源
```shell
./AnimeGo
```

会在程序所在目录输出`data`文件夹，其中`data/animego.yaml`为配置文件。
### 2 修改配置
打开并编辑`data/animego.yaml`

其中主要需要修改的配置项为：
- `setting.client.qbittorrent` : 必选，qBittorrent客户端webapi信息
- `setting.download_path`: 下载器下载保存位置。临时位置，移动后将会删除
- `setting.save_path`: 重命名后移动到位置。此时将会改名
- `plugin.feed`中的`builtin_mikan_rss.py`插件 : 可选，内置自动订阅插件
  - `vars.__url`: 订阅地址，如Mikan的rss订阅地址
  - `vars.__cron`: 订阅时间，Cron格式，参考[Feed订阅插件帮助](assets/plugin/feed/README.md)
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
- python编写插件
  - 定时订阅，[帮助文档](assets/plugin/feed/README.md)
  - 筛选器，[帮助文档](assets/plugin/filter/README.md)
  - 定时任务，[帮助文档](assets/plugin/schedule/README.md)
- 支持Tampermonkey(油猴)插件 [AnimeGo\[Mikan快速订阅\]](https://greasyfork.org/zh-CN/scripts/449596) 快速订阅下载
- Jellyfin支持

## 开发计划
- [x] 增加读取网站离线Archive的缓存功能 降低网站请求
  - [x] [Bangumi数据](https://github.com/bangumi/Archive)
  - [ ] [Mikan数据](https://github.com/MikanProject/bangumi-data/blob/master/dist/data.json)
- [x] 类似[AutoBangumi](https://github.com/EstrellaXD/Auto_Bangumi) 的追番功能
  - [x] [Mikan Project](https://mikanani.me) 订阅支持
  - [x] [qBittorrent](https://qbittorrent.org) 等下载支持
  - [x] [Jellyfin](https://jellyfin.org/) 媒体库软件识别 会写入bgmid到tvshow.nfo 可以配合[jellyfin-plugin-bangumi](https://github.com/kookxiang/jellyfin-plugin-bangumi)使用
  - [ ] ...
- [ ] Web界面支持
- [ ] 模块化与高级自定义功能支持
  - [ ] 独立的订阅支持
  - [x] 独立下载控制
  - [ ] ...

## 开发日志

## v0.9.2
- **配置文件版本号为`1.5.1`**
- 优化webapi
- 支持websocket
  - 支持websocket查看实时日志
- 修复torrent文件重复下载问题
- 补充单测
- 支持数据源重定向

## v0.9.1
- 修复tmdb解析失败时流程结束的问题
- 优化代码
- py支持format方法

## v0.9.0
- **配置文件版本号为`1.5.0`**
- 支持多内容torrent解析下载
  - 优化下载管理器流程
  - 预解析torrent文件
- 单测改造
  - testdata统一路径
  - 更好的单测覆盖
- 新增parser插件
- 使用go1.20

### v0.8.4
- 修改订阅流程
  - 移除多协程订阅功能
- 优化下载管理器
  - 优化性能
  - 改进流程
- 完善下载管理器单测，覆盖更多场景
- 插件执行器新增重命名插件
- 准确的获取magnet和torrent的BT hash
- 修复qBittorrent4.5.0以上删除错误
- 修复下载管理器部分bug

### v0.8.3
- 优化重命名流程
  - 使用一个全局协程，定时接收下载状态变更的通知，执行对应重命名操作
- 重命名插件化
- 新增rename插件文档

### v0.8.2
- 更新gpython
- 插件支持debug模式
- 新增插件调试工具

### v0.8.1
- feed插件支持设置header

### v0.8.0
- 支持feed订阅插件，可以自行解析rss
- 支持在配置文件中设置插件变量、参数
- 支持内置插件
- 更改插件源码结构
- 补充插件文档
- 补充单测用例

### v0.7.4
- 内部统一使用unix风格路径
- 修复python插件中int类型错误

### v0.7.3
- 修复windows下路径拼接错误

### v0.7.2
- **配置文件版本号为`1.3.0`**
- 更改插件配置文件结构
- 支持配置定时任务插件

### v0.7.1
- 移除JavaScript插件支持

### v0.7.0
- 实现MikanTool(原AnimeGoHelperParser)的Python版插件
- Python插件支持多函数
- Python插件支持全局变量
- Python插件新增core模块，支持json和yaml的编码和解码，以及Mikan url解析
- 统一使用filepath已解决windows上可能存在的问题
- 修复部分Bug

### v0.6.8 
- Schedule定时任务
  - 优化参数传递
  - 支持定时执行插件脚本
- 优化部分单测
  - mock了qbt客户端

### v0.6.7
- 修复js plugin调用py问题
- 增加try方法，修改部分错误处理
- 封装全局log方法
  - 统一使用含format的日志输出
  - 日志文件固定为为INFO级别，不含DEBUG日志
- 定时任务支持失败重试
- AnimeGo版本号编译时设置

### v0.6.6
- 优化过滤器流程
- 启动时覆盖内置脚本
- 支持过滤器返回标题解析结果并在后续中直接使用
  - 可自定义解析标题中的信息

### v0.6.5
- **配置文件版本号为`1.2.0`**
- 支持Python编写过滤插件
- 更改配置文件中过滤器插件的格式
- 修复webapi无法删除cache的问题
- 优化代码

### v0.6.4 优化代码更新
- 移除process、store包
  - 初始化功能全部在main中进行
  - 功能模块中使用通过Init方法传递的配置项
- 优化单测，使用testdata完成单测
- 调整import引用顺序
- 调整部分文件夹结构

### v0.6.3
- 修复一个小bug

### v0.6.2
- 新增更新配置文件web api
- 修改获取配置文件web api
- 修复下载Bangumi缓存定时任务解压失败的问题

### v0.6.1
- 修复python脚本中使用CRLF导致无法执行的问题
- 增加bolt delete API
- API返回值和参数修复，去除 " 和 [ 等符号

### v0.6.0 (2023.1.2)
- 支持[gpython](https://github.com/go-python/gpython)扩展
  - 不完整的Python3.4
  - 增加re正则表达式库
  - 暂未开放设置接口
- 使用[Auto_Bangumi](https://github.com/EstrellaXD/Auto_Bangumi)的[raw_parser.py](https://github.com/EstrellaXD/Auto_Bangumi/blob/main/src/parser/analyser/raw_parser.py)进行解析番剧名
  - 移除poketto依赖
  - 稍微修改以适配gpython
- 修改部分单测，使用单独测试数据
- 更改部分代码结构
- 修复torrent内含有多个文件时，重命名失败的问题
  - TODO: 支持外挂字幕文件的重命名和移动
- 启动时检查bangumi缓存修改时间，大于24小时则执行更新

### v0.5.4
- 初始化或升级配置文件后直接退出

### v0.5.3 (2022.12.21)
- 修复重复打开bolt导致死锁的问题
- 新增查询数据库相关接口
  - GET /api/bolt
  - GET /api/bolt/value

### v0.5.2 (2022.12.21)
- 增加休眠机制
  - 下载器中无正在下载、正在做种或已下载项目时，将不会扫描本地文件

### v0.5.1 (2022.12.21)
- 新增schedule定时任务
  - 每周四固定更新AnimeGoData数据
- 移除从bangumi获取当前ep信息流程
- 移动部分代码位置
- 清理无用models

### v5.0.0 (2022.12.20)
- **配置文件版本号为`1.1.0`**
- 全新的downloader manager
  - 更加准确的判断是否重复下载
  - 移除无用配置项
- 取消对下载器的依赖
  - 以AnimeGo自身数据库为准
  - 权限的rename模块，根据下载状态判断重命名
- bolt中使用json存储
  - 移除gob依赖
- 移除不完全的dmhy支持

### v0.4.3 (2022.11.13)
- 更改webapi插件相关接口中，插件名搜索规则
  - 不需要传递 'plugin' 这一层文件夹
  - 插件名可以忽略'.js'后缀
  - 插件名可以使用上层文件夹名，会自动加载文件夹内部的 'main.js' 或 'plugin.js'
- webapi支持Swagger
- webapi增加配置项值获取、配置项注释获取和yaml配置文件获取接口

### v0.4.2 (2022.11.12)
- 修改代码兼容性
- 简单支持dmhy(未开放)

### v0.4.1 (2022.11.9)
- 支持根据插件名搜索插件文件
  - 插件名可以忽略'.js'后缀
  - 插件名可以使用上层文件夹名，会自动加载文件夹内部的 'main.js' 或 'plugin.js'
- 使用最新bolt分支: bbolt
- 支持bolt cache批量写入数据
  - 使AnimeGoData写入数据库速度大大提升

### v0.4.0 (2022.11.6)
- **配置文件版本号为`1.0.0`**
- 修复下载路径为相对路径时，qbt下载位置错误的问题
- 首次使用自动释放资源
- 配置文件动态创建
- 支持缓存时间自定义
- 彩色日志输出

### 2022.10.30 (v0.3.0)
- 修改配置文件结构
- 使用[gorequest](https://github.com/parnurzeal/gorequest)作为网络请求库
  - 更好的重试等待和超时
- 默认关闭debug模式

### 2022.10.23 (v0.2.3)
- 修改插件读取文件规范，现在仅能够读取所在路径文件
- 请求响应非200不再进行缓存
- 修复webapi参数绑定问题

### 2022.10.23 (v0.2.2)
- 修复错误信息嵌套问题，优化错误提示
- UserAgent

### 2022.10.15
- 修改和增加webapi
  - 支持access_key认证
  - 修改url地址
  - 支持设置和获取插件的json配置文件

### 2022.10.14
- 统一错误提示\[未完成]
- 修正js和bangumi部分bug

### 2022.10.6 alpha-0.1
- qBittorrent关闭重连功能
  - qbt退出期间下载项会暂存于下载队列中，重启后恢复下载
  - 下载队列在qbt客户端关闭期间满的话，会停止解析、停止下载
- 更好的日志分类
  - 一般提示[INFO]
  - 已知错误[WARN]，根据具体错误类型自动判断是否继续解析此项
  - 未知错误[ERROR]，可能会影响功能的正常使用
- tmdb默认值设置

### 2022.10.5
- 支持WebAPI
  - 支持Tampermonkey(油猴)插件 [AnimeGo\[Mikan快速订阅\]](https://greasyfork.org/zh-CN/scripts/449596-animego-mikan%E5%BF%AB%E9%80%9F%E8%AE%A2%E9%98%85) 快速订阅下载
- 整理项目初始化相关的代码结构
- 整理测试文件初始化
- 修复重复下载的问题

### 2022.10.4
- 增加种子大小Length字段
- 增加资源释放程序
- 修改部分配置结构

### 2022.10.2
- 初步的插件模型：内置javascript虚拟机引擎 [goja](https://github.com/dop251/goja) 
- 支持js脚本编写Rss过滤器
  - 支持筛选
  - 支持日志、获取Mikan信息等基础内置函数
- 支持 [poketto](https://github.com/3verness/poketto)初步解析下载项名

### 2022.8.28
- **项目正式更名为AnimeGo**
- 更改项目结构
- 增加filter接口（暂无实现）
- 将feed更新manager移动到filter manager
- 增加部分包注释
- 增加http请求超时重试机制（mikan除外）

### 2022.8.27
- 更改项目结构
- 更改缓存方式
  - 新增memorizer组件
  - 改用函数缓存，入参为key，返回值为value
- anisource使用单例模式，避免不必要的对象创建

### 2022.8.17
- 更改项目结构
  - 将anisource抽离到pkg，可单独使用
- 优化错误处理
- 使用goreq来进行网络请求

### 2022.8.14
- 增加主函数
  - 支持安全退出
- 修复bug
  - tmdb报错问题
  - 自动创建不存在的文件夹
  - 安全退出错误
  - ...

### 2022.8.13
- 优化下载流程，修复重复下载问题
- 更改项目结构

### 2022.8.4
- **重构项目结构**
  - 删除不必要`interface`定义
  - 修改为较规范的go项目结构（未完成）
  - `parser`包不使用结构体，直接使用函数
- 修改关键结构体命名
  - `Bangumi`->`Anime`，和bgm.tv网站作区分
- 统一订阅、下载器的manager结构，均采用协程方式运行
  - 订阅和下载器之间，支持使用chan传递下载项
- 待补充和完善...

### 2022.8.3
- 下载流程管理基本完成
- 完整的mikan rss自动下载基本完成

### 2022.7.31
- 修改项目结构
  - 优化config包结构，删除不必要函数
  - 将config和cache放在store文件夹中，并交由store包统一管理

### 2022.7.26 - 7.30
- 修改部分代码结构、细节，增加部分字段
- 日志使用zap
- cache等放在全局变量store中
- cache bucket命名常量化
- 修正qbittorrent方法
- **新增manager**
  - 支持使用client进行下载
  - 支持下载流程的管理，包括排队、进度获取等
  - 支持重命名、移动正在下载项

### 2022.6.27
- 番剧信息完全缓存，即同一个url、同一个番剧、同一集仅需请求一次
  - 使用gob来序列化与反序列化数据
- 支持高级设置，包括各种数据的缓存过期时间等细致配置

### 2022.6.22
- 增加[boltDB](https://github.com/boltdb/bolt) 作为缓存数据库的调用接口
- 调整models.Bangumi结构，使数据传递更合理
- 支持获取当前ep的信息

### 2022.6.21
- 调整项目结构
- 调整配置文件结构
- 支持设置代理
- 支持标签通配符

### 2022.6.19
- 调整项目结构
- 完善[TheMovieDB](https://www.themoviedb.org/) 信息获取，多次查询
  - 完成番剧别名处理
- 完善部分错误处理

### 2022.6.14 2
- 完成[Bangumi](https://bgm.tv) 信息获取
- 完成[TheMovieDB](https://www.themoviedb.org/) 信息获取
  - 搜索存在缺陷待修复
- 增加process包来调用core下功能

### 2022.6.14
- 完成Mikan Rss信息获取与解析

### 2022.6.13
- 配置文件读取
- qBittorrent客户端api的简单再封装

### 2022.6.12
- 项目框架搭建
