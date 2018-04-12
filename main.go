package main

import (
	"flag"
	"fmt"
	. "github.com/qichengzx/sshh/core"
	"os"
	"path/filepath"
)

const (
	appName = "sshh"
	version = "0.1.2"
)

var (
	c string
	h bool
)

func init() {
	flag.BoolVar(&h, "h", false, "This help")
	flag.StringVar(&c, "c", "", "Set configuration `filename` (default ./servers.json)")
	flag.Parse()
	flag.Usage = usage
}

func main() {
	switch {
	case c != "":
		c = parseConfig(c)
		app := New(c)
		app.Start()
		break
	case h:
		flag.Usage()
		break
	default:
		flag.Usage()
		break
	}
}

func parseConfig(c string) string {
	if !filepath.IsAbs(c) {
		appPath, _ := os.Getwd()
		c = filepath.Join(appPath, c)
	}
	return c
}

func usage() {
	fmt.Fprintf(os.Stdout, appName+" "+version+`
Usage: `+appName+` [options]

Options:
`)
	flag.PrintDefaults()
}
