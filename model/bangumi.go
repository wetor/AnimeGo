package model

import "fmt"

type Bangumi struct {
	ID     int // bgm id
	SubID  int // 其他id
	Name   string
	Season int
	Ep     int
}

func (b *Bangumi) FullName() string {
	str := fmt.Sprintf("%s[第%d季][第%d集]", b.Name, b.Season, b.Ep)
	return str
}
