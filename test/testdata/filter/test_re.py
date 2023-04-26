import log
import re


def filter_all(argv):
    result_index = []
    for i, item in enumerate(argv['items']):
        log.info(item)
        if re.search('1080', item['name']):
            result_index.append(i)

    return {
        'index': result_index,
        'error': None
    }
