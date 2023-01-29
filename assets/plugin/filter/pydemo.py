import log
import re
import pylib


def main(argv):
    result = []
    for i, item in enumerate(argv['feedItems']):
        # 解析标题
        parsed = pylib.analyse(item)
        # log.info(i, item, parsed)
        # 解析失败
        if not parsed.episode:
            # -----------------
            # 这里进行二次解析处理
            # -----------------
            log.infof('%d %s 「%s」', i, 'ep解析错误', item.Name)
            continue
        # 跳过非1080
        if not re.search('1080', parsed.resolution):
            continue
        result.append({
            'index': i,
            'parsed': parsed,
        })

    return {
        'data': result,
        'error': None,
    }
