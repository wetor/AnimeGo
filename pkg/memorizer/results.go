package memorizer

type Results struct {
	*Params
}

func NewResults(params ...interface{}) *Results {
	return &Results{
		Params: NewParams(params...),
	}
}

func (r *Results) Set(params ...interface{}) {
	r.Params = NewParams(params...)
}

func (r *Results) Add(params ...interface{}) {
	if len(params)%2 != 0 {
		panic(`params格式为: ["value_name1", value1, ]...`)
	}
	for i := 0; i < len(params); i += 2 {
		r.Keys = append(r.Keys, params[i].(string))
		r.Values = append(r.Values, params[i+1])
	}
}
