
def filter_all(argv):
    result = []
    for i, item in enumerate(argv['items']):
        result.append(i)

    return {
        'index': result,
        'error': None,
    }
