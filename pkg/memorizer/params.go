package memorizer

type Params struct {
	CacheTTL int64
	Keys     []string
	Values   []interface{}
}

func NewParams(params ...interface{}) *Params {
	if len(params)%2 != 0 {
		panic(`params格式为: ["value_name1", value1, ]...`)
	}
	p := &Params{
		Keys:   make([]string, 0, 4),
		Values: make([]interface{}, 0, 4),
	}
	for i := 0; i < len(params); i += 2 {
		p.Keys = append(p.Keys, params[i].(string))
		p.Values = append(p.Values, params[i+1])
	}
	return p
}

func (p *Params) TTL(ttl int64) *Params {
	p.CacheTTL = ttl
	return p
}

func (p Params) Get(paramName string) interface{} {
	for i, key := range p.Keys {
		if key == paramName {
			return p.Values[i]
		}
	}
	return nil
}

func (p Params) Key() interface{} {
	return p.Values
}
