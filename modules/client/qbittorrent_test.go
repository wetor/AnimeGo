package client

import (
	"GoBangumi/config"
	"GoBangumi/models"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"testing"
)

var qbt Client

func TestMain(m *testing.M) {
	fmt.Println("begin")
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "log")
	flag.Set("v", "10")
	flag.Parse()
	defer glog.Flush()

	config.Init("../../data/config/conf.yaml")

	conf := config.ClientQBt()
	qbt = NewQBittorrent(conf.Url, conf.Username, conf.Password)
	m.Run()
	fmt.Println("end")
}

func Test_QBittorrent(t *testing.T) {
	list := qbt.List(&models.ClientListOptions{
		Status:   QBtStatusAll,
		Category: "test",
	})
	fmt.Println(len(list))
	for _, i := range list {
		fmt.Println(i.Name, i.Hash)
	}

}
func TestQBittorrent_Add(t *testing.T) {
	qbt.Add(&models.ClientAddOptions{
		Urls: []string{
			"https://mikanani.me/Download/20220612/4407d51f30f6033513cbe56cae0120881b0a7406.torrent",
			"https://mikanani.me/Download/20220611/56e13c0c4788b77782722ee46d3c6f27233f676b.torrent",
		},
		SavePath:    "/srv/dev-disk-by-uuid-317ec4d4-2933-4fa6-9d7b-fcdfc339de04/share/downloads/complete/",
		Category:    "test",
		Tag:         "test_tag",
		SeedingTime: 60,
	})

}

func TestQBittorrent_Delete(t *testing.T) {

	qbt.Delete(&models.ClientDeleteOptions{
		Hash: []string{"4407d51f30f6033513cbe56cae0120881b0a7406"},
	})

}

func TestQBittorrent_Get(t *testing.T) {

	qbt.Get(&models.ClientGetOptions{Hash: "4407d51f30f6033513cbe56cae0120881b0a7406"})

}
