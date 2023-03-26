import core


__write_tvshow__ = True

def rename(args):
    anime = args['anime']
    ext = args['filename'].split(".")[-1]
    if anime['name_cn'] != '':
        name = anime['name_cn']
    elif anime['name'] != '':
        name = anime['name']
    else:
        name = str(anime['id'])

    name = core.filename(name)
    filepath = '%s/S%02d/E%d.%s' % (name, anime['season'], anime['ep'], ext)
    return {
        'error': None,
        'filepath': filepath,
        'tvshow_dir': name
    }
