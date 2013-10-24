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
	assertEqual(t, res.Code, http.StatusUnauthorized)
	assertEqual(t, res.Header().Get("Content-Type"), "application/json")
	assertEqual(t, res.Body.String(), `{"meta":{"error":"ezpass: unauthorized"}}`+"\n")

	// ok
	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Add("X-Access-Token", ezpasstest.TokenOk)
	res = httptest.NewRecorder()

	AuthHandler(handler)(res, req)
	assertEqual(t, res.Code, http.StatusOK)
	assertEqual(t, res.Body.String(), `Brandon Keene`)

	// membership
	req, _ = http.NewRequest("GET", "/?group_id=1", nil)
	req.Header.Add("X-Access-Token", ezpasstest.TokenOk)
	res = httptest.NewRecorder()
	AuthHandler(handler)(res, req)
	assertEqual(t, res.Code, http.StatusOK)
	assertEqual(t, res.Body.String(), `B-money`)

	// membership - github.com/bmizerany/pat style
	// this exploits the fact that RawQuery will parse :group_id
	// in actual use, the URL would be: /a/:group_id/b
	req, _ = http.NewRequest("GET", "/?:group_id=1", nil)
	req.Header.Add("X-Access-Token", ezpasstest.TokenOk)
	res = httptest.NewRecorder()
	AuthHandler(handler)(res, req)
	assertEqual(t, res.Code, http.StatusOK)
	assertEqual(t, res.Body.String(), `B-money`)
}

func TestGet(t *testing.T) {
	URL = ezpasstest.NewServer().URL

	// ok
	pass, err := Get(ezpasstest.TokenOk)

	assertEqual(t, err, nil)
	assertEqual(t, pass.User.Id, "100")
	assertEqual(t, pass.User.Name, "Brandon Keene")
	assertEqual(t, pass.User.AvatarUrl, "http://i.groupme.com/100")

	pass, err = Get(ezpasstest.TokenUnauthorized)
	assertEqual(t, err, ErrUnauthorized)

	pass, err = Get(ezpasstest.TokenTimeout)
	assertEqual(t, err, ErrTimeout)

	pass, err = Get(ezpasstest.TokenError)
	assertEqual(t, err, ErrUnknown)
}

func TestGetMembership(t *testing.T) {
	URL = ezpasstest.NewServer().URL
	groupId := "1"

	// ok
	pass, err := GetMembership(ezpasstest.TokenOk, groupId)

	assertEqual(t, err, nil)
	assertEqual(t, pass.User.Id, "100")
	assertEqual(t, pass.User.Name, "Brandon Keene")
	assertEqual(t, pass.User.AvatarUrl, "http://i.groupme.com/100")
	assertEqual(t, pass.Membership.Nickname, "B-money")

	pass, err = GetMembership(ezpasstest.TokenUnauthorized, groupId)
	assertEqual(t, err, ErrUnauthorized)

	pass, err = GetMembership(ezpasstest.TokenNotFound, groupId)
	assertEqual(t, err, ErrNotFound)

	pass, err = GetMembership(ezpasstest.TokenTimeout, groupId)
	assertEqual(t, err, ErrTimeout)

	pass, err = GetMembership(ezpasstest.TokenError, groupId)
	assertEqual(t, err, ErrUnknown)
}

func TestToken(t *testing.T) {
	var r *http.Request
	reader := strings.NewReader("")

	r, _ = http.NewRequest("GET", "http://example.com?token=foo", reader)
	assertEqual(t, Token(r), "foo")

	r, _ = http.NewRequest("GET", "http://example.com?access_token=foo", reader)
	assertEqual(t, Token(r), "foo")

	r, _ = http.NewRequest("GET", "http://example.com", reader)
	r.Header.Set("X-Access-Token", "foo")
	assertEqual(t, Token(r), "foo")

	r, _ = http.NewRequest("GET", "http://example.com", reader)
	assertEqual(t, Token(r), "")
}

func assertEqual(t *testing.T, actual interface{}, expected interface{}) {
	if actual != expected {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
