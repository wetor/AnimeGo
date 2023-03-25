import log

import sys
sys.path.append("../../../../../assets/plugin/filter")
from Auto_Bangumi.raw_parser import analyse


def filter_all(argv):
    result_index = []
    for i, item in enumerate(argv['items']):
        result = analyse(item['name'])
        log.info(item['name'], result)
        if result.group == 'NC-Raws':
            result_index.append(i)
    return {
        'index': result_index,
        'error': None
    }
