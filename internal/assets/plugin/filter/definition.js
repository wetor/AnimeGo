// 分辨率过滤器
// 过滤分辨率tag中含有definition的条目
// 需要是能被解析的标准文件名

const definition = '1080'

function main(argv){
    resultIndex = []
    argv.feedItems.forEach(function(item, index, self){
        try {
            let result = animeGo.parseName(item.Name)
            if('Definition' in result && result.Definition.indexOf(definition)>=0){
                resultIndex.push(index)
            }
        }catch (err){
            log.error(err,' ',item.Name)
        }
    })
    return {
        index: resultIndex,
        error: null
    }
}