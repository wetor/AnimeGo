package torrent

import (
	"bytes"
	"io"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/pkg/exceptions"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/request"
	"github.com/wetor/AnimeGo/pkg/utils"
)

const (
	PaddingFilePrefix = "_____padding_file"

	TypeUnknown = "unknown"
	TypeTorrent = "torrent"
	TypeMagnet  = "magnet"
	TypeFile    = "file"
)

var (
	torrentUriRegx1 = regexp.MustCompile(`^http.*\.torrent.*$`)
	torrentUriRegx2 = regexp.MustCompile(`^http.*([a-fA-F0-9]{40})\.torrent`)
	magnetUriRegx1  = regexp.MustCompile("^magnet:")
	magnetUriRegx2  = regexp.MustCompile("urn:btih:([a-fA-F0-9]{40})")
)

type File struct {
	Name   string `json:"name"`
	Dir    string `json:"dir"`
	Length int64  `json:"length"`
}

func (f File) Path() string {
	return path.Join(f.Dir, f.Name)
}

type Torrent struct {
	Type   string  `json:"type"`
	Url    string  `json:"url"`
	Hash   string  `json:"hash"`
	Name   string  `json:"name"`
	Length int64   `json:"length"`
	Files  []*File `json:"files"`
}

func LoadTorrent(r io.Reader) (*Torrent, error) {
	m, err := metainfo.Load(r)
	if err != nil {
		log.DebugErr(err)
		return nil, errors.WithStack(&exceptions.ErrTorrentLoad{Message: "加载metainfo"})
	}
	info, err := m.UnmarshalInfo()
	if err != nil {
		log.DebugErr(err)
		return nil, errors.WithStack(&exceptions.ErrTorrentLoad{Message: "加载bencode info"})
	}

	infoFiles := info.UpvertedFiles()
	infoName := info.BestName()
	files := make([]*File, 0, len(infoFiles))
	for _, file := range infoFiles {
		filePath := file.DisplayPath(&info)
		if strings.HasPrefix(filePath, PaddingFilePrefix) {
			continue
		}
		dir, name := path.Split(filePath)
		if len(infoFiles) > 1 {
			dir = path.Join(infoName, dir)
		}
		files = append(files, &File{
			Name:   name,
			Dir:    dir,
			Length: file.Length,
		})

	}
	t := &Torrent{
		Type:   TypeTorrent,
		Hash:   m.HashInfoBytes().HexString(),
		Name:   infoName,
		Length: info.TotalLength(),
		Files:  files,
	}
	return t, nil
}

func LoadMagnetUri(uri string) (*Torrent, error) {
	m, err := metainfo.ParseMagnetUri(uri)
	if err != nil {
		log.DebugErr(err)
		return nil, errors.WithStack(&exceptions.ErrTorrentLoad{Message: "加载Magnet metainfo"})
	}
	t := &Torrent{
		Type:   TypeTorrent,
		Url:    uri,
		Hash:   m.InfoHash.HexString(),
		Name:   m.DisplayName,
		Length: 0,
		Files:  nil,
	}
	return t, nil
}

func LoadUri(uri string) (t *Torrent, err error) {
	hash, uriType := tryParseHash(uri)
	// 已存在缓存文件，直接读取文件
	file := path.Join(TempPath, hash+".torrent")
	if len(hash) == 40 && utils.IsExist(file) {
		f, err := os.Open(file)
		if err != nil {
			log.DebugErr(err)
			log.Debugf("载入缓存文件失败")
			// 重新加载
			goto load
		}
		t, err = LoadTorrent(f)
		if err != nil {
			log.Debugf("%s", err)
			// 重新加载
			goto load
		}
		err = f.Close()
		if err != nil {
			log.DebugErr(err)
			log.Debugf("关闭缓存文件失败")
			// 继续执行
		}
		t.Type = TypeFile
		t.Url = file
		return t, nil
	}
load:
	switch uriType {
	case TypeMagnet:
		t, err = LoadMagnetUri(uri)
		if err != nil {
			log.Debugf("%s", err)
			// 继续执行
		}
		t.Url = uri
	case TypeTorrent:
		w := bytes.NewBuffer(nil)
		err = request.GetWriter(uri, w)
		if err != nil {
			log.DebugErr(err)
			return nil, errors.WithStack(&exceptions.ErrRequest{Name: uri})
		}
		tw := bytes.NewBuffer(w.Bytes())
		t, err = LoadTorrent(tw)
		if err != nil {
			return nil, err
		}
		file = path.Join(TempPath, t.Hash+".torrent")
		err = os.WriteFile(file, w.Bytes(), 0666)
		if err != nil {
			log.DebugErr(err)
			// 继续执行
		}
		t.Type = TypeFile
		t.Url = file
	default:
		return nil, errors.WithStack(&exceptions.ErrTorrentUrl{Url: uri})
	}
	return t, nil
}

func tryParseHash(uri string) (string, string) {
	hash, isMagnet := parseMagnetUriHash(uri)
	if isMagnet {
		return hash, TypeMagnet
	}
	hash, isTorrent := parseTorrentUriHash(uri)
	if isTorrent {
		return hash, TypeTorrent
	}
	return "", TypeUnknown
}

func parseTorrentUriHash(link string) (hash string, isTorrent bool) {
	if !torrentUriRegx1.MatchString(link) {
		return
	}
	isTorrent = true
	hashMatches := torrentUriRegx2.FindStringSubmatch(link)
	if len(hashMatches) > 1 {
		hash = hashMatches[1]
	}
	return
}

func parseMagnetUriHash(link string) (hash string, isMagnet bool) {
	if !magnetUriRegx1.MatchString(link) {
		return
	}
	isMagnet = true
	hashMatches := magnetUriRegx2.FindStringSubmatch(link)
	if len(hashMatches) > 1 {
		hash = hashMatches[1]
	}
	return
}
