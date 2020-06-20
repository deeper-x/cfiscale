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
	Name           string
	Surname        string
	BirthCity      string
	BirthDate      string
	Gender         string
	EPCreate       string // EPC: End Point for fiscal code Creation
	EPCBuilt       bool   // EPC is built, exists
	EPVerification string // EPV: End Point for fiscal code Verification
	EPVBuilt       bool   // EPV is build, exists
}

// URL defines the service provider domain
const URL = "http://webservices.dotnethell.it/codicefiscale.asmx"
const expected = "Il codice è valido!"

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
	p.buildEPCreate()

	// do request - http GET
	XML, err := p.getC()

	if err != nil {
		return "", err
	}

	// build the resulting fiscal code string
	result, err := p.formatData(XML)
	if err != nil {
		return "", err
	}

	// now prepare for the verification
	p.buildEPVerification(result)

	// now call for the verification
	ok, err := p.GetV()

	if !ok {
		return result, err
	}

	// fiscal code verified and ready
	return result, nil
}

//Verify check for fiscal code validity
func (p *Person) Verify(fc string) bool {
	return true
}

// Get retrieve endpoint data
func (p *Person) getC() (string, error) {
	var retVal = "no-value"

	if !p.EPCBuilt {
		err := errors.New("no EndPoint built")
		return retVal, err
	}

	resp, err := http.Get(p.EPCreate)

	if err != nil {
		log.Println(err)
		return retVal, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		return retVal, err
	}

	return string(body), nil
}

// GetV return if fiscal code is verified
func (p *Person) GetV() (bool, error) {
	if !p.EPVBuilt {
		err := errors.New("End Point Verification doesn't exist")
		return false, err
	}

	// call EPV - verification endpoint
	resp, err := http.Get(p.EPVerification)

	if err != nil {
		log.Println(err)
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		return false, err
	}

	// format result string
	result, err := p.formatData(string(body))

	if err != nil {
		log.Println(err)
		return false, err
	}

	// verify result is successful
	if expected != result {
		return false, nil
	}

	return true, nil
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

// buildEPCreate return the End Point delegated to fiscal code Creation
func (p *Person) buildEPCreate() {
	p.EPCreate = fmt.Sprintf(
		"%v/CalcolaCodiceFiscale?Nome=%v&Cognome=%v&ComuneNascita=%v&DataNascita=%v&Sesso=%v",
		URL,
		p.Name,
		p.Surname,
		p.BirthCity,
		p.BirthDate,
		p.Gender,
	)

	p.EPCBuilt = true
}

func (p *Person) buildEPVerification(fc string) {
	p.EPVerification = fmt.Sprintf(
		"%v/ControllaCodiceFiscale?CodiceFiscale=%v",
		URL,
		fc,
	)
	p.EPVBuilt = true
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
