import core


__scrape__ = True

def rename(args):
    anime = args['anime']
    dir = core.filename(anime['name_cn'])
    ext = args['filename'].split(".")[-1]
    if anime['ep_type'] == 0:
        # 无法解析ep时，使用原名保存
        filepath = '%s/S%02d/%s' % (dir, anime['season'], args['filename'])
    else:
        filepath = '%s/S%02d/E%03d.%s' % (dir, anime['season'], anime['ep'], ext)
    return {
        'error': None,
        'filename': filepath,
        'dir': dir
    }
