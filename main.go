package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"
)

type Data struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type JsonData struct {
	Status Data `json:"status"`
}

func main() {
	path := "data.json"

	go updateJson(path)

	http.HandleFunc("/", welcome)
	http.ListenAndServe(":8080", nil)

}

func welcome(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("data.json")

	if err != nil {
		log.Println("Error reading data from json", err)
		return
	}

	var data Data

	var jsonData JsonData

	if err := json.Unmarshal(file, &jsonData); err != nil {
		log.Println("Error unmarshaling data", err)
		return
	}

	// Ambil data dari struktur jsonData dan letakkan di dalam struktur Data yang sesuai
	data = jsonData.Status

	// Parse template
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	var statusWater string
	var statusWind string

	if data.Water < 5 {
		statusWater = "aman"
	} else if data.Water >= 6 && data.Water <= 8 {
		statusWater = "siaga"
	} else {
		statusWater = "bahaya"
	}

	if data.Wind < 6 {
		statusWind = "aman"
	} else if data.Wind >= 7 && data.Wind <= 15 {
		statusWind = "siaga"
	} else {
		statusWind = "bahaya"
	}

	err = tmpl.Execute(w, struct {
		JSONData    Data
		StatusWater string
		StatusWind  string
	}{
		JSONData:    data,
		StatusWater: statusWater,
		StatusWind:  statusWind,
	})
	if err != nil {
		log.Fatal("Error executing template:", err)
	}
}

func updateJson(path string) {

	for {
		file, err := os.ReadFile(path)

		if err != nil {
			log.Println("Error reading data from json", err)
			return
		}

		var data Data

		if err := json.Unmarshal(file, &data); err != nil {
			log.Println("Error unmarshaling data", err)
			return
		}

		data.Water = rand.Intn(101)
		data.Wind = rand.Intn(101)

		updated, err := json.MarshalIndent(
			map[string]Data{
				"status": data,
			}, "", "   ")

		if err != nil {
			log.Println("Error convert to json: ", err)
		}

		if err := os.WriteFile(path, updated, 0644); err != nil {
			log.Println("Error writing to JSON file:", err)
			return
		}
		time.Sleep(15 * time.Second)
	}
}
