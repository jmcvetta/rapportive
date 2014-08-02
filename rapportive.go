// Copyright (c) 2014 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for
// details.  Resist intellectual serfdom - the ownership of ideas is akin to
// slavery.

// Package rapportive is a client for the undocumented Rapportive API.
// Inspired by:  https://github.com/jordan-wright/rapportive
package rapportive

import (
	"errors"
	"fmt"
	"github.com/jmcvetta/napping"
	"github.com/kr/pretty"
	"log"
	"net/http"
	"runtime"
	"strconv"
)

const (
	statusUrl   = "https://rapportive.com/login_status?user_email=%s"
	contactsUrl = "https://profiles.rapportive.com/contacts/email/%s"
)

type loginResult struct {
	Error string
	Token string `json:"session_token"`
}

type contactsResult struct {
	Name    string
	Contact Contact
}

type Contact struct {
	Name        string
	Email       string
	Twitter     string `json:"twitter_username"`
	Location    string
	Headline    string
	Occupations []*Occupation
	Memberships []*Membership
}

type Occupation struct {
	JobTitle string `json:"job_title"`
	Company  string
}

type Membership struct {
	Site       string `json:"site_name"`
	Username   string
	ProfileUrl string `json:"profile_url"`
}

func login(email string) (sessionToken string, err error) {
	p := napping.Params{
		"user_email": email,
	}
	result := loginResult{}
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

func Query(email string) (c *Contact, err error) {
	// log.SetFlags(log.Ltime | log.Lshortfile)
	token, err := login(email)
	if err != nil {
		return
	}
	c, err = getContact(token, email)
	if err != nil {
		return
	}
	logPretty(c)
	return
}

var RateLimitError = errors.New("Rate Limit Error")

func getContact(sessionToken, email string) (*Contact, error) {
	h := http.Header{}
	h.Add("X-Session-Token", sessionToken)
	s := napping.Session{
		Header: &h,
	}
	url := fmt.Sprintf(contactsUrl, email)
	result := contactsResult{}
	resp, err := s.Get(url, nil, &result, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if resp.Status() == 429 {
		log.Println(RateLimitError)
		return nil, RateLimitError
	}
	if resp.Status() != 200 {
		msg := fmt.Sprintf("Bad response status: %v", resp.Status())
		err = errors.New(msg)
		log.Println(err)
		return nil, err
	}
	return &result.Contact, nil
}

func logPretty(x interface{}) {
	_, file, line, _ := runtime.Caller(1)
	lineNo := strconv.Itoa(line)
	s := file + ":" + lineNo + ": %# v\n"
	pretty.Logf(s, x)
}
