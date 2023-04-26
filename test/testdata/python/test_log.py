import log

def main(argv):
    print(argv)
    log.debug(argv)
    log.info(int('07'))
    log.debugf('test:%d', 11)
