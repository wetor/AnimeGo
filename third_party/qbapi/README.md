qbapi
===

golang 版本的qbittorrent webui api, 实现API v2.8.3的大部分功能(除RSS, Search)

## 基本使用

```golang
func main() {
	var opts []Option
	opts = append(opts, WithAuth("xxxtest", "ruok123"))
	opts = append(opts, WithHost("https://torrent.abc.com"))
	api, err := NewAPI(opts...)
	if err != nil {
		panic(err)
	}
	if err := api.Login(context.Background()); err != nil {
		panic(err)
	}
    rsp, err := api.GetApplicationVersion(context.Background(), &GetApplicationVersionReq{})
    if err != nil {
        panic(err)
    }
    fmt.Printf("rsp:%+v", rsp)
}

```
