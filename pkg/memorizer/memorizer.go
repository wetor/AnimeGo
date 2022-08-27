package memorizer

type Func func(*Params, *Results) error

type Memorizer interface {
	Add(bucket string)
	Put(bucket string, key, val interface{}, ttl int64)
	Get(bucket string, key, val interface{}) error
}

func Memorized(bucket string, mem Memorizer, fn Func) Func {
	mem.Add(bucket)
	return func(input *Params, res *Results) error {
		err := mem.Get(bucket, input.Key(), res)
		if err == nil {
			return nil
		}
		err = fn(input, res)
		if err != nil {
			return err
		}
		mem.Put(bucket, input.Key(), res, input.CacheTTL)
		return nil
	}
}
