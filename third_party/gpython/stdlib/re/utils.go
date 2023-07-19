package re

import "github.com/go-python/gpython/py"

func fixStringParam(string1, string2 py.Object, oldName string) (py.Object, error) {
	if string2 != nil {
		if string1 != nil {
			return nil, py.ExceptionNewf(py.TypeError, "Argument given by name ('%s') and position (1)", oldName)
		}
		return string2, nil
	}
	if string1 == nil {
		return nil, py.ExceptionNewf(py.TypeError, "Required argument 'string' (pos 1) not found")
	}
	return string1, nil
}

func getSlice(String, Start, End py.Object) (res py.Object, start int, end int) {
	if Start == nil && End == nil {
		return String, 0, 0
	}
	str := string(String.(py.String))
	length := len(str)
	start = 0
	if Start != nil {
		start = int(Start.(py.Int))
	}
	end = length
	if End != nil {
		end = int(End.(py.Int))
	}
	if start < 0 {
		start = 0
	} else if start > length {
		start = length
	}

	if end < 0 {
		end = 0
	} else if end > length {
		end = length
	}
	return py.String(str[start:end]), start, end
}

func toString(obj py.Object) (str string, isBytes bool) {
	switch t := obj.(type) {
	case py.String:
		isBytes = false
		str = string(t)
	case py.Bytes:
		isBytes = true
		str = string(t)
	}
	return str, isBytes
}
