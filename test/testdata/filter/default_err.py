

def filter_all(argv):
    result_index = []
    for i, item in enumerate(argv['items']):
        result_index.append(i)
    return {
        'index': result_index,
        'error': '错误'
    }
