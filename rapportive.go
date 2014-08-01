// Package rapportive is a client for the undocumented Rapportive API.
// Inspired by:  https://github.com/jordan-wright/rapportive
package main

import (
	"github.com/jmcvetta/napping"
	"log"
	"runtime"
	"strconv"
	"github.com/kr/pretty"
	"errors"
	"fmt"
	"net/http"
)

const (
	statusUrl = "https://rapportive.com/login_status?user_email=%s"
	contactsUrl = "https://profiles.rapportive.com/contacts/email/%s"

)

type LoginResult struct {
	Error string
	Token string `json:"session_token"`
}

type ContactsResult struct {
	Name string
	Contact Contact
}

type Contact struct {
	Name string
	Email string
	Twitter string `json:"twitter_username"`
	Location string
	Headline string
	Occupations []*Occupation
	Memberships []*Membership
}

type Occupation struct {
	JobTitle string `json:"job_title"`
	Company string
}

type Membership struct {
	Site string `json:"site_name"`
	Username string
	ProfileUrl string `json:"profile_url"`
}


func login(email string) (sessionToken string, err error){
	p := napping.Params{
		"user_email": email,
	}
	result := LoginResult{}
	resp, err := napping.Get(statusUrl, &p, &result, nil)
	if err != nil {
		log.Println(err)
		return
	}
	if resp.Status() != 200 {
		msg := fmt.Sprintf("Bad response status: %v", resp.Status())
		err = errors.New(msg)
		log.Println(err)
		return
	}
	sessionToken = result.Token
	return

}

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	email := "jason.mcvetta@gmail.com"
	token, err := login(email)
	if err != nil {
		log.Fatalln(err)
	}
	h := http.Header{}
	h.Add("X-Session-Token", token)
	s := napping.Session{
		Header: &h,
		Log: true,
	}
	url := fmt.Sprintf(contactsUrl, email)
	result := ContactsResult{}
	resp, err := s.Get(url, nil, &result, nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(resp.Status())
	c := result.Contact
	logPretty(c.Name)
	logPretty(c.Occupations)
	logPretty(c.Memberships)

}


func logPretty(x interface{}) {
	_, file, line, _ := runtime.Caller(1)
	lineNo := strconv.Itoa(line)
	s := file + ":" + lineNo + ": %# v\n"
	pretty.Logf(s, x)
}
