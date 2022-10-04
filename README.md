# AnimeGo

使用Golang编写的全自动追番工具，简单的部署和使用，方便的模块化扩展

## 目的
- 简化部署和使用，以及模块化扩展
- 学习

## 目前进度
- 可使用配置、筛选和下载等基本功能
- javascript编写筛选器，[帮助文档](internal/animego/plugin/javascript/README.md)
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