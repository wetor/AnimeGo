import log
import core

class Test:
    a = '111'
    b = 101
    c = {'d': '111', 'e': '222'}

def main(args):
    obj = Test()
    assert core.dumps(obj) == '{"a":"111","b":101,"c":{"d":"111","e":"222"}}'
    assert core.dumps({"ttt": True, "objs": [{"a": 123}, {"vb": "Test"}]}) == \
           '{"objs":[{"a":123},{"vb":"Test"}],"ttt":true}'

    j = core.loads(args.json)
    y = core.loads(args.yaml, 'yaml')

    jo = core.dumps(j)
    log.info(jo)
    yo = core.dumps(y, 'yaml')
    log.info(yo)

    return {
        'json': j,
        'yaml': y
    }
