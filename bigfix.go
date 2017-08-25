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

//An Computer represents a computer object returend from bigfix
type Computer struct {
	LastReport string `xml:"LastReportTime"`
	ID         string `xml:"ID"`
}

//An Computers represents the XML returned from bigfix
type Computers struct {
	XMLName   string     `xml:"BESAPI"`
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
	//println(len(bigfixcomputers(*bigfixurl, session, *user, *password).Computers))
	bigfixquery(*bigfixurl, session, *user, *password, "names of bes computers")
}

//Sets up http session for interaction with bigfix.
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

//Returns Computer structure containing computer with lastreporttime, and ID
func bigfixcomputers(bigfixurl string, client *http.Client, user string, password string) Computers {
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
	return comps
}

type Result string
type Answer string

type Eval struct {
	//XMLName     xml.Name `xml:"Evalutation"`
	Time        string `xml:"Time"`
	Pluralality string `xml:"Plurality"`
}

type Query struct {
	XMLName    xml.Name `xml:"Query"`
	Resource   string   `xml:"Resource,attr"`
	Answers    []Answer `xml:"Result>Answer"`
	Evaluation []Eval   `xml:"Evalutation`
}

//An query represents the XML returned from bigfix
type BESQuery struct {
	XMLName string `xml:"BESAPI"`
	Query   Query  `xml:"Query"`
	//Answers []Answer `xml:"Query>Result>Answer"`
}

//Runs relevance query against bigfix server
func bigfixquery(bigfixurl string, client *http.Client, user string, password string, relevance string) BESQuery {
	req, err := http.NewRequest("GET", bigfixurl+"/api/query", nil)
	req.SetBasicAuth(user, password)
	q := req.URL.Query()
	q.Add("relevance", relevance)
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	//println(s)
	f := []byte(s)
	var answer BESQuery
	xml := xml.Unmarshal(f, &answer)
	if xml != nil {
		log.Fatal(err)
	}
	fmt.Println(answer)
	return answer
}
