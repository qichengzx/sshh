package core

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type App struct {
	ConfigPath string
	servers    []Server
}

func (app *App) Start(direct bool) {
	fmt.Println("Your servers")
	fmt.Printf("[ID]\t%s\t%s\n", "server name", "IP")
	for i, server := range app.servers {
		fmt.Printf("[%d]\t%s\t%s\n", i, server.Name, server.IP)
	}

	server := Server{}
	if direct && len(app.servers) > 0 {
		server = app.servers[0]
	} else {
		server = app.readInput()
	}

	fmt.Println("你选择了: " + server.IP)
	fmt.Println("Connecting...")
	server.Key = app.keyFile(server.Key)
	server.Connect()
}

func New(c, g, f string) *App {
	servers := readConf(c, g, f)

	return &App{
		ConfigPath: c,
		servers:    servers,
	}
}

func (app *App) readInput() Server {
	fmt.Print("输入序号: ")

	input := ""
	fmt.Scanln(&input)
	if input == "q" {
		os.Exit(0)
	}

	num, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("输入有误，请重新输入")
		return app.readInput()
	}

	if num > len(app.servers)-1 {
		fmt.Println("输入有误，请重新输入")
		return app.readInput()
	}

	return app.servers[num]
}

func (app *App) keyFile(k string) string {
	if !filepath.IsAbs(k) {
		appPath, _ := os.Getwd()

		path := filepath.Join(appPath, k)
		if _, err := os.Stat(path); err == nil {
			return path
		}

		path = filepath.Join(filepath.Dir(app.ConfigPath), k)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return k
}

type Conf struct {
	//Common is all the common info,
	//So you can no longer fill in a detailed list.
	Common struct {
		User     string `yaml:user`
		Password string `yaml:"password"`
		Port     int    `yaml:"port"`
		Method   string `yaml:"method"`
	}

	Groups Groups   `yaml:"group"`
	Single []Server `yaml:"single"`
}

type Groups struct {
	Groups map[string]Group `yaml:"groups"`
}

type Group struct {
	Name     string
	Method   string   `yaml:"method"`
	User     string   `yaml:"user"`
	Port     int      `yaml:"port"`
	Password string   `yaml:"password"`
	IP       []string `yaml:"ip"`
	Key      string   `yaml:"key"`
}

func readConf(c, g, f string) []Server {
	b, err := ioutil.ReadFile(c)
	if err != nil {
		panic(err)
	}

	conf := Conf{}
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var servers []Server
	for name, server := range conf.Groups.Groups {
		//group check
		if g != "" && name != g {
			continue
		}

		s := Server{
			Name:     name,
			Method:   server.Method,
			User:     server.User,
			Port:     server.Port,
			Password: server.Password,
			Key:      server.Key,
		}

		for _, ip := range server.IP {
			//IP check
			if f != "" && !strings.Contains(ip, f) {
				continue
			}
			s.IP = ip
			servers = append(servers, s)
		}
	}
	for _, server := range conf.Single {
		//IP check
		if f != "" && !strings.Contains(server.IP, f) {
			continue
		}
		servers = append(servers, server)
	}

	for i := 0; i < len(servers); i++ {
		if servers[i].User == "" && conf.Common.User != "" {
			servers[i].User = conf.Common.User
		}
		if servers[i].Port == 0 && conf.Common.Port != 0 {
			servers[i].Port = conf.Common.Port
		}
		if servers[i].Password == "" && conf.Common.Password != "" {
			servers[i].Password = conf.Common.Password
		}
		if servers[i].Method == "" && conf.Common.Method != "" {
			servers[i].Method = conf.Common.Method
		}
	}

	return servers
}
