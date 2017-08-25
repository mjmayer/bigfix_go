package bigfix

import (
	"crypto/tls"
	"encoding/xml"
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


//Sets up http session for interaction with bigfix.
func Bigfixlogin(user string, password string, bigfixurl string) *http.Client {
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
func Bigfixcomputers(bigfixurl string, client *http.Client, user string, password string) Computers {
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

//an Eval represents the time and plurality from a bigfix query
type Eval struct {
	Time        string `xml:"Time"`
	Pluralality string `xml:"Plurality"`
}

//An Query represent the query xml from bigfix
type Query struct {
	XMLName    xml.Name `xml:"Query"`
	Resource   string   `xml:"Resource,attr"`
	Result     []string `xml:"Result>Answer"`
	Evaluation Eval     `xml:"Evalutation`
}

//An BESQuery represents the XML returned from bigfix
type BESQuery struct {
	XMLName    string `xml:"BESAPI"`
	Query      Query  `xml:"Query"`
}

//Runs relevance query against bigfix server
func Bigfixquery(bigfixurl string, client *http.Client, user string, password string, relevance string) Query {
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
	f := []byte(s)
	var answer BESQuery
	xml := xml.Unmarshal(f, &answer)
	if xml != nil {
		log.Fatal(err)
	}
	return answer.Query
}
