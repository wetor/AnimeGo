
def rename(args):
    anime = args['anime']
    ext = args['src'].split(".")[-1]
    if anime['name_cn'] != '':
        name = anime['name_cn']
    elif anime['name'] != '':
        name = anime['name']
    else:
        name = str(anime['id'])
    dst = '%s/S%02d/E%d.%s' % (name, anime['season'], anime['ep'], ext)
    return {
        'error': None,
        'dst': dst
    }
