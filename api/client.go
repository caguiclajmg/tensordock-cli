package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Response struct {
	Success bool `json:"success"`
}

type ServersGetSingleResponse struct {
	Response
	Server Server `json:"server"`
}

type ServersListResponse struct {
	Response
	Servers map[string]Server `json:"servers"`
}

type BillingDetailsResponse struct {
	Response
	BillingDetails
}

type Client struct {
	BaseUrl  string
	ApiKey   string
	ApiToken string
}

func (client *Client) get(path string, params map[string]string) (*json.RawMessage, error) {
	query := url.Values{}
	query.Add("api_key", client.ApiKey)
	query.Add("api_token", client.ApiToken)
	for key, elem := range params {
		query.Add(key, elem)
	}

	url := fmt.Sprintf("%v/%v?%v", client.BaseUrl, path, query.Encode())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bytes, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode == 200 && strings.HasPrefix(res.Header.Get("Content-Type"), "text/html") {
		bytes, err := json.Marshal(Response{Success: false})
		if err != nil {
			return nil, err
		}
		msg := json.RawMessage(bytes)
		return &msg, nil
	}

	msg := json.RawMessage(bytes)
	return &msg, nil
}

func (client *Client) ListServers() (*ServersListResponse, error) {
	raw, err := client.get("list", nil)
	if err != nil {
		return nil, err
	}

	var res ServersListResponse
	if err := json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (client *Client) StopServer(server string) (*Response, error) {
	raw, err := client.get("stop/single", map[string]string{"server": server})
	if err != nil {
		return nil, err
	}

	var res Response
	if err := json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (client *Client) StartServer(server string) (*Response, error) {
	raw, err := client.get("start/single", map[string]string{"server": server})
	if err != nil {
		return nil, err
	}

	var res Response
	if err := json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (client *Client) GetServer(server string) (*ServersGetSingleResponse, error) {
	raw, err := client.get("get/single", map[string]string{"server": server})
	if err != nil {
		return nil, err
	}

	var res ServersGetSingleResponse
	if err := json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (client *Client) GetBillingDetails() (*BillingDetailsResponse, error) {
	raw, err := client.get("billing", nil)
	if err != nil {
		return nil, err
	}

	var res BillingDetailsResponse
	if err := json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func NewClient(baseUrl string, apiKey string, apiToken string) *Client {
	return &Client{baseUrl, apiKey, apiToken}
}
