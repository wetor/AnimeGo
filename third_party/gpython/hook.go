package gpython

import (
	"github.com/brahma-adshonor/gohook"
	"github.com/go-python/gpython/py"
	"math/big"
	"strconv"
	"strings"
)

func Hook() {
	err := gohook.Hook(py.IntFromString, IntFromString, SrcIntFromString)
	if err != nil {
		panic(err)
	}
}

func IntFromString(str string, base int) (py.Object, error) {
	var x *big.Int
	var ok bool
	s := str
	negative := false
	convertBase := base

	// Get rid of padding
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		goto error
	}

	// Get rid of sign
	if s[0] == '+' || s[0] == '-' {
		if s[0] == '-' {
			negative = true
		}
		s = s[1:]
		if len(s) == 0 {
			goto error
		}
	}

	// Get rid of leading sigils and set convertBase
	if len(s) > 1 && s[0] == '0' {
		switch s[1] {
		case 'x', 'X':
			convertBase = 16
		case 'o', 'O':
			convertBase = 8
		case 'b', 'B':
			convertBase = 2
		default:
			goto nosigil
		}
		if base != 0 && base != convertBase {
			// int("0xFF", 10)
			// int("0b", 16)
			convertBase = base // ignore sigil
			goto nosigil
		}
		s = s[2:]
		if len(s) == 0 {
			goto error
		}
	nosigil:
	}
	if convertBase == 0 {
		convertBase = 10
	}

	// Use int64 conversion for short strings since 12**36 < IntMax
	// and 10**18 < IntMax
	if len(s) <= 12 || (convertBase <= 10 && len(s) <= 18) {
		i, err := strconv.ParseInt(s, convertBase, 64)
		if err != nil {
			goto error
		}
		if negative {
			i = -i
		}
		return py.Int(i), nil
	}

	// The base argument must be 0 or a value from 2 through
	// 36. If the base is 0, the string prefix determines the
	// actual conversion base. A prefix of “0x” or “0X” selects
	// base 16; the “0” prefix selects base 8, and a “0b” or “0B”
	// prefix selects base 2. Otherwise the selected base is 10.
	x, ok = new(big.Int).SetString(s, convertBase)
	if !ok {
		goto error
	}
	if negative {
		x.Neg(x)
	}
	return (*py.BigInt)(x).MaybeInt(), nil
error:
	return nil, py.ExceptionNewf(py.ValueError, "invalid literal for int() with base %d: '%s'", convertBase, str)
}

func SrcIntFromString(str string, base int) (py.Object, error) {
	return py.IntFromString(str, base)
}
