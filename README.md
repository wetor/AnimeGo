# GoBangumi

使用Golang编写的全自动追番工具，相较于[AutoBangumi](https://github.com/EstrellaXD/Auto_Bangumi) 简化部署和使用，以及更方便的模块化扩展

## 目的
- 简化部署和使用，以及模块化扩展
- 学习

## 开发计划
- [ ] 完整的[AutoBangumi](https://github.com/EstrellaXD/Auto_Bangumi) 追番功能
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