package main

import (
	"encoding/xml"
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
