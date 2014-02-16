package ghubic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type right struct {
	Text  string `json:"text"`
	Value string `json:"value"`
}

func (r *right) String() string {
	return fmt.Sprintf("%s", r.Value)
}

type endpoint struct {
	Name   string   `json:"name"`
	Rights []*right `json:"rights"`
}

func (e *endpoint) String() string {
	rights := ""
	for _, value := range e.Rights {
		rights += value.String()
	}

	return fmt.Sprintf("%s.%s", e.Name, rights)
}

type scope []endpoint

func (s scope) String() string {
	ep := ""
	// Join with ',' char
	for _, value := range s {
		ep += value.String() + ","
	}

	if len(ep) > 0 {
		// Remove trailing character
		ep = ep[0 : len(ep)-1]
	}

	return fmt.Sprintf("%s", ep)
}

// Request full scope from hubiC API
func getFullScope() (fullScope string, err error) {
	resp, err := http.Get(ScopeUrl)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var endpoints scope
	err = json.Unmarshal(body, &endpoints)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", endpoints), nil
}
