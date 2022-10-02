// getRow
//  @Description: 获取当前js脚本所执行的行数信息，用于日志输出
//
function getRow(){
    try {
        const callstack = new Error().stack.split("\n")
        for(let i = callstack.length-1; i >= 0; i--) {
            let matchArray = callstack[i].match(/at (.+?) \((.+?)\)$/)
            if (matchArray){
                let nameSplit = matchArray[2].split("/")
                return nameSplit[nameSplit.length-1]
            }
        }
    }catch (err) {
        return ""
    }
    return ""
}

log = {
    debug(...params){
        goLog.debug(`[${getRow()}]\t`, ...params)
    },
    info(...params){
        goLog.info(`[${getRow()}]\t`, ...params)
    },
    error(...params){
        goLog.error(`[${getRow()}]\t`, ...params)
    },
}

filter = {
    definition(height){

    }
}