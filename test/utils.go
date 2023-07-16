package test

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func LogBatchCompare(r io.Reader, match func(line, match string) bool, args ...any) {
	if match == nil {
		match = func(line, match string) bool {
			return strings.Contains(line, match)
		}
	}
	index := 0
	lineIndex := 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		skip := 0
		line := scanner.Text()
		switch val := args[index].(type) {
		case int:
			skip = val - 1
		case string:
			if !match(line, val) {
				panic(fmt.Sprintf(`"%s" not match "%s"`, val, line))
			}
		case map[string]int:
			size := 0
			for _, v := range val {
				size += v
			}
			startLineIndex := lineIndex
			endLineIndex := lineIndex + size
			for j := 0; j < size; j++ {
				for k, v := range val {
					if match(line, k) {
						if v > 0 {
							val[k]--
							break
						} else {
							panic(fmt.Sprintf(`"%v" not match in [%d,%d]`, val, startLineIndex, endLineIndex))
						}
					}
				}
				if j < size-1 {
					scanner.Scan()
					lineIndex++
					line = scanner.Text()
				}
			}
			size = 0
			for _, v := range val {
				size += v
			}
			if size > 0 {
				panic(fmt.Sprintf(`"%v" not match in [%d,%d]`, val, startLineIndex, endLineIndex))
			}
		case []string:
			// 判断接下来等长度的日志中，是否包含数组内成员内容
			size := len(val)
			startLineIndex := lineIndex
			endLineIndex := lineIndex + size
			for j := 0; j < size; j++ {
				deleteIndex := -1
				for i, l := range val {
					if match(line, l) {
						deleteIndex = i
						break
					}
				}
				if deleteIndex == -1 {
					panic(fmt.Sprintf(`"%v" not match in [%d,%d]`, val, startLineIndex, endLineIndex))
				} else {
					val = append(val[:deleteIndex], val[deleteIndex+1:]...)
				}
				if j < size-1 {
					scanner.Scan()
					lineIndex++
					line = scanner.Text()
				}
			}
			if len(val) > 0 {
				panic(fmt.Sprintf(`"%v" not match in [%d,%d]`, val, startLineIndex, endLineIndex))
			}
		}
		lineIndex++
		index++
		if index >= len(args) {
			break
		}
		for i := 0; i < skip; i++ {
			scanner.Scan()
			lineIndex++
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
