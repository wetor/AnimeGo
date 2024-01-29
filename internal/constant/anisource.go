package constant

// bangumi

const (
	BangumiSubjectBucket         = "bangumi_sub"
	BangumiMinSimilar    float64 = 0.75
	BangumiDefaultHost           = "https://api.bgm.tv"
	BangumiBucket                = "bangumi"
)

// mikan

const (
	MikanIdXPath         = "//a[@class='mikan-rss']"                                 // Mikan番剧id获取XPath
	MikanGroupXPath      = "//p[@class='bangumi-info']/a[@class='magnet-link-wrap']" // Mikan番剧信息获取group字幕组id和name
	MikanBangumiUrlXPath = "//p[@class='bangumi-info']/a[contains(@href, 'bgm.tv')]" // Mikan番剧信息中bangumi id获取XPath

	MikanAuthCookie = ".AspNetCore.Identity.Application"

	MikanDefaultHost = "https://api.bgm.tv"
	MikanBucket      = "bangumi"
)

// themoviedb

const (
	ThemoviedbDefaultHost             = "https://api.themoviedb.org"
	ThemoviedbBucket                  = "themoviedb"
	ThemoviedbMatchSeasonDays         = 90
	ThemoviedbMinSimilar      float64 = 0.75
)
