# Schedule定时任务插件
使用Cron表达式定时执行任务
## 插件设置
### \_\_name\_\_
全局变量`__name__`，str类型，可选  
定时任务的名称，方便日志中区分  

### \_\_cron\_\_
全局变量`__cron__`，str类型，必要  
定时任务的定时规则，支持 [秒] [分] [小时] [日] [月] [周] 六项的Cron表达式    
[Cron表达式文档（维基百科，推荐）](https://zh.wikipedia.org/wiki/Cron)  
[Cron表达式文档（百度百科）](https://baike.baidu.com/item/cron/10952601)  

## 入口函数
### run
`run(args) -> dict`

#### 入参
暂不支持参数传递，`args`固定为`None`，但是无法省略

#### 返回值
暂不支持返回值，即使返回了数据也不会做任何处理  
