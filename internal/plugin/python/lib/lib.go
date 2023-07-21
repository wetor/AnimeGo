package lib

var isInit = false

func Init() {
	if !isInit {
		InitLog()
		InitAnimeGo()
		isInit = true
	}
}
