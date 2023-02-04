package try

const rethrowPanic = "_____rethrow"

type (
	E         interface{}
	Exception struct {
		finally func()
		Error   E
	}
)

func Throw() {
	panic(rethrowPanic)
}

func This(f func()) (e Exception) {
	e = Exception{nil, nil}
	// catch error in
	defer func() {
		e.Error = recover()
	}()
	f()
	return
}

func (e Exception) Catch(f func(err E)) {
	if e.Error != nil {
		defer func() {
			// call finally
			if e.finally != nil {
				e.finally()
			}

			// rethrow exceptions
			if err := recover(); err != nil {
				if err == rethrowPanic {
					err = e.Error
				}
				panic(err)
			}
		}()
		f(e.Error)
	} else if e.finally != nil {
		e.finally()
	}
}

func (e Exception) Finally(f func()) (e2 Exception) {
	if e.finally != nil {
		panic("finally was only set")
	}
	e2 = e
	e2.finally = f
	return
}
