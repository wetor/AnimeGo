
def filter_all(argv):
    result = []
    for i, item in enumerate(argv['items']):
        result.append({
            'index': i
        })

    return {
        'data': result,
        'error': None,
    }
