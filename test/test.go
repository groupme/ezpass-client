package test

import (
	"fmt"
	"github.com/bmizerany/pat"
	"log"
	"net/http"
	"net/http/httptest"
	"time"
)

const (
	TokenOk           = "200"
	TokenUnauthorized = "401"
	TokenNotFound     = "404"
	TokenTimeout      = "408"
	TokenError        = "500"
)

const (
	ResponseUnauthorized = "Unauthorized"
	ResponseNotFound     = "Not Found"
	ResponseError        = "Internal Server Error"
	ResponseUser         = `{
	  "user": {
	    "id":           "100",
	    "name":         "Brandon Keene",
	    "avatar_url":   "http://i.groupme.com/100",
	    "access_token": "success"
	  }
	}`
	ResponseMembership = `{
	  "user": {
	    "id":           "100",
	    "name":         "Brandon Keene",
	    "avatar_url":   "http://i.groupme.com/100",
	    "access_token": "success"
	  },
	  "membership": {
	  	"nickname": "B-money"
	  }
	}`
)

// Set up a test ezpass service in the style of net/http/httptest
func NewServer() *httptest.Server {
	mux := pat.New()
	mux.Get("/user", handler(ResponseUser))
	mux.Get("/groups/:group_id", handler(ResponseMembership))
	ts := httptest.NewServer(mux)
	return ts
}

func handler(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := token(r)
		log.Print("token:", t)
		switch t {
		case TokenOk:
			fmt.Fprint(w, body)
		case TokenNotFound:
			http.Error(w, ResponseNotFound, http.StatusNotFound)
		case TokenTimeout:
			time.Sleep(time.Second)
			http.Error(w, ResponseNotFound, http.StatusRequestTimeout)
		case TokenError:
			http.Error(w, ResponseError, http.StatusInternalServerError)
		default:
			http.Error(w, ResponseUnauthorized, http.StatusUnauthorized)
		}
	}
}

func token(r *http.Request) (token string) {
	token = r.Header.Get("X-Access-Token")
	if len(token) > 0 {
		return
	}

	token = r.FormValue("token")
	if len(token) > 0 {
		return
	}

	token = r.FormValue("access_token")
	if len(token) > 0 {
		return
	}

	return ""
}
