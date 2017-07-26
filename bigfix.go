package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
)

var bigfixapiurl = "https://bigfix.ucdavis.edu:52311/api/login"
var user = ""
var password = ""

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", bigfixapiurl, nil)
	req.SetBasicAuth(user, password)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	println(s)
}
