package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"bytes"
)

func ErrorCheck(e error) {
    if e != nil {
		log.Fatalln(e)
        panic(e)
    }
}

func main() {
	url :="https://hooks.slack.com/services/"
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBufferString("{\"text\": \"TERMIN VERFÜGBAR - PLZ: 73730 Code: CHKK-2XAJ-9WNY Link: https://229-iz.impfterminservice.de/impftermine/suche/CHKK-2XAJ-9WNY/73730/\"}"))
	ErrorCheck(err)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36")
	
	resp, err := client.Do(req)
	ErrorCheck(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	ErrorCheck(err)
	if resp.StatusCode < 300 {
		fmt.Println(string(body))
	}else{
		fmt.Println("{}")
	}
}