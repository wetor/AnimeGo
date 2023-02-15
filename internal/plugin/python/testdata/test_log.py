import log

def main(argv):
    print(argv)
    log.debug(argv)
    log.info(argv)
    log.debugf('test:%d', 11)
