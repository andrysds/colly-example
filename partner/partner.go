package partner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/andrysds/dropship-checker/product"
)

const (
	usernameEnvKey          = "USERNAME"
	passwordEnvKey          = "PASSWORD"
	loginUrlEnvKey          = "LOGIN_URL"
	getProductBaseUrlEnvKey = "GET_PRODUCT_BASE_URL"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Partner struct {
	httpClient        httpClient
	authToken         string
	username          string
	password          string
	loginUrl          string
	getProductBaseUrl string
}

func NewPartner() *Partner {
	return &Partner{
		httpClient:        &http.Client{},
		username:          os.Getenv(usernameEnvKey),
		password:          os.Getenv(passwordEnvKey),
		loginUrl:          os.Getenv(loginUrlEnvKey),
		getProductBaseUrl: os.Getenv(getProductBaseUrlEnvKey),
	}
}

type loginResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

func (p *Partner) Login() error {
	jsonBody, err := json.Marshal(
		map[string]string{
			"username": p.username,
			"password": p.password,
		},
	)
	if err != nil {
		return err
	}
	reqBody := bytes.NewBuffer(jsonBody)

	req, err := http.NewRequest(http.MethodPost, p.loginUrl, reqBody)
	if err != nil {
		return err
	}

	res, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("got this status code: %d", res.StatusCode)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var loginRes *loginResponse
	err = json.Unmarshal(resBody, &loginRes)
	if err != nil {
		return err
	}

	p.authToken = loginRes.Data.Token
	return nil
}

func (p *Partner) GetProduct(slug string) (*product.Product, error) {
	req, err := http.NewRequest(http.MethodGet, p.getProductBaseUrl+slug, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", p.authToken)

	res, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("got this status code: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var pp *product.Product
	err = json.Unmarshal(body, &pp)
	return pp, err
}
