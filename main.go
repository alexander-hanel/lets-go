package main

import (
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

const staticConfigPath = "./configs/config.yaml"

const helloWorld = `package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
`

type config struct {
	ProjectPath string `yaml:"project_path"`
	MainName    string `yaml:"source_name"`
	IdeName     string `yaml:"ide"`
}

func readConfigs(configPath string) config {
	yfile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal("ERROR: Could Not Find Config: ", configPath, " ", err)
	}

	var c config
	err2 := yaml.Unmarshal(yfile, &c)
	if err2 != nil {
		log.Fatal("ERROR: yaml.Unmarshal:", err2)
	}
	
	return c
}

func writeMain(MainSourcePath string) {
	f, err := os.Create(MainSourcePath)
	if err != nil {
		log.Fatal("ERROR: File Creation Failed:", err)
	}

	defer f.Close()
	_, err2 := f.WriteString(helloWorld)
	if err2 != nil {
		log.Fatal("ERROR: Write Hello World Failed:", err2)
	}
}

func initModule(dirName string, srcPath string) {
	cmd := exec.Command("go", "mod", "init", srcPath)
	cmd.Dir = dirName
	err := cmd.Run()
	if err != nil {
		log.Fatal("ERROR: Go mod init failed:", err)
	}
}

func openIDE(MainSourcePath string, ideName string) {
	cmd := exec.Command(ideName, ".")
	cmd.Dir = MainSourcePath
	err := cmd.Run()
	if err != nil {
		log.Fatal("ERROR: Opening IDE Failed. Check $PATH", err)
	}
}

func main() {
	gopath := os.Getenv("GOROOT")
	if gopath == "" {
		gopath = build.Default.GOROOT
	}

	if gopath == "" {
		log.Fatal(`ERROR: Environmental variable "GOROOT" is not set`)
		return
	}

	var configPath string
	var projectName string

	flag.StringVar(&configPath, "config", "", "Configuration file for LetsGo")
	flag.StringVar(&projectName, "name", "", "Project name")
	flag.Parse()

	if configPath == "" {
		configPath = staticConfigPath
	}

	// if the user does not pass a project name, create a directory name based on the time
	if projectName == "" {
		now := time.Now()
		projectName = fmt.Sprintf("%d%02d%02d%02d%02d%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	}

	config := readConfigs(configPath)
	dirPath := filepath.Join(config.ProjectPath, projectName)
	if dirPath == "" {
		// TODO add GOPATH as dirpath

	}

	err := os.Mkdir(dirPath, 0755)
	if err != nil {
		log.Fatal("ERROR: Directory Could Not Be Created:", err)
	}

	MainSourcePath := filepath.Join(dirPath, config.MainName)
	writeMain(MainSourcePath)
	initModule(dirPath, projectName)
	fmt.Println("Project Dir Created: ", dirPath)
	openIDE(dirPath, config.IdeName)
}
ope