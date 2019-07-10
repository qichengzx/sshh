package main

import (
	"flag"
	"fmt"
	. "github.com/qichengzx/sshh/core"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

const (
	appName   = "sshh"
	version   = "0.3.0"
	cacheFile = ".sshh_profile"
)

var (
	c string
	h bool
)

func init() {
	flag.BoolVar(&h, "h", false, "This help")
	flag.StringVar(&c, "c", "", "Set configuration `filename` (default ./servers.yaml)")
	flag.Parse()
	flag.Usage = usage
}

func main() {
	switch {
	case c != "":
		c = parseConfig(c)
		profileWrite(c)
		appRun(c)
		break
	case h:
		flag.Usage()
		break
	default:
		f, err := profileRead()
		if err != nil {
			log.Fatal(err)
		}
		appRun(string(f))
		break
	}
}

func appRun(c string) {
	app := New(c)
	app.Start()
}

func parseConfig(c string) string {
	if !filepath.IsAbs(c) {
		appPath, _ := os.Getwd()
		c = filepath.Join(appPath, c)
	}
	return c
}

func home() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return usr.HomeDir
}

func profileWrite(c string) error {
	homeStr := home()
	return ioutil.WriteFile(filepath.Join(homeStr, cacheFile), []byte(c), 0644)
}

func profileRead() ([]byte, error) {
	homeStr := home()
	return ioutil.ReadFile(filepath.Join(homeStr, cacheFile))
}

func usage() {
	fmt.Fprintf(os.Stdout, appName+" "+version+`
Usage: `+appName+` [options]

Options:
`)
	flag.PrintDefaults()
}
