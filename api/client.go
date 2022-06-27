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
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type Response struct {
	Success bool `json:"success"`
}

type GetServerResponse struct {
	Response
	Server Server `json:"server"`
}

type ListServersResponse struct {
	Response
	Servers map[string]Server `json:"servers"`
}

type GetBillingDetailsResponse struct {
	Response
	BillingDetails
}

type DeployServerRequest struct {
	AdminUser    string `mapstructure:"admin_user"`
	AdminPass    string `mapstructure:"admin_pass"`
	InstanceType string `mapstructure:"instance_type"`
	GPUModel     string `mapstructure:"gpu_model"`
	GPUCount     int    `mapstructure:"gpu_count"`
	VCPUs        int    `mapstructure:"vcpus"`
	RAM          int    `mapstructure:"ram"`
	Storage      int    `mapstructure:"storage"`
	StorageClass string `mapstructure:"storage_class"`
	OS           string `mapstructure:"os"`
	Location     string `mapstructure:"location"`
	Name         string `mapstructure:"name"`
}

type DeployServerResponse struct {
	Response
	Server struct {
		Id    string              `json:"id"`
		Ip    string              `json:"ip"`
		Links []map[string]string `json:"links"`
	} `json:"server"`
}

type ModifyServerRequest struct {
	ServerId     *string `mapstructure:"server_id,omitempty"`
	InstanceType *string `mapstructure:"instance_type,omitempty"`
	GPUModel     *string `mapstructure:"gpu_model,omitempty"`
	GPUCount     *int    `mapstructure:"gpu_count,omitempty"`
	VCPUs        *int    `mapstructure:"vcpus,omitempty"`
	RAM          *int    `mapstructure:"ram,omitempty"`
	Storage      *int    `mapstructure:"storage,omitempty"`
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

	// FIXME: Workaround for API issue which causes endpoint to
	// return an HTML Page with a 200 Status code
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

func (client *Client) ListServers() (*ListServersResponse, error) {
	raw, err := client.get("list", nil, true)
	if err != nil {
		return nil, err
	}

	var res ListServersResponse
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

func (client *Client) GetServer(server string) (*GetServerResponse, error) {
	raw, err := client.get("get/single", map[string]string{"server": server}, true)
	if err != nil {
		return nil, err
	}

	var res GetServerResponse
	if err := json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (client *Client) DeployServer(req DeployServerRequest) (*DeployServerResponse, error) {
	var rawBody map[string]interface{}
	err := mapstructure.Decode(req, &rawBody)
	if err != nil {
		return nil, err
	}

	// convert to map[string]string
	body := map[string]string{}
	for key, elem := range rawBody {
		str := fmt.Sprintf("%v", elem)
		body[key] = str
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

func (client *Client) GetBillingDetails() (*GetBillingDetailsResponse, error) {
	raw, err := client.get("billing", nil, true)
	if err != nil {
		return nil, err
	}

	var res GetBillingDetailsResponse
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

func (client *Client) ModifyServer(req ModifyServerRequest) (*Response, error) {
	var rawBody map[string]interface{}
	err := mapstructure.Decode(req, &rawBody)
	if err != nil {
		return nil, err
	}

	// conver to map[string]string resolving pointers as necessary
	body := map[string]string{}
	for key, elem := range rawBody {
		val := reflect.ValueOf(elem)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		body[key] = fmt.Sprintf("%v", val)
	}

	raw, err := client.post("modify/single/custom", body, true)
	if err != nil {
		return nil, err
	}

	var res Response
	if err := json.Unmarshal(*raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
