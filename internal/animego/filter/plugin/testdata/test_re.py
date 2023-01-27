import log
import re


def main(argv):
    result_index = []
    for i, item in enumerate(argv['feedItems']):
        log.info(item)
        if re.search('1080', item.Name):
            result_index.append(i)

    return {
        'index': result_index,
        'error': None
    }
