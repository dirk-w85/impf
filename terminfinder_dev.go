package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"os"
	"log"
	"time"
	"strings"
)

type Termine struct {
	Terminevorhanden            bool          `json:"termineVorhanden"`
	Vorhandeneleistungsmerkmale []interface{} `json:"vorhandeneLeistungsmerkmale"`
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

func ErrorCheck(e error) {
    if e != nil {
		log.Fatalln(e)
        panic(e)
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
	body, err := ioutil.ReadAll(resp.Body)
	ErrorCheck(err)
	if resp.StatusCode < 300 {
		return string(body)
	}else{
		return "{}"
	}
}

func GetDates (resp string, Plz string, Url string) (bool){
	var Dates Termine
	err := json.Unmarshal([]byte(resp), &Dates)
	ErrorCheck(err)
	return Dates.Terminevorhanden
}

func main() {
	//fmt.Println("Terminfinder started. Pulling data...")

	f, err := os.OpenFile("/var/log/terminfinder.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	ErrorCheck(err)
	defer f.Close()

	log.SetOutput(f)
	resp := GetRequest("https://005-iz.impfterminservice.de/assets/static/impfzentren.json", "")
	var Zentren Impfzentrum
	err = json.Unmarshal([]byte(resp), &Zentren)
	ErrorCheck(err)

	log.Printf("%d Impfzentren in Baden-Wuerttemberg\n", len(Zentren.BadenW))
	for z := 0; z < len(Zentren.BadenW); z++  {
		url := Zentren.BadenW[z].URL+"rest/suche/termincheck?plz="+Zentren.BadenW[z].Plz+"&leistungsmerkmale=L920,L921,L922,L923"
		resp := GetRequest(url, "")
		dateAvailable := GetDates(resp,Zentren.BadenW[z].Plz,Zentren.BadenW[z].URL)
		//dateAvailable = true
		if dateAvailable {
			fmt.Printf("terminfinder,ort=%s,plz=%s termin=%d\n", strings.ReplaceAll(Zentren.BadenW[z].Ort, " ", "_"),Zentren.BadenW[z].Plz, 1)
			log.Printf("terminfinder,ort=%s,plz=%s termin=%d\n", strings.ReplaceAll(Zentren.BadenW[z].Ort, " ", "_"),Zentren.BadenW[z].Plz, 1)
		}else{
			fmt.Printf("terminfinder,ort=%s,plz=%s termin=%d\n", strings.ReplaceAll(Zentren.BadenW[z].Ort, " ", "_"),Zentren.BadenW[z].Plz, 0)
		}
		time.Sleep(200*time.Millisecond)		
	}	
}