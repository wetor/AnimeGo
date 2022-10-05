
function main(argv){
    resultIndex = []

    argv.feedItems.forEach(function(item, index, self){
        try {
            let result = animeGo.parseName(item.Name)
            log.info(item.Name,' ',item.Length)
            // let mikanInfo = animeGo.getMikanInfo(item.Url)
            // log.info(mikanInfo)
            if('Definition' in result && result.Definition.indexOf('1080')>=0 &&
                item.Length > 100*1024*1024){
                resultIndex.push(index)
            }
            // sleep(1000)
        }catch (err){
            log.error(err,' ',item.Name)
        }
    })

    return {
        index: resultIndex,
        error: null
    }
}