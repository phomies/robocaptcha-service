package main

import (
	"encoding/xml"
	"strconv"
)

// TwiML Structs
type TwimlResponse struct {
	XMLName  xml.Name     `xml:"Response"`
	Say      *TwimlSay    `xml:",omitempty"`
	Gather   *TwimlGather `xml:",omitempty"`
	Dial     *TwimlDial   `xml:",omitempty"`
	Redirect string       `xml:",omitempty"`
}

type TwimlGather struct {
	XMLName   xml.Name  `xml:"Gather"`
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

type TwimlDial struct {
	XMLName    xml.Name `xml:"Dial"`
	DialNumber string   `xml:",chardata"`
	CallerId   string   `xml:"callerId,attr"`
}

func structEndCall() *TwimlResponse {
	return &TwimlResponse{
		Say: &TwimlSay{
			Voice:   "Polly.Brian",
			Content: "Call has timed out, goodbye.",
		},
	}
}

func structUserNotFound() *TwimlResponse {
	return &TwimlResponse{
		Say: &TwimlSay{
			Voice:   "Polly.Brian",
			Content: "There are no such user on this platform, please try again later.",
		},
	}
}

func structIncorrectEntered(times int) *TwimlResponse {
	return &TwimlResponse{
		Say: &TwimlSay{
			Content: "Incorrect, please try again!",
			Voice:   "Polly.Brian",
		},
		Redirect: "/incoming/" + strconv.Itoa(times+1),
	}
}

func structHelloWorld() *TwimlResponse {
	return &TwimlResponse{
		Say: &TwimlSay{
			Content: "Hello world!",
		},
	}
}

func structVerifyCall(actionUrl string, voicePrompt string, times int) TwimlResponse {
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
		Redirect: "/incoming/" + strconv.Itoa(times+1),
	}
	return twimlStruct
}

func structForwardingCall(forwardTo string, numberFrom string) TwimlResponse {
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
	return twimlStruct
}
