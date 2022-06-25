package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
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

type DeployServerResponse struct {
	Response
	Server ServerDeploy `json:"server"`
}

type Client struct {
	BaseUrl  string
	ApiKey   string
	ApiToken string
	Debug    bool
}

func (client *Client) do(method string, path string, params map[string]string, headers map[string]string, body []byte) (*json.RawMessage, error) {
	query := url.Values{}
	for key, elem := range params {
		query.Add(key, elem)
	}

	url := fmt.Sprintf("%v/%v?%v", client.BaseUrl, path, query.Encode())
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	for key, elem := range headers {
		req.Header.Add(key, elem)
	}

	if client.Debug {
		reqDump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		fmt.Println(string(reqDump))
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if client.Debug {
		resDump, err := httputil.DumpResponse(res, true)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		fmt.Println(string(resDump))
	}

	bytes, _ := ioutil.ReadAll(res.Body)
	if strings.HasPrefix(res.Header.Get("Content-Type"), "text/html") {
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

func (client *Client) get(path string, params map[string]string, auth bool) (*json.RawMessage, error) {
	newParams := map[string]string{}

	if auth {
		newParams["api_key"] = client.ApiKey
		newParams["api_token"] = client.ApiToken
	}

	for key, elem := range params {
		newParams[key] = elem
	}

	return client.do(http.MethodGet, path, newParams, nil, nil)
}

func (client *Client) post(path string, body map[string]string, auth bool) (*json.RawMessage, error) {
	newBody := url.Values{}

	if auth {
		newBody.Add("api_key", client.ApiKey)
		newBody.Add("api_token", client.ApiToken)
	}

	for key, elem := range body {
		newBody.Add(key, elem)
	}

	return client.do(
		http.MethodPost,
		path,
		nil,
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		[]byte(newBody.Encode()),
	)
}

func (client *Client) ListServers() (*ServersListResponse, error) {
	raw, err := client.get("list", nil, true)
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
	raw, err := client.get("stop/single", map[string]string{"server": server}, true)
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
	raw, err := client.get("start/single", map[string]string{"server": server}, true)
	if err != nil {
		return nil, err
	}

	var res Response
	if err := json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (client *Client) DeleteServer(server string) (*Response, error) {
	raw, err := client.get("delete/single", map[string]string{"server": server}, true)
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
	raw, err := client.get("get/single", map[string]string{"server": server}, true)
	if err != nil {
		return nil, err
	}

	var res ServersGetSingleResponse
	if err := json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (client *Client) DeployServer(adminUser string, adminPass string, instanceType string, gpuModel string, gpuCount int, vcpus int, ram int, storage int, storageClass string, os string, location string, name string) (*DeployServerResponse, error) {
	body := map[string]string{
		"admin_user":    adminUser,
		"admin_pass":    adminPass,
		"instance_type": instanceType,
		"gpu_model":     gpuModel,
		"gpu_count":     strconv.Itoa(gpuCount),
		"vcpus":         strconv.Itoa(vcpus),
		"ram":           strconv.Itoa(ram),
		"storage":       strconv.Itoa(storage),
		"storage_class": storageClass,
		"os":            os,
		"location":      location,
		"name":          name,
	}
	raw, err := client.post("deploy/single/custom", body, true)
	if err != nil {
		return nil, err
	}

	var res DeployServerResponse
	if err := json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (client *Client) GetBillingDetails() (*BillingDetailsResponse, error) {
	raw, err := client.get("billing", nil, true)
	if err != nil {
		return nil, err
	}

	var res BillingDetailsResponse
	if err := json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func NewClient(baseUrl string, apiKey string, apiToken string, debug bool) *Client {
	return &Client{baseUrl, apiKey, apiToken, debug}
}

func (client *Client) RestartServer(server string) (*Response, error) {
	raw, err := client.get("restart/single", map[string]string{"server": server}, true)
	if err != nil {
		return nil, err
	}

	var res Response
	if err := json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
