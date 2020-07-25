package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Initializer struct
type Initializer struct {
	token   string
	baseurl string
}

// New creates a new struct
func New(t string, b string) Initializer {
	u := Initializer{
		token:   t,
		baseurl: b,
	}
	if u.token == "" || u.baseurl == "" {
		panic("Token and Base url is needed to call this Function")
	}
	return u
}

func postRequest(url string, reqBody []byte, token string) (body string, err error) {
	var bearer = "Bearer " + token
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "Error", fmt.Errorf("error making http call: %w", err)
	}
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "Error", fmt.Errorf("error doing request: %w", err)
	}

	defer resp.Body.Close()

	bod, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Error", fmt.Errorf("error reading body: %w", err)
	}
	if resp.StatusCode != 200 {
		body = "Status code returned was: " + string(resp.StatusCode)
	} else {
		body = string(bod)
	}

	return
}

// RetrieveAuth retrieves authentication of a user
func (w Initializer) RetrieveAuth() (body string, err error) {

	endpoint := w.baseurl + "products/auths"

	body, err = postRequest(endpoint, nil, w.token)
	if err != nil {
		return "Error", fmt.Errorf("error retrieving auth token: %w", err)
	}
	return

}

// ByID fetch authentication info using the id of the authentication record.
func (w Initializer) ByID(i string) (body string, err error) {

	reqBody, err := json.Marshal(map[string]string{
		"id": i,
	})
	if err != nil {
		return "Error", fmt.Errorf("error converting json: %w", err)
	}

	endpoint := w.baseurl + "auth/getById"

	body, err = postRequest(endpoint, reqBody, w.token)
	if err != nil {
		return "Error", fmt.Errorf("error fetching auth using id: %w", err)
	}
	return
}

// ByOptions fetch authentication info using the options metadata you provided when setting up the widget.
func (w Initializer) ByOptions(page string, limit string, firstname string, lastname string) (body string, err error) {

	type Option struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	type payload struct {
		Page    string `json:"page"`
		Limit   string `json:"limit"`
		Options Option `json:"options"`
	}

	pl := payload{
		Page:  page,
		Limit: limit,
		Options: Option{
			FirstName: firstname,
			LastName:  lastname,
		},
	}

	reqBody, err := json.Marshal(pl)
	if err != nil {
		return "Error", fmt.Errorf("error marshalling json: %w", err)
	}
	url := w.baseurl + "auth/getByOptions"

	body, err = postRequest(url, reqBody, w.token)
	if err != nil {
		return "Error", fmt.Errorf("error retrieving auth byoptions: %w", err)
	}
	return
}

// ByCustomer fetch authentication info using the customer id
func (w Initializer) ByCustomer(customerID string) (body string, err error) {

	reqBody, err := json.Marshal(map[string]string{
		"customer": customerID,
	})
	if err != nil {
		return "Error", fmt.Errorf("error marshalling json: %w", err)
	}

	endpoint := w.baseurl + "auth/getByCustomer"

	body, err = postRequest(endpoint, reqBody, w.token)
	if err != nil {
		return "Error", fmt.Errorf("error retrieving auth bycustomer: %w", err)
	}
	return
}
