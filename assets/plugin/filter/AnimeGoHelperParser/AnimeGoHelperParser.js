//https://stackoverflow.com/questions/29085197/how-do-you-json-stringify-an-es6-map
function replacer(key, value) {
    if (value instanceof Map) {
        return {
            dataType: 'Map',
            value: Array.from(value.entries()), // or with spread: value: [...value]
        };
    } else {
        return value;
    }
}
function reviver(key, value) {
    if (typeof value === 'object' && value !== null) {
        if (value.dataType === 'Map') {
            return new Map(value.value);
        }
    }
    return value;
}

function getIsPush(value, key,titlename) {
    log.info('| key:' + key)
    log.debug('| is_enable_whitelist:' + value.is_enable_whitelist)
    log.debug('| is_enable_blacklist:' + value.is_enable_blacklist)
    var is_push = true
    var is_whitelist_has_word = false
    var is_blacklist_has_word = false
    log.debug('| whitelist')
    if (value.is_enable_whitelist) {
        log.debug('| '+value.whitelist)
        value.whitelist.forEach(function (arr_item, arr_index, arr) {
            if (titlename.includes(arr_item.value)) {
                is_whitelist_has_word = true;
            }
        })
    }
    log.debug('| blacklist')
    if (value.is_enable_blacklist) {
        log.debug(value.blacklist)
        value.blacklist.forEach(function (arr_item, arr_index, arr) {
            if (titlename.includes(arr_item.value)) {
                is_blacklist_has_word = true;
            }
        })
    }
    //白名单
    if (value.is_enable_whitelist && !value.is_enable_blacklist) {
        if (is_whitelist_has_word) {
            is_push = true
        } else {
            is_push = false
        }
    }
    //黑名单
    if (!value.is_enable_whitelist && value.is_enable_blacklist) {
        if (is_blacklist_has_word) {
            is_push = false
        } else {
            is_push = true
        }
    }
    //白名单+黑名单
    if (value.is_enable_whitelist && value.is_enable_blacklist) {
        if (is_whitelist_has_word && !is_blacklist_has_word) {
            is_push = true
        } else {
            is_push = false
        }
    }
    log.debug('| is_whitelist_has_word:' + is_whitelist_has_word)
    log.debug('| is_blacklist_has_word:' + is_blacklist_has_word)
    log.info('| is_push:' + is_push)
    return is_push
}

function main(argv) {
    resultIndex = []

    var jsonstr = os.readFile(variable.name + '.json')
    var myFiliters = JSON.parse(jsonstr, reviver);
    var isNeedGetMikanInfo = false;
    log.debug('==================================')
    log.debug('| Filiter0 ' + myFiliters.Filiter0.size)
    log.debug('| Filiter1 ' + myFiliters.Filiter1.size)
    log.debug('| Filiter2 ' + myFiliters.Filiter2.size)
    log.debug('| Filiter3 ' + myFiliters.Filiter3.size)
    log.debug('| Filiter4 ' + myFiliters.Filiter4.size)

    if (myFiliters.Filiter1.size > 0 || myFiliters.Filiter2.size > 0 || myFiliters.Filiter3.size > 0) {
        isNeedGetMikanInfo = true
    }

    argv.feedItems.forEach(function (item, index, self) {
        try {
            let result = animeGo.parseName(item.Name)
            log.debug(' ')
            log.debug(' ')
            log.debug('==================================')
            log.info('| '+item.Name, ' ', item.Length)
	        log.debug('==================================')
            var isPush0 = true;
            var isPush1 = true;
            var isPush2 = true;
            var isPush3 = true;
            var isPush4 = true;
            //0
            log.info( '| isNeedGetMikanInfo:'+isNeedGetMikanInfo)
            log.debug('==========0.Gobal.Start===========')
            myFiliters.Filiter0.forEach(function (value, key, map) {
                isPush0 = getIsPush(value, key , item.Name)
            })
            //1,2,3
            log.debug('==1,2,3.MikanID,SubGroupID.Start==')
            if (isNeedGetMikanInfo) {
                let mikanInfo = animeGo.getMikanInfo(item.Url)
                //log.info(mikanInfo)
                // ID         int
                // SubGroupID int
                // PubGroupID int
                // GroupName  string
                var key1 = mikanInfo.ID + '+' + mikanInfo.SubGroupID
                var key2 = mikanInfo.ID
                var key3 = mikanInfo.SubGroupID
                log.debug('| key1:'+key1)
                log.debug('| key2:'+key2)
                log.debug('| key3:'+key3)
                if (myFiliters.Filiter1.has(key1)) {
                    isPush1 = getIsPush(myFiliters.Filiter1.get(key1), key1, item.Name)
                } else if (myFiliters.Filiter2.has(key2)) {
                    isPush2 = getIsPush(myFiliters.Filiter2.get(key2), key2, item.Name)
                } else if (myFiliters.Filiter3.has(key3)) {
                    isPush3 = getIsPush(myFiliters.Filiter3.get(key3), key3, item.Name)
                }else{
                    log.debug('| no fetch')
                }
            }
            //4
            log.debug('========4.GroupName.Start=========')
            log.debug('| '+result.Group)
            var key4 = result.Group
            if (myFiliters.Filiter4.has(key4)) {
                isPush4 = getIsPush(myFiliters.Filiter4.get(key4), key4, item.Name)
            }

            log.debug('==============Allend==============')
            if (isPush0 && isPush1 && isPush2 && isPush3 && isPush4) {
                log.info('| push index:' + index)
                resultIndex.push(index)
                log.info('| pushed resultIndex:' + resultIndex)
            }else{
                log.info('| drop index:' + index+','+item.Name)
            }
            log.debug('==================================')
            log.debug(' ')
            log.debug(' ')

            //  if('Definition' in result && result.Definition.indexOf('720')>=0 ||
            //     'Group' in result && result.Group === 'ANi'){
            //resultIndex.push(index)
            //}
            
            if(isNeedGetMikanInfo){
                sleep(1000)
            }
        } catch (err) {
            log.error(err, ' ', item.Name)
        }
    })

    return {
        index: resultIndex,
        error: null
    }
}
