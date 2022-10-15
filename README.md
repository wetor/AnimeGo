# AnimeGo

使用Golang编写的全自动追番工具，简单的部署和使用，方便的模块化扩展

## 使用帮助
```text
-config string
    配置文件路径；配置文件中的相对路径均是相对与程序的位置 (default "data/config/animego.yaml")
-debug
    Debug模式，将会输出更多的日志 (default true)
-init-path string
    [初始化]输出资源/配置文件到的根目录
-init-replace
    [初始化]输出资源/配置文件时是否自动替换
```
### 0.安装插件
- AnimeGo过滤器插件：https://github.com/deqxj00/AnimeGoFilterPlugin
- AnimeGo网页快速订阅Tampermonkey(油猴)插件：https://github.com/deqxj00/AnimeGoHelper

### 1.释放资源
```shell
AnimeGo -init-path=./data
```
> 可选`-init-replace`，启用后遇到已存在文件将不提示直接覆盖，慎用
### 2.修改配置
打开并编辑`./data/config/animego.yaml`
>路径和`1.释放资源`所释放位置有关

### 3.启动程序
```shell
AnimeGo -config=./data/config/animego.yaml
```
> 可选`-debug`，启用后将输出更详细的日志

可以使用如screen等工具放至后台执行，也可以创建服务并启动
### 更多待补充...

## 文档
1. [配置文件](assets/config/animego.yaml)
2. [插件函数文档（仅过滤器）](internal/animego/plugin/javascript/README.md)
3. [webapi接口文档](internal/web/README.md)

## 目的
- 简化部署和使用，以及模块化扩展
- 学习

## 目前进度
- 可使用配置、筛选和下载等基本功能
- javascript编写筛选器，[帮助文档](internal/animego/plugin/javascript/README.md)
- 支持Tampermonkey(油猴)插件 [AnimeGo\[Mikan快速订阅\]](https://greasyfork.org/zh-CN/scripts/449596-animego-mikan%E5%BF%AB%E9%80%9F%E8%AE%A2%E9%98%85) (v0.3.3)快速订阅下载
- Jellyfin支持

## 开发计划
- [ ] 类似[AutoBangumi](https://github.com/EstrellaXD/Auto_Bangumi) 的追番功能
  - [ ] [Mikan Project](https://mikanani.me) 订阅支持
  - [ ] [qBittorrent](https://qbittorrent.org) 等下载支持
  - [ ] [Jellyfin](https://jellyfin.org/) 等媒体库软件识别
  - [ ] ...
- [ ] Web界面支持
- [ ] 模块化与高级自定义功能支持
  - [ ] 独立的订阅支持
  - [ ] 独立下载控制
  - [ ] ...

## 开发日志

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