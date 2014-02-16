package ghubic

import (
	"encoding/json"
	"errors"
	"fmt"
	"goauth2"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
    "io"
)

type Account struct {
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Email        string `json:"email"`
	Activated    bool   `json:"activated"`
	CreationDate string `json:"creationDate"`
	Language     string `json:"language"`
	Status       string `json:"status"`
	Offer        string `json:"offer"`

	// This field is only available after GetCredentials is called.
	Credentials *AccountCredentials
	// This field is only available after GetUsage is called.
	Usage *AccountUsage

	Token *goauth2.OAuthToken
}

type AccountCredentials struct {
	Token    string `json:"token"`
	Endpoint string `json:"endpoint"`
	Expires  string `json:"expires"`
}

type AccountUsage struct {
	Used  int `json:"used"`
	Quota int `json:"quota"`
}

func newAccount(token *goauth2.OAuthToken) (account *Account, err error) {
	account = new(Account)
	account.Token = token

	body, err := account.Call("account", "GET", nil, false)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (account *Account) GetCredentials() (credentials *AccountCredentials, err error) {
	if account.Token == nil {
		return nil, errors.New("Account is not linked to token")
	}

	body, err := account.Call("account/credentials", "GET", nil, false)
	if err != nil {
		return nil, err
	}

	credentials = new(AccountCredentials)
	err = json.Unmarshal(body, credentials)
	if err != nil {
		return nil, err
	}

	account.Credentials = credentials

	return credentials, nil
}

func (account *Account) GetUsage() (usage *AccountUsage, err error) {
	if account.Token == nil {
		return nil, errors.New("Account is not linked to token")
	}

	body, err := account.Call("account/usage", "GET", nil, false)
	if err != nil {
		return nil, err
	}

	usage = new(AccountUsage)
	err = json.Unmarshal(body, usage)
	if err != nil {
		return nil, err
	}

	account.Usage = usage

	return usage, nil
}

func (account *Account) Call(endpoint, method string, params map[string]string, refreshed bool) (body []byte, err error) {

	// Building query
	u := url.Values{}

	for p, v := range params {
		u.Set(p, v)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", CallsUrl, endpoint), strings.NewReader(u.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", account.Token.Type, account.Token.AccessToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		// Catch response
		apiErr := new(ApiError)
		err = json.Unmarshal(body, apiErr)
		if err != nil {
			return nil, err
		}

		// Check expired token
		if apiErr.ApiError == "invalid_token" && apiErr.ApiDesc == "expired" && !refreshed {
			err = account.Token.Refresh()
			if err != nil {
				return nil, err
			}

			return account.Call(endpoint, method, params, true)
		}

		apiErr.HTTPCode = resp.StatusCode
		return nil, apiErr
	}

	return body, nil
}

func (account* Account) AddFile(path, filename string, content io.ReadCloser) error {
    _, err := account.GetCredentials()
    if err != nil {
        return err
    }

    if path == "" {
        path = "/"
    }

    req, err := http.NewRequest("PUT", account.Credentials.Endpoint + "/default" + path + filename, content)
    req.Header.Set("X-Auth-Token", account.Credentials.Token)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }

    if resp.StatusCode != 201 {
        return nil
    }

    return nil
}
