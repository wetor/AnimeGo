package main

import (
	"flag"
	"github.com/golang/glog"
)

func main() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "log")
	flag.Set("v", "10")
	flag.Parse()
	defer glog.Flush()
	glog.V(5).Infoln("hello world")
}
