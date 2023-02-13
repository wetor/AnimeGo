import time

import log
import core
import pylib


def get_is_push(value, key, title):
    log.infof('| key: %v', key)
    log.debugf('| is_enable_whitelist: %v', value.is_enable_whitelist)
    log.debugf('| is_enable_blacklist: %v', value.is_enable_blacklist)

    is_push = True
    is_whitelist_has_word = False
    is_blacklist_has_word = False

    log.debug('| whitelist')
    if value.is_enable_whitelist:
        log.debugf('| %v', value.whitelist)
        for item in value.whitelist:
            if title.find(item) >= 0:
                is_whitelist_has_word = True
                break

    log.debug('| blacklist')
    if value.is_enable_blacklist:
        log.debugf('| %v', value.blacklist)
        for item in value.blacklist:
            if title.find(item) >= 0:
                is_blacklist_has_word = True
                break

    # 白名单
    if value.is_enable_whitelist and not value.is_enable_blacklist:
        if is_whitelist_has_word:
            is_push = True
        else:
            is_push = False

    # 黑名单
    if not value.is_enable_whitelist and value.is_enable_blacklist:
        if is_blacklist_has_word:
            is_push = False
        else:
            is_push = True

    # 白名单+和名单
    if value.is_enable_whitelist and value.is_enable_blacklist:
        if is_whitelist_has_word and not is_blacklist_has_word:
            is_push = True
        else:
            is_push = False

    log.debugf('| is_whitelist_has_word: %v', is_whitelist_has_word)
    log.debugf('| is_blacklist_has_word: %v', is_blacklist_has_word)
    log.infof('| is_push: %v', is_push)
    return is_push


def main(args):
    result = []
    myFiliters = _get_config()
    isNeedGetMikanInfo = False

    log.debugf('==================================')
    log.debugf('| Filiter0 %v', len(myFiliters.Filiter0))
    log.debugf('| Filiter1 %v', len(myFiliters.Filiter1))
    log.debugf('| Filiter2 %v', len(myFiliters.Filiter2))
    log.debugf('| Filiter3 %v', len(myFiliters.Filiter3))
    log.debugf('| Filiter4 %v', len(myFiliters.Filiter4))

    if len(myFiliters.Filiter1) > 0 or len(myFiliters.Filiter2) > 0 or len(myFiliters.Filiter3) > 0:
        isNeedGetMikanInfo = True

    for index, item in enumerate(args.feedItems):
        try:
            parsed = pylib.analyse(item)
            log.debug(' ')
            log.debug('==================================')
            log.info('| '+item.Name, ' ', item.Length)
            log.debug('==================================')
            isPush0 = True
            isPush1 = True
            isPush2 = True
            isPush3 = True
            isPush4 = True
            # 0
            log.infof('| isNeedGetMikanInfo: %v', isNeedGetMikanInfo)
            log.debug('==========0.Gobal.Start===========')
            for key, filters in myFiliters.Filiter0.items():
                isPush0 = get_is_push(filters, key, item.Name)

            # 1,2,3
            log.debug('==1,2,3.MikanID,SubGroupID.Start==')
            if isNeedGetMikanInfo:
                mikanInfo = core.parse_mikan(item.Url)
                key1 = 'key_' + str(int(mikanInfo.id)) + '_' + str(int(mikanInfo.sub_group_id))
                key2 = str(int(mikanInfo.id))
                key3 = str(int(mikanInfo.sub_group_id))
                log.debug('| key1: '+key1)
                log.debug('| key2: '+key2)
                log.debug('| key3: '+key3)
                if key1 in myFiliters.Filiter1.keys():
                    isPush1 = get_is_push(myFiliters.Filiter1.get(key1), key1, item.Name)
                elif key2 in myFiliters.Filiter2.keys():
                    isPush2 = get_is_push(myFiliters.Filiter2.get(key2), key2, item.Name)
                elif key3 in myFiliters.Filiter3.keys():
                    isPush3 = get_is_push(myFiliters.Filiter3.get(key2), key3, item.Name)
                else:
                    log.debug('| no fetch')
            # 4
            log.debug('========4.GroupName.Start=========')
            log.debugf('| %v', parsed.group)
            key4 = parsed.group
            if key4 in myFiliters.Filiter4.keys():
                isPush4 = get_is_push(myFiliters.Filiter4.get(key4), key4, item.Name)
            log.debug('==============Allend==============')
            if isPush0 and isPush1 and isPush2 and isPush3 and isPush4:
                log.infof('| push index: %v', index)
                result.append({
                    'index': index,
                    'parsed': parsed,
                })
                log.infof('| pushed index: %v', index)
            else:
                log.infof('| drop index: %v, %v', index, item.Name)

            log.debug('==================================')
            log.debug(' ')

            if isNeedGetMikanInfo:
                time.sleep(1)

        except Exception as e:
            log.debugf('%v %v', item.Name, e)

    return {
        'data': result,
        'error': None
    }
