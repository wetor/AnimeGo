
def main(argv):
    result = []
    for i, item in enumerate(argv['feedItems']):
        result.append({
            'index': i
        })

    return {
        'data': result,
        'error': None,
    }
