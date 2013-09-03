package ezpass

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	URL     string
	Timeout = 50 * time.Millisecond
)

var (
	ErrUnauthorized = errors.New("ezpass: unauthorized")
	ErrNotFound     = errors.New("ezpass: not found (user is not in group)")
	ErrTimeout      = errors.New("ezpass: timeout")
	ErrUnknown      = errors.New("ezpass: unknown error")
)

type Pass struct {
	User struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		AvatarUrl   string `json:"avatar_url"`
		AccessToken string `json:"access_token"`
	} `json:"user"`
	Membership struct {
		Nickname string `json:"nickname"`
	} `json:"membership,omitempty"`
}

func init() {
	// maybe convert this jankiness into a Client singleton
	URL = os.Getenv("EZPASS_URL") // override me in test, production, etc.
}

func Get(token string) (*Pass, error) {
	return performWithTimeout(userUrl(token))
}

func GetMembership(token string, groupId string) (*Pass, error) {
	return performWithTimeout(groupUrl(token, groupId))
}

func performWithTimeout(url string) (*Pass, error) {
	p := make(chan *Pass)
	e := make(chan error)

	go func() {
		pass, err := perform(url)
		if err != nil {
			e <- err
		}
		p <- pass
	}()

	select {
	case pass := <-p:
		return pass, nil
	case err := <-e:
		return nil, err
	case <-(time.After(Timeout)):
		return nil, ErrTimeout
	}
}

func perform(url string) (*Pass, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	switch res.StatusCode {
	case http.StatusOK:
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		pass := &Pass{}
		err = json.Unmarshal(body, pass)
		if err != nil {
			return nil, err
		}
		return pass, nil
	case http.StatusNotFound:
		return nil, ErrNotFound
	case http.StatusUnauthorized:
		return nil, ErrUnauthorized
	default:
		return nil, ErrUnknown
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

func userUrl(token string) string {
	return fmt.Sprintf("%s/user?access_token=%s", URL, token)
}

func groupUrl(token string, groupId string) string {
	return fmt.Sprintf("%s/groups/%s?access_token=%s", URL, groupId, token)
}
