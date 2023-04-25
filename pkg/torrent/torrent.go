package torrent

import (
	"bytes"
	"io"
	"path"
	"strings"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/wetor/AnimeGo/pkg/request"
)

const (
	PaddingFilePrefix = "_____padding_file"

	TypeTorrent = "torrent"
	TypeMagnet  = "magnet"
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
		return nil, err
	}
	info, err := m.UnmarshalInfo()
	if err != nil {
		return nil, err
	}

	infoFiles := info.UpvertedFiles()
	files := make([]*File, 0, len(infoFiles))
	for _, file := range infoFiles {
		filePath := file.DisplayPath(&info)
		if strings.HasPrefix(filePath, PaddingFilePrefix) {
			continue
		}
		dir, name := path.Split(filePath)
		files = append(files, &File{
			Name:   name,
			Dir:    dir,
			Length: file.Length,
		})

	}
	t := &Torrent{
		Type:   TypeTorrent,
		Hash:   m.HashInfoBytes().HexString(),
		Name:   info.BestName(),
		Length: info.TotalLength(),
		Files:  files,
	}
	return t, nil
}

func LoadMagnetUri(uri string) (*Torrent, error) {
	m, err := metainfo.ParseMagnetUri(uri)
	if err != nil {
		return nil, err
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
	if strings.HasPrefix(uri, TypeMagnet) {
		t, err = LoadMagnetUri(uri)
	} else if strings.HasPrefix(uri, "http") {
		w := bytes.NewBuffer(nil)
		err = request.GetWriter(uri, w)
		if err != nil {
			return nil, err
		}
		t, err = LoadTorrent(w)
	}
	t.Url = uri
	return
}
