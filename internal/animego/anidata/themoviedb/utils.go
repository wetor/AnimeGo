package themoviedb

import "time"

func StrTimeSubAbs(t1, t2 string) int {
	if len(t1)*len(t2) == 0 {
		return 0
	}
	s := StrTimeSub(t1, t2)
	if s < 0 {
		return -int(s / (24 * 60 * 60))
	} else {
		return int(s / (24 * 60 * 60))
	}
}

func StrTimeSub(t1, t2 string) int64 {
	time1, _ := time.Parse("2006-01-02", t1)
	time2, _ := time.Parse("2006-01-02", t2)
	ut1 := time1.Unix()
	ut2 := time2.Unix()
	return ut1 - ut2
}
