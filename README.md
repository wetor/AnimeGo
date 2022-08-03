# GoBangumi

使用Golang编写的全自动追番工具，简单的部署和使用，方便的模块化扩展

## 目的
- 简化部署和使用，以及模块化扩展
- 学习

## 目前进度
### 因个人原因暂停一段时间开发，目前仍无法完整的运行！第一个可使用版本开发完成后会放出Releases
- 基本配置文件、rss解析、番剧信息获取完成。
- 调用下载器进行下载、Jellyfin支持，以及自动控制等待开发。

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
### 2022.8.3
- 下载流程管理基本完成

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