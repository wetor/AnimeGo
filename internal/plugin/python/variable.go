package python

type Variable struct {
	Name     string
	Nullable bool
	Getter   func(name string) interface{}
	Setter   func(name string, val interface{})
}
