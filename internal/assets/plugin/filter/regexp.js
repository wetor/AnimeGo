// 正则表达式过滤器
// 过滤条目名符合regexp表达式的条目

const regexp = /1080/

function main(argv){
    resultIndex = []
    argv.feedItems.forEach(function(item, index, self){
        if (regexp.test(item.Name)){
            resultIndex.push(index)
        }
    })
    return {
        index: resultIndex,
        error: null
    }
}