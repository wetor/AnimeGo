package test

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

const (
	FlagNone = 1 << iota
	FlagSkipable
	FlagUnknown
)

type MatchFunc func(line, match string) bool

type Range struct {
	Min, Max int
}

func NewRange(min, max int) *Range {
	return &Range{
		Min: min,
		Max: max,
	}
}

func (r *Range) In(count int) bool {
	if count >= r.Min && count <= r.Max {
		return true
	}
	return false
}

func parseFlags(str string) (string, int) {
	if len(str) <= 2 {
		return str, FlagNone
	}

	if str[0] == '#' {
		flag := FlagUnknown
		switch str[1] {
		case '!':
			flag = FlagSkipable
		}
		if flag != FlagUnknown {
			str = str[2:]
		}
		return str, flag
	} else {
		return str, FlagNone
	}
}

func parseSingle(lines []string, index int, pattern string, match MatchFunc) int {
	line := lines[index]
	pattern, flag := parseFlags(pattern)
	result := match(line, pattern)
	switch flag {
	case FlagSkipable:
		// skip
	default:
		if !result {
			panic(fmt.Sprintf(`"%s" not match "%s"`, pattern, line))
		}
	}
	if result {
		index++
	}
	return index
}

func parseMulti(lines []string, index int, pattern []string, match MatchFunc) int {
	// 判断接下来等长度的日志中，是否包含数组内成员内容
	line := lines[index]
	size := len(pattern)

	startLineIndex := index
	endLineIndex := index + size
	for j := 0; j < size; j++ {
		deleteIndex := -1
		for i, l := range pattern {
			if match(line, l) {
				deleteIndex = i
				break
			}
		}
		if deleteIndex == -1 {
			panic(fmt.Sprintf("\"%v\" not match in \n\t%v\n", pattern,
				strings.Join(lines[startLineIndex:endLineIndex], "\n\t")))
		} else {
			pattern = append(pattern[:deleteIndex], pattern[deleteIndex+1:]...)
		}
		if index < len(lines)-1 {
			index++
			line = lines[index]
		}
	}
	if len(pattern) > 0 {
		panic(fmt.Sprintf("\"%v\" not match in \n\t%v\n", pattern,
			strings.Join(lines[startLineIndex:endLineIndex], "\n\t")))

	}
	return index
}

func parseRange(lines []string, index int, pattern map[string]any, match MatchFunc) int {
	line := lines[index]
	size := 0
	countMap := make(map[string]int, len(pattern))
	for k, v := range pattern {
		countMap[k] = 0
		switch val := v.(type) {
		case int:
			size += val
		case *Range:
			size += val.Max
		}
	}

	startLineIndex := index
	for j := 0; j < size; j++ {
		for k := range pattern {
			if match(line, k) {
				countMap[k]++
				break
			}
		}
		if index < len(lines)-1 {
			index++
			line = lines[index]
		}
	}
	endLineIndex := index
	matchedCount := 0
	for k, v := range pattern {
		count := countMap[k]
		matchedCount += count
		switch val := v.(type) {
		case int:
			if count != val {
				panic(fmt.Sprintf("\"%v\" not match in \n\t%v\n", val,
					strings.Join(lines[startLineIndex:endLineIndex], "\n\t")))
			}
		case *Range:
			if !val.In(count) {
				panic(fmt.Sprintf("\"%v(%v)\" not match in \n\t%v\n", k, v,
					strings.Join(lines[startLineIndex:endLineIndex], "\n\t")))
			}
		}
	}
	index = index - (endLineIndex - startLineIndex) + matchedCount
	return index
}

func LogBatchCompare(r io.Reader, match MatchFunc, args ...any) {
	if match == nil {
		match = func(line, match string) bool {
			return strings.Contains(line, match)
		}
	}
	all, _ := io.ReadAll(r)
	all = bytes.ReplaceAll(all, []byte("\r\n"), []byte("\n"))
	lines := strings.Split(string(all), "\n")
	argsIndex := 0
	for lineIndex := 0; lineIndex < len(lines); {
		switch val := args[argsIndex].(type) {
		case int:
			lineIndex += val
		case string:
			lineIndex = parseSingle(lines, lineIndex, val, match)
		case []string:
			lineIndex = parseMulti(lines, lineIndex, val, match)
		case map[string]any:
			lineIndex = parseRange(lines, lineIndex, val, match)
		}
		argsIndex++
		if argsIndex >= len(args) {
			break
		}

	}
}
