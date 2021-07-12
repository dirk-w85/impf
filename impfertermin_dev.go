package main

import (
	//"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"  
	"encoding/base64"
)

type Impftermine struct {
	Gesuchteleistungsmerkmale []string      `json:"gesuchteLeistungsmerkmale"`
	Termine                   []interface{} `json:"termine"`
	Terminetss                []interface{} `json:"termineTSS"`
	Praxen                    struct {
	} `json:"praxen"`
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

	req.Header.Set("Authorization", "Basic "+token)
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

func GetDates (resp string, Plz string, Code string) (){
	var Dates Impftermine
	err := json.Unmarshal([]byte(resp), &Dates)
	ErrorCheck(err)
	//fmt.Println(len(Dates.Termine))
	if len(Dates.Termine) == 0 && len(Dates.Terminetss) == 0 {
		log.Println("Kein Termin verfügbar... - PLZ: "+Plz+" Code: "+Code)
	}else{
		log.Println("TERMIN VERFÜGBAR - PLZ: "+Plz+" Code: "+Code)
	}



}

// ------------------

func main()  {

	Plz := []string{"70174","70176","71334","71065","72072","70629"}
	Codes := []string{"MPXB-8CZU-MZM4","LQUK-VT7U-XVD7","CHKK-2XAJ-9WNY"}
	Urls := []string{"001-iz.impfterminservice.de","002-iz.impfterminservice.de","003-iz.impfterminservice.de", "229-iz.impfterminservice.de"}

	for i := 0; i < len(Plz); i++  {
		//fmt.Printf("Checking for PLZ %s\n", Plz[i])

		for g := 0; g < len(Codes); g++  {
			for h := 0; h < len(Urls); h++  {
				//fmt.Printf("Checking for PLZ %s with Code %s\n", Plz[i], Codes[g])
				token := base64.StdEncoding.EncodeToString([]byte(":"+Codes[g]))
				resp := GetRequest("https://"+Urls[h]+"/rest/suche/impfterminsuche?plz="+Plz[i], token)
				GetDates(resp,Plz[i],Codes[g])
				//fmt.Println(resp)
				time.Sleep(2*time.Second)
			}
		}
	}
} 