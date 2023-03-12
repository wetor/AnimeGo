import core
import log

__name__ = "Mikan_Rss"

__cron__ = ""

__url__ = ""

def parse(args):
    items = core.parse_mikan_rss(args['data'])
    return {
        "items": items,
        "error": None
    }

