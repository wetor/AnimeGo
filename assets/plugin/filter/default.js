// default默认过滤器
// 无任何过滤，返回所有条目

function main(argv){
    resultIndex = []
    argv.feedItems.forEach(function(item, index, self){
        resultIndex.push(index)
    })
    return {
        index: resultIndex,
        error: null
    }
}