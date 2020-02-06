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
	version   = "0.4.1"
	cacheFile = ".sshh_profile"
)

var (
	c string
	h bool
	g string
	f string
)

func init() {
	flag.BoolVar(&h, "h", false, "This help")
	flag.StringVar(&c, "c", "", "Use specified config file (default ./servers.yaml)")
	flag.StringVar(&g, "g", "", "Only show specificd group(s) in config file")
	flag.StringVar(&f, "f", "", "Only show servers that ip matched the given words")
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
	app := New(c, g, f)
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
Usage: `+appName+` [-c path/to/your/config-file.yaml] [-g group-name] [-f ip]

Options:
`)
	flag.PrintDefaults()
}
