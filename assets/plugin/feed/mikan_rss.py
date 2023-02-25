import core
import log

__name__ = "Mikan_Rss"

# 每分钟第10秒执行
__cron__ = "10 0/1 * * * ?"

__url__ = "https://mikanani.me/RSS/Bangumi?bangumiId=228&subgroupid=1"

def parse(args):
    log.info(len(args['data']))
    items = core.parse_mikan_rss(args['data'])
    for item in items:
        log.info(item)

    return {
        "items": items,
        "error": None
    }

