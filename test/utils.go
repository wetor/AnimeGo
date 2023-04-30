package test

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func LogBatchCompare(r io.Reader, args ...any) {
	index := 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		skip := 0
		line := scanner.Text()
		switch val := args[index].(type) {
		case int:
			skip = val - 1
		case string:
			if !strings.Contains(line, val) {
				panic(fmt.Sprintf(`"%s" not in "%s"`, val, line))
			}
		case func(string) bool:
			if !val(line) {
				panic(line)
			}
		}
		index++
		if index >= len(args) {
			break
		}
		for i := 0; i < skip; i++ {
			scanner.Scan()
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
