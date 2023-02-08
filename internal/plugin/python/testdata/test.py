import log

def main(args):
    log.info(args['params'])
    return {
        'result': len(args['params'])
    }


def test(args):
    log.info("test")
    log.info(args)
