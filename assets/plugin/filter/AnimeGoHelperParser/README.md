关联项目

AnimeGo本体:   [AnimeGo](https://github.com/wetor/AnimeGo)

AnimeGo过滤插件:   [AnimeGoFilterPlugin](https://github.com/deqxj00/AnimeGoFilterPlugin)

Tampermonkey插件:   [AnimeGoHelper\[Mikan快速订阅\]](https://greasyfork.org/zh-CN/scripts/449596) 

AnimeGo过滤插件 和 Tampermonkey插件需要配套使用

------------------------

1.使用前复制到AnimeGo中plugin目录下，并在config下animego.yaml文件中修改插件路径以及名字
```
  filter:
    javascript: plugin/filter/AnimeGoHelperParser/AnimeGoHelperParser.js
```

2.网页筛选规则导出需要配合Tampermonkey插件使用 

3.过滤规则同步到AnimeGo就行

