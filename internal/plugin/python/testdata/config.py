import log


def test(args):
    log.info(args)
    log.info(__plugin_name__)
    log.info(__plugin_path__)
    log.info(__animego_version__)
    config = _get_config()
    log.info(config)
    for k, v in config.Filiter1.items():
        log.info(k)
        for kk, vv in v.items():
            log.infof('%v = %v', kk, vv)

    return config
