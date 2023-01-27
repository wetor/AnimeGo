# 日志包
import log

# 使用 plugin/anisource/Auto_Bangumi/raw_parser.py
import sys
sys.path.append("../anisource/Auto_Bangumi")
import raw_parser


def main(argv):
    result_index = []
    for i, item in enumerate(argv['feedItems']):
        log.info(item)
        result = raw_parser.parser.analyse(item.Name)
        if result.group == 'NC-Raws':
            result_index.append(i)
    return {
        'index': result_index,
        'error': None
    }
