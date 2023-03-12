
__name__ = "Refresh_Test"
__name2__ = ""

# 5s 执行一次
__cron__ = "*/5 * * * * ?"


def run(args):
    print(args)
    print(__name__)
    print(__name2__)
