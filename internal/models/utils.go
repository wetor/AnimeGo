package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

func Format(format string, p map[string]any) string {
	args, i := make([]string, len(p)*2), 0
	for k, v := range p {
		args[i] = "{" + k + "}"
		args[i+1] = fmt.Sprint(v)
		i += 2
	}
	return strings.NewReplacer(args...).Replace(format)
}

var filenameMap = map[string]any{
	`/`: "",
	`\`: "",
	`[`: "(",
	`]`: ")",
	`:`: "-",
	`;`: "-",
	`=`: "-",
	`,`: "-",
}

func Filename(filename string) string {
	return Format(filename, filenameMap)
}

func Md5Str(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
