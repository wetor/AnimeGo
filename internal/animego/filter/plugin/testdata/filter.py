import log
import raw_parser


def main(argv):
    result_index = []
    for i, item in enumerate(argv['feedItems']):
        result = raw_parser.parser.analyse(item.Name)
        if result.group == 'NC-Raws':
            result_index.append(i)
    return {
        'index': result_index,
        'error': None
    }
