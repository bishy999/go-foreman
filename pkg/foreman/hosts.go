/*
Package foreman provide a client to the Foreman API.
*/
package foreman

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

const (
	hostsapi = "api/hosts"
	message  = "Resource host not found by i"
)

// HostsReq contains parent field for data payload for creating a host
type hostsReq struct {
	Host hostsReqBody `json:"host"`
}

// hostsReqBody contains fields neccessary for creating a host
type hostsReqBody struct {
	Name                    string            `json:"name"`
	HostgroupID             int               `json:"hostgroup_id"`
	OrganisationID          int               `json:"organization_id"`
	LocationID         		int               `json:"location_id"`
	Managed                 bool              `json:"managed"`
	ComputeProfileID        string            `json:"compute_profile_id"`
	ProvisionMethod         string            `json:"provision_method"`
	Build                   bool              `json:"build"`
	Enabled                 bool              `json:"enabled"`
	Comment                 string            `json:"comment"`
	Overwrite               string            `json:"overwrite"`
	ComputeAttributes       map[string]string `json:"compute_attributes"`
	HostParameterAttributes interface{}       `json:"host_parameters_attributes"`
}

// CheckHost checks if instance already exists
func (ci *ConnectionInfo) CheckHost(ctx context.Context) (bool, string, error) {

	var jsonData []byte
	url, _ := url.Parse(ci.BaseURL)
	url.Path = path.Join(url.Path, hostsapi, ci.Hostname)
	apiused := url.String()
	log.Printf("API Used: %s", apiused)

	exists, status, err := sendRequest(ctx, ci.Hostname, ci.Username, ci.Password, ci.Client, jsonData, apiused, http.MethodGet)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	return exists, status, err
}

// CreateHost create host with the name provided
func (ci *ConnectionInfo) CreateHost(ctx context.Context) (bool, string, error) {

	requestBody := hostsReqBody{
		Name:                    ci.Hostname,
		HostgroupID:             ci.Group,
		OrganisationID:			 9,
		LocationID:              15,
		Managed:                 true,
		ComputeProfileID:        ci.Profile,
		ProvisionMethod:         "image",
		Build:                   true,
		Enabled:                 true,
		Comment:                 "Built by Jenkins",
		Overwrite:               "false",
		ComputeAttributes:       map[string]string{"flavor_id": ci.Size},
		HostParameterAttributes: []interface{}{map[string]string{"name": "disksize", "value": "512"}},
	}

	requestHead := hostsReq{
		Host: requestBody,
	}

	jsonData, err := json.MarshalIndent(requestHead, "", "    ")
	log.Printf("%s", string(jsonData))
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	url, _ := url.Parse(ci.BaseURL)
	url.Path = path.Join(url.Path, hostsapi)
	apiused := url.String()
	log.Printf("API Used:(%s) %s", http.MethodPost, apiused)

	exists, _, err := sendRequest(ctx, ci.Hostname, ci.Username, ci.Password, ci.Client, jsonData, apiused, http.MethodPost)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	status := "The host [" + ci.Hostname + "] was created successfully"

	return exists, status, err

}

// DeleteHost deletes the host with name provided
func (ci *ConnectionInfo) DeleteHost(ctx context.Context) (bool, string, error) {

	var jsonData []byte
	url, _ := url.Parse(ci.BaseURL)
	url.Path = path.Join(url.Path, hostsapi, ci.Hostname)
	apiused := url.String()
	log.Printf("API Used: %s", apiused)

	exists, status, err := sendRequest(ctx, ci.Hostname, ci.Username, ci.Password, ci.Client, jsonData, apiused, http.MethodDelete)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	return exists, status, err
}

// sendRequest send http request to specified endpoints and returns response
func sendRequest(ctx context.Context, name string, user string, password string, client *http.Client, data []byte, api string, method string) (bool, string, error) {

	var status string
	var err error
	var req *http.Request
	var exists = true

	if method == http.MethodGet {
		req, err = http.NewRequest(http.MethodGet, api, nil)

	} else if method == http.MethodDelete {
		req, err = http.NewRequest(http.MethodDelete, api, nil)

	} else {
		req, err = http.NewRequest(http.MethodPost, api, bytes.NewBuffer(data))
	}
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	req = req.WithContext(ctx)
	req.SetBasicAuth(user, password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	} else {
		status = string(body)
		if strings.Contains(status, message) {
			exists = false
		}
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 404 {
		log.Printf("Http status code is: %v", resp.StatusCode)
		log.Printf("Http request sent was: %v", resp.Request)
	}

	return exists, status, err
}