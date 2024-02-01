package constant

// bangumi

const (
	BangumiHost                  = "https://api.bgm.tv"
	BangumiBucket                = "bangumi"
	BangumiSubjectBucket         = "bangumi_sub"
	BangumiMinSimilar    float64 = 0.75
)

// mikan

const (
	MikanHost            = "https://mikanani.me"
	MikanBucket          = "mikan"
	MikanAuthCookie      = ".AspNetCore.Identity.Application"                        // Mikan 认证cookie名
	MikanIdXPath         = "//a[@class='mikan-rss']"                                 // Mikan番剧id获取XPath
	MikanGroupXPath      = "//p[@class='bangumi-info']/a[@class='magnet-link-wrap']" // Mikan番剧信息获取group字幕组id和name
	MikanBangumiUrlXPath = "//p[@class='bangumi-info']/a[contains(@href, 'bgm.tv')]" // Mikan番剧信息中bangumi id获取XPath
)

// themoviedb

const (
	ThemoviedbHost                    = "https://api.themoviedb.org"
	ThemoviedbBucket                  = "themoviedb"
	ThemoviedbApiKey                  = "api_key" // Themoviedb api_key参数名
	ThemoviedbMatchSeasonDays         = 90
	ThemoviedbMinSimilar      float64 = 0.75
)
