import log

import sys
sys.path.append("../../../../../assets/plugin/filter")
from Auto_Bangumi.raw_parser import analyse


def main(argv):
    result_index = []
    for i, item in enumerate(argv['feedItems']):
        result = analyse(item.Name)
        log.info(item.Name, result)
        if result.group == 'NC-Raws':
            result_index.append(i)
    return {
        'index': result_index,
        'error': None
    }
