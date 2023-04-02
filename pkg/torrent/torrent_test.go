package torrent_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/wetor/AnimeGo/pkg/json"
	"github.com/wetor/AnimeGo/pkg/torrent"
)

func TestLoadTorrent(t *testing.T) {
	file := "C:\\Users\\wetor\\Downloads\\b88ba4c86278b05ddfa025b94b0c65597ff72968.torrent"
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	ts, err := torrent.LoadTorrent(f)
	if err != nil {
		panic(err)
	}

	d, _ := json.Marshal(ts)

	fmt.Println(string(d))

}

func TestLoadMagnet(t *testing.T) {

	ts, _ := torrent.LoadMagnetUri("magnet:?xt=urn:btih:f6aa232b3024073c90d04614fcbf050d94fe8ad6")
	fmt.Println(ts)
	ts, _ = torrent.LoadMagnetUri("magnet:?xt=urn:btih:62VCGKZQEQDTZEGQIYKPZPYFBWKP5CWW&dn=&tr=http%3A%2F%2F104.143.10.186%3A8000%2Fannounce&tr=udp%3A%2F%2F104.143.10.186%3A8000%2Fannounce&tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&tr=http%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce&tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce&tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce&tr=http%3A%2F%2Ftracker.prq.to%2Fannounce&tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce&tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud&tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce&tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce&tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce&tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce&tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce&tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce&tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce")
	fmt.Println(ts)
}
