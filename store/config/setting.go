package config

//type Setting struct {
//	models.Setting
//}
//
//func (s *Setting) Tag(info *models.Bangumi) string {
//	if len(s.TagSrc) == 0 || info == nil {
//		return ""
//	}
//	date, _ := time.Parse("2006-01-02", info.AirDate)
//	mouth := (int(date.Month()) + 2) / 3
//	str := utils.Format(s.TagSrc, utils.TagFormat{
//		"year":          date.Year(),
//		"quarter":       (mouth-1)*3 + 1,
//		"quarter_index": mouth,
//		"quarter_name":  []string{"冬", "春", "夏", "秋"}[mouth-1],
//		"ep":            info.Ep,
//		"week":          (int(date.Weekday())+6)%7 + 1,
//		"week_name":     []string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}[date.Weekday()],
//	})
//	return str
//}
