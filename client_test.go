package ezpass

import (
	"github.com/groupme/ezpass-client/test"
	"net/http"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	URL = test.NewServer().URL

	// ok
	pass, err := Get(test.TokenOk)
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
	pass, err = Get(test.TokenUnauthorized)
	if err != ErrUnauthorized {
		t.Error("User: error should be ErrUnauthorized")
	}

	// timeout
	pass, err = Get(test.TokenTimeout)
	if err != ErrTimeout {
		t.Error("User: error should be ErrTimeout")
	}

	// error
	pass, err = Get(test.TokenError)
	if err != ErrUnknown {
		t.Error("User: error should be ErrUnknown")
	}
}

func TestGetMembership(t *testing.T) {
	URL = test.NewServer().URL
	groupId := "1"

	// ok
	pass, err := GetMembership(test.TokenOk, groupId)
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
	pass, err = GetMembership(test.TokenUnauthorized, groupId)
	if err != ErrUnauthorized {
		t.Error("User: error should be ErrUnauthorized")
	}

	// not found
	pass, err = GetMembership(test.TokenNotFound, groupId)
	if err != ErrNotFound {
		t.Error("User: error should be ErrNotFound")
	}

	// timeout
	pass, err = GetMembership(test.TokenTimeout, groupId)
	if err != ErrTimeout {
		t.Error("User: error should be ErrTimeout")
	}

	// error
	pass, err = GetMembership(test.TokenError, groupId)
	if err != ErrUnknown {
		t.Error("User: error should be ErrUnknown")
	}
}

func TestToken(t *testing.T) {
	var r *http.Request
	reader := strings.NewReader("")

	r, _ = http.NewRequest("GET", "http://example.com?token=foo", reader)
	if token(r) != "foo" {
		t.Error("failed to get 'token' param")
	}

	r, _ = http.NewRequest("GET", "http://example.com?access_token=foo", reader)
	if token(r) != "foo" {
		t.Error("failed to get 'access_token' param")
	}

	r, _ = http.NewRequest("GET", "http://example.com", reader)
	r.Header.Set("X-Access-Token", "foo")
	if token(r) != "foo" {
		t.Error("failed to get 'X-Access-Token' header")
	}

	r, _ = http.NewRequest("GET", "http://example.com", reader)
	if token(r) != "" {
		t.Error("failed to detect missing token")
	}
}
