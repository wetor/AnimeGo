package torrent_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/wetor/AnimeGo/pkg/torrent"
)

func TestLoadMagnet(t *testing.T) {

	ts, _ := torrent.LoadMagnetUri("magnet:?xt=urn:btih:f6aa232b3024073c90d04614fcbf050d94fe8ad6")
	fmt.Println(ts)
	ts, _ = torrent.LoadMagnetUri("magnet:?xt=urn:btih:62VCGKZQEQDTZEGQIYKPZPYFBWKP5CWW&dn=&tr=http%3A%2F%2F104.143.10.186%3A8000%2Fannounce&tr=udp%3A%2F%2F104.143.10.186%3A8000%2Fannounce&tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&tr=http%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce&tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce&tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce&tr=http%3A%2F%2Ftracker.prq.to%2Fannounce&tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce&tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud&tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce&tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce&tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce&tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce&tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce&tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce&tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce")
	fmt.Println(ts)
}

func TestLoadTorrent(t *testing.T) {

	s, _ := os.Open("testdata/40b003ab90b1f7145abeec15e636b901e317a572.torrent")
	ts, _ := torrent.LoadTorrent(s)
	_ = s.Close()
	fmt.Println(ts)
	for _, i := range ts.Files {
		fmt.Println(i)
	}

	m, _ := os.Open("testdata/b7f570888e8967744b399361429ede46d1c0e484.torrent")
	ts, _ = torrent.LoadTorrent(m)
	_ = m.Close()
	fmt.Println(ts)
	for _, i := range ts.Files {
		fmt.Println(i)
	}
}
