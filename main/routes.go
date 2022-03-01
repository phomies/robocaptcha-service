package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

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
	w.Write([]byte(xml.Header))
	w.Write(twimlOutput)
}

// TwiML Incoming Call Handler
func httpIncoming(w http.ResponseWriter, r *http.Request) {

	randomWord := generateRandomWord()
	randomNumber := generateRandomNumber()
	randomNumberSpaced := addSpaces(strconv.Itoa(randomNumber))
	actionUrl := fmt.Sprintf("/verify/%d/%s", randomNumber, randomWord)
	voicePrompt := fmt.Sprintf(
		`Your call is being screened for human verification. Please enter the numbers %s or say the word %s.`,
		randomNumberSpaced,
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
	w.Write([]byte(xml.Header))
	w.Write(twimlOutput)
}

// Perform verification on incoming calls
func httpVerify(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	speechResult := strings.ToLower(r.Form.Get("SpeechResult"))
	digitsEntered := r.Form.Get("Digits")
	numberDialed := r.Form.Get("Called")
	numberFrom := r.Form.Get("From")

	var userDialed User
	if numberDialed != "" {
		userDialed = *getUserFromMaskedNumber(numberDialed)
	} else {
		// Use mock user if anonymous phone
		numberFrom = os.Getenv("DEFAULT_NUMBER_FROM")
		userDialed = User{
			PhoneNumber: os.Getenv("ANON_FORWARD_TO"),
		}
	}

	forwardTo := userDialed.PhoneNumber

	// Removes the /request/ from "/request/123/abc"
	pathParams := strings.TrimPrefix(r.URL.String(), "/verify/")
	correct := strings.Split(pathParams, "/")

	if len(correct) != 2 {
		w.WriteHeader(500)
		w.Write([]byte("Error, non-2 arguments received"))
		return
	}

	fmt.Println(r.URL.String())
	fmt.Println(correct)
	fmt.Println(speechResult, digitsEntered)

	speechOk := strings.Contains(speechResult, correct[1])
	digitsOk := digitsEntered == correct[0]

	if speechOk || digitsOk {

		// Correct, allow the call to pass through
		twimlStruct := TwimlResponse{
			Say: &TwimlSay{
				Content: "You did it, redirecting your call now.",
				Voice:   "Polly.Brian",
			},
			Dial: &TwimlDial{
				DialNumber: forwardTo,
				CallerId:   numberFrom,
			},
		}

		twimlOutput, err := xml.Marshal(twimlStruct)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(xml.Header))
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
	w.Write([]byte(xml.Header))
	w.Write(twimlOutput)
}
