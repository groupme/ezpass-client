package ezpass

import (
	"fmt"
	"github.com/groupme/ezpass-client/ezpasstest"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthHandler(t *testing.T) {
	URL = ezpasstest.NewServer().URL

	handler := func(w http.ResponseWriter, r *http.Request, pass *Pass) {
		if len(pass.Membership.Nickname) > 0 {
			fmt.Fprintf(w, pass.Membership.Nickname)
		} else {
			fmt.Fprintf(w, pass.User.Name)
		}
	}

	// unauthorized
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	AuthHandler(handler)(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Error("should be unauthorized")
	}

	// ok
	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Add("X-Access-Token", ezpasstest.TokenOk)
	res = httptest.NewRecorder()

	AuthHandler(handler)(res, req)

	if res.Code != http.StatusOK {
		t.Error("should be ok")
	}

	if res.Body.String() != "Brandon Keene" {
		t.Error("failed to get pass and user data")
	}

	// membership
	req, _ = http.NewRequest("GET", "/?group_id=1", nil)
	req.Header.Add("X-Access-Token", ezpasstest.TokenOk)
	res = httptest.NewRecorder()

	AuthHandler(handler)(res, req)

	if res.Code != http.StatusOK {
		t.Error("should be ok")
	}

	if res.Body.String() != "B-money" {
		t.Error("failed to get pass and membership data")
	}
}

func TestGet(t *testing.T) {
	URL = ezpasstest.NewServer().URL

	// ok
	pass, err := Get(ezpasstest.TokenOk)
	if err != nil {
		t.Error(err)
	}

	if pass.User.Id != "100" {
		t.Error("User: id is incorrect")
	}

	if pass.User.Name != "Brandon Keene" {
		t.Error("User: name is incorrect")
	}

	if pass.User.AvatarUrl != "http://i.groupme.com/100" {
		t.Error("User: avatar_url is incorrect")
	}

	// unauthorized
	pass, err = Get(ezpasstest.TokenUnauthorized)
	if err != ErrUnauthorized {
		t.Error("User: error should be ErrUnauthorized")
	}

	// timeout
	pass, err = Get(ezpasstest.TokenTimeout)
	if err != ErrTimeout {
		t.Error("User: error should be ErrTimeout")
	}

	// error
	pass, err = Get(ezpasstest.TokenError)
	if err != ErrUnknown {
		t.Error("User: error should be ErrUnknown")
	}
}

func TestGetMembership(t *testing.T) {
	URL = ezpasstest.NewServer().URL
	groupId := "1"

	// ok
	pass, err := GetMembership(ezpasstest.TokenOk, groupId)
	if err != nil {
		t.Error(err)
	}

	if pass.User.Id != "100" {
		t.Error("User: id is incorrect")
	}

	if pass.User.Name != "Brandon Keene" {
		t.Error("User: name is incorrect")
	}

	if pass.User.AvatarUrl != "http://i.groupme.com/100" {
		t.Error("User: avatar_url is incorrect")
	}

	if pass.Membership.Nickname != "B-money" {
		t.Error("User: nickname is incorrect")
	}

	// unauthorized
	pass, err = GetMembership(ezpasstest.TokenUnauthorized, groupId)
	if err != ErrUnauthorized {
		t.Error("User: error should be ErrUnauthorized")
	}

	// not found
	pass, err = GetMembership(ezpasstest.TokenNotFound, groupId)
	if err != ErrNotFound {
		t.Error("User: error should be ErrNotFound")
	}

	// timeout
	pass, err = GetMembership(ezpasstest.TokenTimeout, groupId)
	if err != ErrTimeout {
		t.Error("User: error should be ErrTimeout")
	}

	// error
	pass, err = GetMembership(ezpasstest.TokenError, groupId)
	if err != ErrUnknown {
		t.Error("User: error should be ErrUnknown")
	}
}

func TestToken(t *testing.T) {
	var r *http.Request
	reader := strings.NewReader("")

	r, _ = http.NewRequest("GET", "http://example.com?token=foo", reader)
	if Token(r) != "foo" {
		t.Error("failed to get 'token' param")
	}

	r, _ = http.NewRequest("GET", "http://example.com?access_token=foo", reader)
	if Token(r) != "foo" {
		t.Error("failed to get 'access_token' param")
	}

	r, _ = http.NewRequest("GET", "http://example.com", reader)
	r.Header.Set("X-Access-Token", "foo")
	if Token(r) != "foo" {
		t.Error("failed to get 'X-Access-Token' header")
	}

	r, _ = http.NewRequest("GET", "http://example.com", reader)
	if Token(r) != "" {
		t.Error("failed to detect missing token")
	}
}
