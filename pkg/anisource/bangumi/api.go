package bangumi

import "fmt"

var infoApi = func(id int) string {
	return fmt.Sprintf("%s/v0/subjects/%d", Host, id)
}

var epInfoApi = func(id, ep, eps int) string {
	rang := 2
	if eps < 25 {
		rang = eps / 5
	} else if eps < 50 {
		rang = eps / 8
	} else if eps < 300 {
		rang = eps / 12
	} else {
		rang = 30
	}
	offset := ep - 1 - rang
	if offset < 0 {
		offset = 0
	}
	limit := rang*2 + 1

	epType := 0 // 仅番剧本体
	return fmt.Sprintf("%s/v0/episodes?subject_id=%d&type=%d&limit=%d&offset=%d", Host, id, epType, limit, offset)
}
