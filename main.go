package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gen2brain/beeep"
)

type BandaCalendarResponse struct {
	CalendarInfo []CalendarInfo `json:"calendarInfo"`
}

type CalendarInfo struct {
	Date 			string `json:"date"`
	IsAvailable 	int `json:"isAvailable"`
	WaitingListFull 	int `json:"waitingListFull"`
	IsSpecialDay 		bool `json:"isSpecialDay"`
	Type string `json:"type"`
}

type BandaHoursResponse struct {
	Hours Hour `json:"center"`
}

type Hour struct {
	Type string `json:"type"` // results | no-results
	Slots []Slot `json:"slots"`
}

type Slot struct {
	Hour string `json:"hour"`
	Type string `json:"type"` // deposit | no
}

const searchDateUrl = "https://api.meitre.com/api/calendar-availability-new/1620/2/dinner"

func main() {
	res := fetchData[BandaCalendarResponse](searchDateUrl)

	isAvailable := false
	for i := 0; i < len(res.CalendarInfo); i++ {
		if res.CalendarInfo[i].IsAvailable == 1 || res.CalendarInfo[i].Type== "available" {
			isAvailable = true
			date := strings.Split(res.CalendarInfo[i].Date, "T")[0]

			log.Printf("-> Hay una mesa disponible el %v! \n", date)
			log.Println("-> Buscando horarios disponibles...")
			searchHour(date)
			log.Println()
		}
	}
	if !isAvailable {
		log.Println("-> No hay mesas disponibles")
	} else {
		execNotify()
	}
}

func searchHour(time string) {
	url := "https://api.meitre.com/api/search-all-hours/en/2/"+time+"/dinner/1620"
	res := fetchData[BandaHoursResponse](url)

	if res.Hours.Type == "results" && len(res.Hours.Slots) > 0 {
		for i := 0; i < len(res.Hours.Slots); i++ {
			log.Printf("-> Hay un horario disponible a las %v \n", res.Hours.Slots[i].Hour)
		}
	} else {
		log.Println("-> No hay horarios disponibles")
	}
}

func fetchData[T any](url string) T {
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	// Convert the body to type string
	var responseObject T
	json.Unmarshal(body, &responseObject)

	return responseObject
}

func execNotify() {
	icon := "/Users/fedetoledo/Downloads/b.jpg"
	err := beeep.Notify("Mesa disponible!", "Hay una mesa disponible en Banda! :)", icon)
	if err != nil {
		panic(err)
	}
}
