package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"encoding/json"
)

type Configuration struct {
	Codes    []string `json:"Codes"`
	Logfile  string   `json:"Logfile"`
	Slackurl string   `json:"SlackURL"`
}

func ErrorCheck(e error) {
    if e != nil {
		log.Fatalln(e)
        panic(e)
    }
}

func main()  {
	fmt.Println("Config Test...")

	file, err := ioutil.ReadFile("config.json")
	ErrorCheck(err)
	var Config Configuration
	json.Unmarshal([]byte(file), &Config)

	fmt.Println(Config.Logfile)
}