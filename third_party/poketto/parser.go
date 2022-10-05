package poketto

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	groupRepl      = strings.NewReplacer("【", "[", "】", "]")
	groupRegx      = regexp.MustCompile(`[\[\]]`)
	splitMatchRegx = regexp.MustCompile(`(.*|\[.*])( -? \d+ |\[\d+]|\[\d+.?[vV]\d{1}]|[第]\d+[话話集]|\[\d+.?END])(.*)`)
	seasonRegx1    = regexp.MustCompile(`新番|月?番`)
	seasonRegx2    = regexp.MustCompile(`.*新番.`)
	seasonRegx3    = regexp.MustCompile(`^[^]】]*[]】]`)
	seasonRegx4    = regexp.MustCompile(`S\d{1,2}|Season \d{1,2}|[第].[季期]`)
	seasonRegx5    = regexp.MustCompile(`S|Season`)
	seasonRegx6    = regexp.MustCompile(`[第 ].*[季期]`)
	seasonRegx7    = regexp.MustCompile(`[第季期 ]`)
	nameRepl       = strings.NewReplacer("（仅限港澳台地区）", "")
	nameRegx1      = regexp.MustCompile(`/|  |-  `)
	nameMatchRegx  = regexp.MustCompile(`([^\x00-\xff]{1,})(\s)([\x00-\xff]{4,})`)
	nameRegx2      = regexp.MustCompile(`[aA-zZ]`)
	epRegx         = regexp.MustCompile(`\d{1,4}`)
	tagRegx1       = regexp.MustCompile(`[\[\]()（）]`)
	tagRegx2       = regexp.MustCompile(`[简繁日字幕]|CH|BIG5|GB`)
	tagRepl        = strings.NewReplacer("_MP4", "")
	tagRegx3       = regexp.MustCompile(`1080|720|2160|4K`)
	tagRegx4       = regexp.MustCompile(`B-Global|[Bb]aha|[Bb]ilibili|AT-X|Web`)
	numDict        = map[rune]int{
		'一': 1, '二': 2, '三': 3, '四': 4, '五': 5, '伍': 5,
		'六': 6, '七': 7, '八': 8, '九': 9, '十': 10,
	}
)

type Episode struct {
	TitleRaw string

	Name       string
	Season     int
	Ep         int
	Group      string
	Definition string
	Sub        string
	Source     string

	ParseErr error
}

func NewEpisode(raw string) *Episode {
	return &Episode{TitleRaw: raw}
}

func (ep *Episode) TryParse() {
	ep.ParseErr = ep.parse()
}

func (ep *Episode) ToMap() map[string]interface{} {
	dict := make(map[string]interface{})
	dict["raw"] = ep.TitleRaw
	if len(ep.Name) > 0 {
		dict["name"] = ep.Name
	}
	if ep.Season > 0 {
		dict["season"] = ep.Season
	}
	if ep.Ep > 0 {
		dict["ep"] = ep.Ep
	}
	if len(ep.Group) > 0 {
		dict["group"] = ep.Group
	}
	if len(ep.Definition) > 0 {
		dict["definition"] = ep.Definition
	}
	if len(ep.Sub) > 0 {
		dict["sub"] = ep.Sub
	}
	if len(ep.Source) > 0 {
		dict["source"] = ep.Source
	}
	return dict
}

func (ep *Episode) ToFields() []string {
	return []string{ep.TitleRaw, ep.Name, fmt.Sprint(ep.Season), fmt.Sprint(ep.Ep), ep.Group, ep.Definition, ep.Sub, ep.Source}
}

