import log
import core
import re
from Auto_Bangumi.raw_parser import analyse


def filter_all(argv):
    result = []
    for i, item in enumerate(argv['items']):
        # 解析标题
        parsed = analyse(item.name)
        # log.info(i, item, parsed)
        # 解析失败
        if not parsed.episode:
            # -----------------
            # 这里进行二次解析处理
            # -----------------
            log.infof('%d %s 「%s」', i, 'ep解析错误', item.name)
            continue
        # 跳过非1080
        if not re.search('1080', parsed.resolution):
            continue
        result.append({
            'index': i,
            'parsed': parsed,
        })
        log.info(core.dumps(result[len(result)-1]))

    return {
        'data': result,
        'error': None,
    }
