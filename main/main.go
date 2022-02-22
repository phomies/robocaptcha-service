package main

import (
	"encoding/xml"
	"net/http"
)

// TwiML Structs
type TwimlResponse struct {
	XMLName xml.Name    `xml:"Response"`
	Gather  TwimlGather `xml:",omitempty"`
}

type TwimlGather struct {
	XMLName   xml.Name `xml:"Gather"`
	Say       TwimlSay `xml:",omitEmpty"`
	Input     string   `xml:"input,attr"`
	Timeout   int      `xml:"timeout,attr"`
	NumDigits int      `xml:"numDigits,attr"`
	action    string   `xml:"action,attr"`
	language  string   `xml:"language,attr"`
}

type TwimlSay struct {
	XMLName xml.Name `xml:"Say"`
}

// Main function
func main() {
	http.HandleFunc("/incoming", twiml)
	http.HandleFunc("/", healthcheck)
	http.ListenAndServe(":3000", nil)
}

// Web service health check
func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("It works!!!"))
}

// TwiML Incoming Call Handler
func twiml(w http.ResponseWriter, r *http.Request) {
	// twiml := TwiML{Say: "Hello World!"}

	xmlOutput, err := xml.Marshal(twiml)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(xmlOutput)
}
