package main

import (
	"crypto/tls"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
)

type Computer struct {
	LastReport string `xml:"LastReportTime"`
	ID         string `xml:"ID"`
}

type Computers struct {
	Comp      string     `xml:"BESAPI"`
	Computers []Computer `xml:"Computer"`
}

func main() {
	//command line arguments input and parse
	var (
		bigfixurl = flag.String("bigfixurl", "https://bigfix.ucdavis.edu:52311", "URL for BigFix HTTP Service")
		user      = flag.String("user", "", "BigFix username")
		password  = flag.String("password", "", "BigFixPassword")
	)
	flag.Parse()
	var session = bigfixlogin(*user, *password, *bigfixurl)
	println(bigfixcomputers(*bigfixurl, session, *user, *password))

}

func bigfixlogin(user string, password string, bigfixurl string) *http.Client {
	//allows for invalid certs
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Transport: tr,
		Jar:       cookieJar,
	}
	req, err := http.NewRequest("GET", bigfixurl+"/api/login", nil)
	req.SetBasicAuth(user, password)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	println(s)
	return client
}

func bigfixcomputers(bigfixurl string, client *http.Client, user string, password string) bool {
	req, err := http.NewRequest("GET", bigfixurl+"/api/computers", nil)
	req.SetBasicAuth(user, password)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	f := []byte(s)
	var comps Computers
	xml := xml.Unmarshal(f, &comps)
	if xml != nil {
		log.Fatal(err)
	}
	//s := string(bodyText)
	fmt.Println(comps)
	return (true)
}
