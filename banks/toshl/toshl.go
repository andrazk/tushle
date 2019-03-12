package toshl

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"tushle"
)

const (
	acronym = "toshl"
	baseURL = "https://api.toshl.com"
)

// Toshl defines Toshl client
type Toshl struct {
	personalToken string

	cli   tushle.Cli
	creds tushle.Credentials
}

// New returns Toshl API client.
func New(cli tushle.Cli, creds tushle.Credentials) tushle.Banker {
	return &Toshl{
		cli:   cli,
		creds: creds,
	}
}

// NewClientWithPersonalToken returns a new Toshl API client.
func NewClientWithPersonalToken(pt string) *Toshl {
	return &Toshl{
		personalToken: pt,
	}
}

// Account Model
type Account struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Balance  float64  `json:"balance"`
	Currency Currency `json:"currency"`
}

// Currency holds code and other data
type Currency struct {
	Code string `json:"code"`
}

// Acronym returns bank acronym.
func (t *Toshl) Acronym() string {
	return acronym
}

func (t *Toshl) token() (string, error) {
	if t.personalToken != "" {
		return t.personalToken, nil
	}

	var token string
	err := t.creds.Get(acronym, &token)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return base64.StdEncoding.EncodeToString([]byte(token)), nil
}

func (t *Toshl) do(req *http.Request, buf interface{}) error {
	token, err := t.token()
	if err != nil {
		return errors.WithStack(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return errors.Errorf(
			"request failed with code %d and message %s",
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	// Fetch entity when created
	if resp.StatusCode == 201 {
		url := fmt.Sprintf("%s/%s", baseURL, strings.TrimLeft(resp.Header.Get("Location"), "/"))
		req2, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		return t.do(req2, buf)
	}

	if buf != nil {
		return json.NewDecoder(resp.Body).Decode(&buf)
	}
	return nil
}
