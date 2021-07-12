package main

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time" 
	"encoding/base64"
	"os"
	"bytes"
	"strconv"
)

type Impftermine struct {
	Gesuchteleistungsmerkmale []string      `json:"gesuchteLeistungsmerkmale"`
	Termine                   []interface{} `json:"termine"`
	Terminetss                []interface{} `json:"termineTSS"`
	Praxen                    struct {} `json:"praxen"`
}
	
type Impfzentrum struct {
	BadenW []struct {
		Zentrumsname string `json:"Zentrumsname"`
		Plz          string `json:"PLZ"`
		Ort          string `json:"Ort"`
		Bundesland   string `json:"Bundesland"`
		URL          string `json:"URL"`
		Adresse      string `json:"Adresse"`
	} `json:"Baden-Württemberg"`
}

type Configuration struct {
	Codes    []string `json:"Codes"`
	Logfile  string   `json:"Logfile"`
	Slackurl string   `json:"SlackURL"`
	FromPLZ  int      `json:"FromPLZ"`
	ToPLZ    int      `json:"ToPLZ"`
}

func ErrorCheck(e error) {
    if e != nil {
		log.Fatalln(e)
        panic(e)
    }
}

func SendSlack (msg string){
	client := &http.Client{}
	//data := '{"text": "Hello, world."}'
	req, err := http.NewRequest("POST", Config.Slackurl, bytes.NewBufferString("{\"text\": \""+msg+"\"}"))
	ErrorCheck(err)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36")
	
	resp, err := client.Do(req)
	ErrorCheck(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	ErrorCheck(err)
	if resp.StatusCode < 300 {
		fmt.Println(string(body))
		log.Println("Slack message send...")
	}else{
		fmt.Println("{}")
	}
}

func GetRequest (url string, token string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	ErrorCheck(err)

	if token != ""{
		req.Header.Set("Authorization", "Basic "+token)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36")
	
	resp, err := client.Do(req)
	ErrorCheck(err)
	defer resp.Body.Close()
	//fmt.Println(resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	ErrorCheck(err)
	if resp.StatusCode < 300 {
		return string(body)
	}else{
		return "{}"
	}
}

func GetDates (resp string, Plz string, Code string, Url string) (){
	var Dates Impftermine
	err := json.Unmarshal([]byte(resp), &Dates)
	ErrorCheck(err)
	//fmt.Println(Dates)
	if len(Dates.Termine) == 0 && len(Dates.Terminetss) == 0 && len(Dates.Gesuchteleistungsmerkmale) >= 1 {
		//fmt.Println("Kein Termin verfügbar... - PLZ: "+Plz+" Code: "+Code)
		//log.Println(Dates)
		//SendSlack("Keine Termine")
	}
	if len(Dates.Termine) >= 1 || len(Dates.Terminetss) >= 1 && len(Dates.Gesuchteleistungsmerkmale) >= 1{
		msg:="TERMIN VERFÜGBAR - PLZ: "+Plz+" Code: "+Code+" Link: "+Url+"impftermine/suche/"+Code+"/"+Plz+"/"
		log.Println(msg)
		log.Println(Dates)
		SendSlack(msg)
	}
	if len(Dates.Termine) >= 1 {
		log.Println(Dates)
	}
	if len(Dates.Terminetss) >= 1 {
		log.Println(Dates)
	}
}

// -----------------
	
func main()  {	
	file, err := ioutil.ReadFile("/lue13/coding/impf/config.json")
	ErrorCheck(err)
	fmt.Println("Reading Config succefull...")
	var Config Configuration
	json.Unmarshal([]byte(file), &Config)

	fmt.Println("Pulling data...")

	f, err := os.OpenFile(Config.Logfile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	ErrorCheck(err)
	defer f.Close()

	log.SetOutput(f)
	log.Println("Impf Checker started.")
	for {
		resp := GetRequest("https://005-iz.impfterminservice.de/assets/static/impfzentren.json", "")
		var Zentren Impfzentrum
		err = json.Unmarshal([]byte(resp), &Zentren)
		ErrorCheck(err)	

		for z := 0; z < len(Zentren.BadenW); z++  {
			currentPlz, _ := strconv.Atoi(Zentren.BadenW[z].Plz)			
			if currentPlz >= Config.FromPLZ && currentPlz <= Config.ToPLZ {
				//fmt.Println(currentPlz)
				url := Zentren.BadenW[z].URL+"rest/suche/impfterminsuche?plz="+Zentren.BadenW[z].Plz
				//fmt.Println(url)
				for g := 0; g < len(Config.Codes); g++  {
					token := base64.StdEncoding.EncodeToString([]byte(":"+Config.Codes[g]))
					resp := GetRequest(url, token)
					GetDates(resp,Zentren.BadenW[z].Plz,Config.Codes[g],Zentren.BadenW[z].URL)
					time.Sleep(200*time.Millisecond)
				}
			}
		}
	}
} 