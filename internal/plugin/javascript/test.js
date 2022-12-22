
function main(argv){
    resultIndex = []
    log.info(variable, variable.name, variable.version)
    argv.feedItems.forEach(function(item, index, self){
        try {
            log.info(item.Name,' ',item.Length)
        }catch (err){
            log.error(err,' ',item.Name)
        }
    })

    return {
        index: resultIndex,
        error: null
    }
}