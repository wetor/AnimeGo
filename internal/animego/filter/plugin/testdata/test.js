function main(argv){
    resultIndex = []
    argv.feedItems.forEach(function(item, index, self){
        if(item.Name.indexOf('1080') >= 0){
            resultIndex.push(index)
        }
    })
    return {
        index: resultIndex,
        error: null
    }
}