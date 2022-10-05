package poketto

import (
	"AnimeGo/internal/animego/parser"
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"testing"
)

var rawList = ``

func TestParser(t *testing.T) {
	fr, err := os.Open("./data/test_data.txt")
	if err != nil {
		panic(err)
	}
	defer fr.Close()
	sc := bufio.NewScanner(fr)

	var eps []*Episode
	for sc.Scan() {
		title := sc.Text()
		cmp, err := parser.ParseTitle(title)
		ep := NewEpisode(title)
		ep.TryParse()
		if ep.ParseErr == nil && err == nil {
			fmt.Printf("ep: [%v %v], dpi: [%v %v], sub: [%v %v], source: [%v %v], err: [%v %v]\n",
				ep.Ep, cmp.Ep, ep.Definition, cmp.Definition, ep.Sub, cmp.Subtitle, ep.Source, cmp.Source, ep.ParseErr, err)
		} else {
			fmt.Println(ep.ParseErr, err)
		}
		eps = append(eps, ep)
	}

	fw, err := os.Create("./data/test_out.txt")
	if err != nil {
		panic(err)
	}
	defer fw.Close()
	w := csv.NewWriter(fw)
	for _, ep := range eps {
		if err := w.Write(ep.ToFields()); err != nil {
			panic(err)
		}
	}
	w.Flush()
}
