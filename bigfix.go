package main

import (
	"crypto/tls"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	//command line arguments input and parse
	var (
		bigfixurl = flag.String("bigfixurl", "https://bigfix.ucdavis.edu:52311/api/login", "URL for BigFix HTTP Service")
		user      = flag.String("user", "", "BigFix username")
		password  = flag.String("password", "", "BigFixPassword")
	)
	flag.Parse()

	//allows for invalid certs
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", *bigfixurl, nil)
	req.SetBasicAuth(*user, *password)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	println(s)
}
