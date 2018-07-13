package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	fmt.Println("你选择了: " + server.Name)
	server.Key = app.keyFile(server.Key)
	server.Connect()
}

func New(c string) *App {
	servers, _ := readConf(c)

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

func readConf(c string) ([]Server, error) {
	b, err := ioutil.ReadFile(c)
	if err != nil {
		panic(err)
	}

	servers := []Server{}
	err = json.Unmarshal(b, &servers)
	if err != nil {
		return servers, err
	}

	return servers, nil
}
