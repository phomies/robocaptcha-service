package main

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

// TwiML Structs
type TwimlResponse struct {
	XMLName  xml.Name     `xml:"Response"`
	Say      *TwimlSay    `xml:",omitempty"`
	Gather   *TwimlGather `xml:",omitempty"`
	Redirect string       `xml:",omitempty"`
}

type TwimlGather struct {
	XMLNam    xml.Name  `xml:"Gather"`
	Say       *TwimlSay `xml:",omitempty"`
	Input     string    `xml:"input,attr,omitempty"`
	Timeout   int       `xml:"timeout,attr,omitempty"`
	NumDigits int       `xml:"numDigits,attr,omitempty"`
	Action    string    `xml:"action,attr,omitempty"`
	Language  string    `xml:"language,attr,omitempty"`
}

type TwimlSay struct {
	XMLName xml.Name `xml:"Say"`
	Voice   string   `xml:"voice,attr,omitempty"`
	Content string   `xml:",chardata"`
}

// Main function
func main() {
	http.HandleFunc("/incoming", httpIncoming)
	http.HandleFunc("/verify/", httpVerify)
	http.HandleFunc("/", healthcheck)
	http.ListenAndServe(":8080", nil)
}

// Web service health check
func healthcheck(w http.ResponseWriter, r *http.Request) {
	twimlStruct := TwimlResponse{
		Say: &TwimlSay{
			Content: "Hello world!",
		},
	}
	twimlOutput, err := xml.Marshal(twimlStruct)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/xml")
	w.Write(twimlOutput)
}

// TwiML Incoming Call Handler
func httpIncoming(w http.ResponseWriter, r *http.Request) {

	randomWord := generateRandomWord()
	randomNumber := generateRandomNumber()
	actionUrl := fmt.Sprintf("/verify/%d/%s", randomNumber, randomWord)
	voicePrompt := fmt.Sprintf(
		`Your call is being screened for human verification, on your dialpad, please enter %d or say the word %s.`,
		randomNumber,
		randomWord,
	)

	twimlStruct := TwimlResponse{
		Gather: &TwimlGather{
			Timeout:   3,
			NumDigits: 4,
			Action:    actionUrl,
			Input:     "dtmf speech",
			Language:  "en-SG",
			Say: &TwimlSay{
				Voice:   "Polly.Brian",
				Content: voicePrompt,
			},
		},
		Redirect: "/incoming",
	}

	twimlOutput, err := xml.Marshal(twimlStruct)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/xml")
	w.Write(twimlOutput)
}

// Perform verification on incoming calls
func httpVerify(w http.ResponseWriter, r *http.Request) {
	speechResult := r.Header.Get("SpeechResult")
	digitsEntered := r.Header.Get("Digits")

	// Removes the /request/ from "/request/123/abc"
	pathParams := strings.TrimPrefix(r.URL.String(), "/request/")
	correct := strings.Split(pathParams, "/")

	speechOk := strings.Contains(speechResult, correct[1])
	digitsOk := digitsEntered == correct[0]

	if speechOk || digitsOk {

		// Correct, allow the call to pass through
		twimlStruct := TwimlResponse{
			Say: &TwimlSay{
				Content: "You did it, redirecting your call now.",
				Voice:   "Polly.Brian",
			},
		}

		twimlOutput, err := xml.Marshal(twimlStruct)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/xml")
		w.Write(twimlOutput)
		return

	}

	// Incorrect, try again from scratch
	twimlStruct := TwimlResponse{
		Say: &TwimlSay{
			Content: "Incorrect, please try again!",
			Voice:   "Polly.Brian",
		},
		Redirect: "/incoming",
	}

	twimlOutput, err := xml.Marshal(twimlStruct)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/xml")
	w.Write(twimlOutput)
}

func generateRandomWord() string {
	wordList := []string{
		"singular",
		"designer",
		"agenda",
		"commitment",
		"tradition",
		"conference",
		"potential",
		"producer",
	}
	randomIdx := rand.Intn(len(wordList))
	return wordList[randomIdx]
}

func generateRandomNumber() int {
	randomNumber := 1000 + rand.Intn(9000)
	return randomNumber
}