func (ep *Episode) parse() error {
	raw := ep.TitleRaw
	if raw == "" {
		return errors.New("原始标题为空，无法解析。")
	}
	raw = groupRepl.Replace(raw)
	var group string
	if groupRegx.MatchString(raw) {
		group = groupRegx.Split(raw, -1)[1]
	}
	matcher := splitMatchRegx.FindStringSubmatch(raw)
	if matcher == nil {
		return CannotParseErr
	}
	name, season, err := getSeason(matcher[1])
	if err != nil {
		return CannotParseSeasonErr
	}
	name, err = getName(name)
	if err != nil {
		return CannotParseNameErr
	}
	epNum, err := getEp(matcher[2])
	if err != nil {
		return CannotParseEpErr
	}
	definition, sub, source, err := getTag(matcher[3])
	if err != nil {
		return CannotParseTagErr
	}

	ep.Name = name
	ep.Season = season
	ep.Ep = epNum
	ep.Group = group
	ep.Definition = definition
	ep.Sub = sub
	ep.Source = source
	return nil
}

func getSeason(raw string) (name string, season int, err error) {
	if seasonRegx1.MatchString(raw) {
		raw = seasonRegx2.ReplaceAllString(raw, "")
	} else {
		raw = seasonRegx3.ReplaceAllString(raw, "")
		raw = strings.TrimSpace(raw)
	}
	raw = groupRegx.ReplaceAllString(raw, "")
	seasonRe := seasonRegx4
	seasonMatcher := seasonRe.FindAllString(raw, -1)
	if seasonMatcher == nil {
		return raw, 1, nil
	} else {
		name = seasonRe.ReplaceAllString(raw, "")
		for _, s := range seasonMatcher {
			if seasonRegx5.MatchString(s) {
				season, err = strconv.Atoi(seasonRegx5.ReplaceAllString(s, ""))
				if err == nil {
					return
				}
			} else if seasonRegx6.MatchString(s) {
				seasonBuf := seasonRegx7.ReplaceAllString(s, "")
				if season, err = strconv.Atoi(seasonBuf); err == nil {
					return
				}
				if season, err = getNum(seasonBuf); err == nil {
					return
				}
			}
		}
	}
	return "", 0, errors.New("无法识别季数")
}

func getNum(raw string) (int, error) {
	for _, r := range []rune(raw) {
		if n, ok := numDict[r]; ok {
			return n, nil
		}
	}
	return 0, errors.New("无法转换为数字")
}

func getName(raw string) (name string, err error) {
	raw = strings.TrimSpace(raw)
	raw = nameRepl.Replace(raw)
	slicesRaw := nameRegx1.Split(raw, -1)
	var slices []string
	for _, s := range slicesRaw {
		if s != "" {
			slices = append(slices, s)
		}
	}
	if len(slices) == 1 {
		if strings.Contains(raw, "_") {
			slices = strings.Split(raw, "_")
		} else if strings.Contains(raw, " - ") {
			slices = strings.Split(raw, "-")
		}
	}
	if len(slices) == 1 {
		matcher := nameMatchRegx.FindStringSubmatch(raw)
		if matcher != nil && matcher[3] != "" {
			return matcher[3], nil
		}
	}
	maxLen := 0
	for _, s := range slices {
		if l := len(nameRegx2.FindAllString(s, -1)); l > maxLen {
			maxLen = l
			name = s
		}
	}
	name = strings.TrimSpace(name)
	return name, nil
}

func getEp(raw string) (epNum int, err error) {
	if epRaw := epRegx.FindString(raw); epRaw != "" {
		epNum, err = strconv.Atoi(epRaw)
		return
	}
	return 0, errors.New("无法解析集数")
}

func getTag(raw string) (dpi, sub, source string, err error) {
	raw = tagRegx1.ReplaceAllString(raw, " ")
	tagsRaw := strings.Split(raw, " ")
	var tags []string
	for _, t := range tagsRaw {
		if t != "" {
			tags = append(tags, t)
		}
	}
	for _, t := range tags {
		if tagRegx2.MatchString(t) {
			sub = tagRepl.Replace(t)
		} else if tagRegx3.MatchString(t) {
			dpi = t
		} else if tagRegx4.MatchString(t) {
			source = t
		}
	}
	err = nil
	return
}
