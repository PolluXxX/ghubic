package ghubic

import (
	"net/url"

	"github.com/polluxxx/goauth2"
)

const (
	AuthUrl  = "https://api.hubic.com/oauth/auth/"
	TokenUrl = "https://api.hubic.com/oauth/token/"
	ScopeUrl = "https://api.hubic.com/1.0/scope/scope"
	CallsUrl = "https://api.hubic.com/1.0/"
)

type HubicApi struct {
	Client *goauth2.Client
	Scope  string
}

func NewHubicApi(id, secret, redirectUrl string) (*HubicApi, error) {
	a := goauth2.Api{
		AuthUrl:  AuthUrl,
		TokenUrl: TokenUrl,
	}

	c := goauth2.NewClient(id, secret, redirectUrl, &a)

	// Getting full scope
	fullScope, err := getFullScope()
	if err != nil {
		return nil, err
	}

	ha := HubicApi{
		Client: c,
		Scope:  fullScope,
	}

	return &ha, nil
}

func (hubic *HubicApi) GetAuthUrl(state string) (uri *url.URL, err error) {
	m := make(map[string]string)

	m["state"] = state
	m["scope"] = hubic.Scope

	return hubic.Client.GetAuthUrl(m)
}

func (hubic *HubicApi) FinalizeAuth(code string) (account *Account, err error) {
	token, err := hubic.Client.Exchange(code)
	if err != nil {
		return nil, err
	}

	account, err = newAccount(token)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (hubic *HubicApi) GetAccountFromToken(token *goauth2.OAuthToken) (account *Account, err error) {
	return newAccount(token)
}
