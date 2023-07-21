import log
import core

def main(args):
    result = core.parse_mikan('https://mikanani.me/Home/Episode/a6f48155e7648a945e9bf85949c6cf8d8eb7ad61')
    log.info(result)
    s = core.dumps(result)
    log.info(s)
