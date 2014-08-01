// Package rapportive is a client for the undocumented Rapportive API.
// Inspired by:  https://github.com/jordan-wright/rapportive
package rapportive

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


func Login(email string) (sessionToken string, err error){
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

//func main() {
//	log.SetFlags(log.Ltime | log.Lshortfile)
//	email := "jason.mcvetta@gmail.com"
//	token, err := Login(email)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	c, err := QueryContacts(token, email)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	logPretty(c)
//}

func QueryContacts(sessionToken, email string) (*Contact, error) {
	h := http.Header{}
	h.Add("X-Session-Token", sessionToken)
	s := napping.Session{
		Header: &h,
	}
	url := fmt.Sprintf(contactsUrl, email)
	result := ContactsResult{}
	resp, err := s.Get(url, nil, &result, nil)
	if err != nil {
		log.Println(err)
		return nil, err
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
