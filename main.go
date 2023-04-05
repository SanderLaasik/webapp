package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var baseUrl = "http://localhost:8080"

type Airport struct {
	Key     string `json:"_key"`
	Name    string
	City    string
	State   string
	Country string
}

type Flight struct {
	FlightNum     int
	From          Airport
	To            Airport
	DepTimeUTC    string
	ArrTimeUTC    string
	DurFormatted  string
	TailNum       string
	UniqueCarrier string
}

type Connection struct {
	FlightsCount int
	TotFormatted string
	Flights      []Flight
}

func main() {

	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/connections", connectionsHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))

}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	var resBody = makeRequest(http.MethodGet, "/airports", nil)
	var data []Airport
	json.Unmarshal(resBody, &data)

	var tmplFile = "airports.html"
	tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}

}

func connectionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		params := make(map[string]string)
		params["from"] = r.FormValue("from")
		params["to"] = r.FormValue("to")
		date, err := time.Parse("2006-01-02", r.FormValue("date"))
		if err != nil {
			panic(err)
		}
		params["date"] = date.Format("2006/01/02")
		params["limit"] = "5"

		var resBody = makeRequest(http.MethodGet, "/connections", params)
		var data []Connection
		json.Unmarshal(resBody, &data)
		var tmplFile = "connections.html"
		tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
		if err != nil {
			panic(err)
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
	}

}

func makeRequest(method string, path string, params map[string]string) []byte {
	requestURL := fmt.Sprintf(baseUrl + path)
	req, err := http.NewRequest(method, requestURL, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}
	if params != nil {
		q := req.URL.Query()
		for key, value := range params {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: failed to send request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: status code: %d\n", res.StatusCode)
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	return resBody
}
