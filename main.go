package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Person declare data to be passed
type Person struct {
	Name      string
	Surname   string
	BirthCity string
	BirthDate string
	Gender    string
	EndPoint  string
	EPBuilt   bool
}

// URL defines the service provider domain
const URL = "http://webservices.dotnethell.it"

func main() {
	// Typical usage
	res, err := DoRequest("silvio", "berlusconi", "milano", "29/09/1936", "M")

	if err != nil {
		log.Println(err)
	}

	log.Println(res)
}

// DoRequest is the exit point
func DoRequest(name string, surname string, birthCity string, birthDate string, gender string) (string, error) {
	// instance anagraphic data
	p := newPerson(name, surname, birthCity, birthDate, gender)

	// define endpoint
	p.buildEndPoint()

	// do request - http GET
	XML, err := p.get()

	if err != nil {
		return "", err
	}

	result, err := p.formatData(XML)
	if err != nil {
		return "", err
	}

	return result, nil
}

// Get retrieve endpoint data
func (p *Person) get() (string, error) {
	var retVal = "no-value"

	if !p.EPBuilt {
		err := errors.New("no EndPoint built")
		return retVal, err
	}

	resp, err := http.Get(p.EndPoint)

	if err != nil {
		log.Println(err)
		return retVal, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
		return retVal, err
	}

	return string(body), nil
}

//formatData return string from input xml string
func (p *Person) formatData(inXML string) (string, error) {
	var fc string
	err := xml.Unmarshal([]byte(inXML), &fc)
	if err != nil {
		log.Println(err)
		return fc, err
	}

	return fc, nil
}

// buildEndPoint return the full endpoint
func (p *Person) buildEndPoint() {
	p.EndPoint = fmt.Sprintf(
		"%v/codicefiscale.asmx/CalcolaCodiceFiscale?Nome=%v&Cognome=%v&ComuneNascita=%v&DataNascita=%v&Sesso=%v",
		URL,
		p.Name,
		p.Surname,
		p.BirthCity,
		p.BirthDate,
		p.Gender,
	)

	p.EPBuilt = true
}

// newPerson return Person object
func newPerson(name string, surname string, birthCity string, birthDate string, gender string) Person {
	return Person{
		Name:      strings.ReplaceAll(name, " ", ""),
		Surname:   strings.ReplaceAll(surname, " ", ""),
		BirthCity: strings.ReplaceAll(birthCity, " ", ""),
		BirthDate: strings.ReplaceAll(birthDate, " ", ""),
		Gender:    strings.ReplaceAll(gender, " ", ""),
	}
}
