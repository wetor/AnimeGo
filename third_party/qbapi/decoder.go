package qbapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Decoder func([]byte, interface{}) error

var (
	JsonDec = json.Unmarshal
	StrDec  = strDec
	IntDec  = intDec
)

func strDec(data []byte, v interface{}) error {
	st, ok := v.(*string)
	if !ok {
		return fmt.Errorf("should use string to decode")
	}
	*st = string(data)
	return nil
}

func intDec(data []byte, v interface{}) error {
	st, ok := v.(*int)
	if !ok {
		return fmt.Errorf("should use int to decode")
	}
	strData := string(data)
	rs, err := strconv.ParseInt(strData, 10, 64)
	if err != nil {
		return err
	}
	*st = int(rs)
	return nil
}
