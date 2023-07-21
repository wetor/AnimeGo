package py

import (
	"fmt"
	"github.com/go-python/gpython/py"
	"regexp"
	"strconv"
	"strings"
)

// var formatRegx = regexp.MustCompile(`\{((\d+)|(\w*((\.\w+)+|(\[\w+])+)*)*)((:[^}]*)?)}`)
var formatRegx = regexp.MustCompile(`\{((\d+)|(\w*((\.\w+)+|(\[\w+])+)*)*)((:(.?[<^>]\d+)?[^}]*)?)}`)

func init() {

	py.StringType.Dict["join"] = py.MustNewMethod("join", func(self py.Object, args py.Object) (py.Object, error) {
		argList := args.(*py.List)
		list := make([]string, argList.Len())
		for i, item := range argList.Items {
			list[i] = string(item.(py.String))
		}
		return py.String(strings.Join(list, string(self.(py.String)))), nil
	}, 0, `join(list)`)

	py.StringType.Dict["format"] = py.MustNewMethod("format", format, 0, `format(args...)`)

}

func format(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	var argString string
	var err error
	formatStr := string(self.(py.String))
	matches := formatRegx.FindAllStringSubmatch(formatStr, -1)
	result := formatStr
	argsIndex := 0
	for _, match := range matches {
		fullMatch := match[0]
		indexStr := match[2]
		fieldName := match[3]
		formatSpec := match[7]
		alignSpec := match[9]

		var arg py.Object

		if indexStr != "" {
			// 正则匹配为整数，无需处理错误
			index, _ := strconv.ParseInt(indexStr, 10, 64)
			arg, err = args.M__getitem__(py.Int(index))
			if err != nil {
				return nil, err
			}
		} else if fieldName != "" {
			field := strings.NewReplacer(".", " .", "][", " ", "[", " ", "]", " ").Replace(fieldName)
			parts := strings.Split(field, " ")
			argName := parts[0]
			// 存在参数名
			if argName != "" {
				if argName[0] == '.' {
					argName = argName[1:]
				}
				// 参数名为数字，即索引
				if i, e := strconv.ParseInt(argName, 10, 64); e == nil {
					arg, err = args.M__getitem__(py.Int(i))
				} else {
					arg, err = kwargs.M__getitem__(py.String(argName))
				}
				if err != nil {
					return nil, err
				}
			} else {
				// 不存在参数名使用顺序索引
				arg, err = args.M__getitem__(py.Int(argsIndex))
				if err != nil {
					return nil, err
				}
				argsIndex++
			}
			// 子元素访问
			for _, part := range parts[1:] {
				if part == "" {
					continue
				}
				attrFlag := false
				if part[0] == '.' {
					part = part[1:]
					attrFlag = true
				}
				var keyObj py.Object = py.String(part)
				if i, e := strconv.ParseInt(part, 10, 64); e == nil {
					keyObj = py.Int(i)
				}
				arg, err = getItemOrAttr(arg, keyObj, attrFlag)
				if err != nil {
					return nil, err
				}
			}
		} else {
			// 顺序索引
			arg, err = args.M__getitem__(py.Int(argsIndex))
			if err != nil {
				return nil, err
			}
			argsIndex++
		}

		switch val := arg.(type) {
		case py.Int:
			argString, err = formatInt(val, formatSpec, alignSpec)
		case py.Float:
			argString, err = formatFloat(val, formatSpec, alignSpec)
		case py.String:
			argString, err = formatString(val, formatSpec, alignSpec)
		default:
			str, ok := arg.(py.I__str__)
			if !ok {
				return nil, py.ExceptionNewf(py.TypeError, "'%s' object has no attribute '__str__'", arg.Type().Name)
			}
			strObj, _ := str.M__str__()
			argString = string(strObj.(py.String))
		}
		if err != nil {
			return nil, err
		}
		result = strings.Replace(result, fullMatch, argString, 1)
	}

	return py.String(result), nil
}

func getItemOrAttr(obj py.Object, key py.Object, attr bool) (py.Object, error) {
	if attr {
		if k, ok := key.(py.Int); ok {
			return nil, py.ExceptionNewf(py.AttributeError, "'%s' has no attribute '%d'", obj.Type().Name, int64(k))
		}
		return py.GetAttr(obj, key)
	}
	getter, ok := obj.(py.I__getitem__)
	if !ok {
		return nil, py.ExceptionNewf(py.TypeError, "'%s' object has no attribute '__getitem__'", obj.Type().Name)
	}
	result, err := getter.M__getitem__(key)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func formatString(val py.String, spec, alignSpec string) (string, error) {
	value := string(val)
	if spec == "" {
		return fmt.Sprintf("%s", value), nil
	}
	formatSpec := "%" + strings.TrimLeft(spec[1:], alignSpec)
	end := spec[len(spec)-1]
	if end >= '0' && end <= '9' {
		formatSpec += "s"
	}
	result := fmt.Sprintf(formatSpec, value)
	if strings.HasPrefix(result, "%!"+string(end)) {
		return "", py.ExceptionNewf(py.ValueError, "Unknown format code '%c' for object of type 'str'", end)
	}
	return formatAlign(result, alignSpec), nil
}

func formatInt(val py.Int, spec, alignSpec string) (string, error) {
	value := int64(val)
	if spec == "" {
		return fmt.Sprintf("%d", value), nil
	}
	formatSpec := "%" + strings.TrimLeft(spec[1:], alignSpec)
	end := spec[len(spec)-1]
	if end >= '0' && end <= '9' {
		formatSpec += "d"
	}
	result := fmt.Sprintf(formatSpec, value)
	if strings.HasPrefix(result, "%!"+string(end)) {
		return "", py.ExceptionNewf(py.ValueError, "Unknown format code '%c' for object of type 'int'", end)
	}
	if formatSpec[len(formatSpec)-1] == 'o' && formatSpec[1] == '#' {
		result = "0o" + result[1:]
	}
	return formatAlign(result, alignSpec), nil
}

func formatFloat(val py.Float, spec, alignSpec string) (string, error) {
	value := float64(val)
	if spec == "" {
		return strconv.FormatFloat(value, 'f', -1, 64), nil
	}
	formatSpec := "%" + strings.TrimLeft(spec[1:], alignSpec)
	end := spec[len(spec)-1]
	if end >= '0' && end <= '9' {
		formatSpec += "f"
	}
	result := fmt.Sprintf(formatSpec, value)
	if strings.HasPrefix(result, "%!"+string(end)) {
		return "", py.ExceptionNewf(py.ValueError, "Unknown format code '%c' for object of type 'float'", end)
	}
	return formatAlign(result, alignSpec), nil
}

func formatAlign(result, alignSpec string) string {
	if alignSpec == "" {
		return result
	}
	char := " "
	var align uint8 = '>'
	var sizeStr string
	if alignSpec[0] != '<' && alignSpec[0] != '^' && alignSpec[0] != '>' {
		char = string(alignSpec[0])
		align = alignSpec[1]
		sizeStr = alignSpec[2:]
	} else {
		align = alignSpec[0]
		sizeStr = alignSpec[1:]
	}
	size, _ := strconv.ParseInt(sizeStr, 10, 64)
	padding := int(size) - len(result)
	if align == '^' {
		// 计算左侧填充空格的数量
		leftPad := padding / 2
		rightPad := padding - leftPad
		leftFill := strings.Repeat(char, leftPad)
		rightFill := strings.Repeat(char, rightPad)
		return leftFill + result + rightFill
	}
	fill := strings.Repeat(char, padding)
	if align == '<' {
		return result + fill
	} else if align == '>' {
		return fill + result
	}
	return result
}
