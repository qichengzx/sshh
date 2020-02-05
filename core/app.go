package core

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type App struct {
	ConfigPath string
	servers    []Server
}

func (app *App) Start() {
	fmt.Println("Your servers")
	fmt.Printf("[ID]\t%s\t%s\n", "server name", "IP")
	for i, server := range app.servers {
		fmt.Printf("[%d]\t%s\t%s\n", i, server.Name, server.IP)
	}

	server := app.readInput()
	fmt.Println("你选择了: " + server.IP)
	fmt.Println("Connecting...")
	server.Key = app.keyFile(server.Key)
	server.Connect()
}

func New(c, g string) *App {
	servers, _ := readConf(c, g)

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

func readConf(c, g string) ([]Server, error) {
	b, err := ioutil.ReadFile(c)
	if err != nil {
		panic(err)
	}

	conf := Conf{}
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	servers := []Server{}
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
		s.Name = name
		for _, ip := range server.IP {
			s.IP = ip
			servers = append(servers, s)
		}
	}
	for _, server := range conf.Single {
		servers = append(servers, server)
	}

	return servers, nil
}
