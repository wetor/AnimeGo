package gpython

import (
	"github.com/brahma-adshonor/gohook"
	. "github.com/go-python/gpython/py"
)

func Hook() {
	err := gohook.Hook(IntNew, HookIntNew, nil)
	if err != nil {
		panic(err)
	}
}

func HookIntNew(metatype *Type, args Tuple, kwargs StringDict) (Object, error) {
	var xObj Object = Int(0)
	var baseObj Object
	base := 10
	err := ParseTupleAndKeywords(args, kwargs, "|OO:int", []string{"x", "base"}, &xObj, &baseObj)
	if err != nil {
		return nil, err
	}
	if baseObj != nil {
		base, err = MakeGoInt(baseObj)
		if err != nil {
			return nil, err
		}
		if base != 0 && (base < 2 || base > 36) {
			return nil, ExceptionNewf(ValueError, "int() base must be >= 2 and <= 36")
		}
	}
	// Special case converting string types
	switch x := xObj.(type) {
	// FIXME Bytearray
	case Bytes:
		return IntFromString(string(x), base)
	case String:
		return IntFromString(string(x), base)
	}
	if baseObj != nil {
		return nil, ExceptionNewf(TypeError, "int() can't convert non-string with explicit base")
	}
	return MakeInt(xObj)
}
