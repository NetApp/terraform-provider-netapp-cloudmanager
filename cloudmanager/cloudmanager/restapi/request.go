package restapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Request represents a request to a REST API
type Request struct {
	Method                string      `json:"method"`
	Params                interface{} `json:"params"`
	GCPDeploymentTemplate string
	GCPServiceAccountKey  string
}

// BuildHTTPReq builds an HTTP request to carry out the REST request
func (r *Request) BuildHTTPReq(host string, token string, audience string, baseURL string, paramsNil bool, accountID string, clientID string, gcpType bool, cloudmanagerSimulator bool) (*http.Request, error) {

	url := host + baseURL
	var req *http.Request
	var err error

	// authenticating separately for GCP calls
	if gcpType {
		if r.Method == "POST" {
			if paramsNil {
				req, err = http.NewRequest(r.Method, url, bytes.NewReader([]byte(r.GCPDeploymentTemplate)))
				if err != nil {
					return nil, err
				}
			} else {
				bodyJSON, err := json.Marshal(r.Params)
				if err != nil {
					return nil, err
				}
				req, err = http.NewRequest(r.Method, url, bytes.NewReader([]byte(bodyJSON)))
				if err != nil {
					return nil, err
				}
			}
		} else {
			var err error
			req, err = http.NewRequest(r.Method, url, nil)
			if err != nil {
				return nil, err
			}
		}
		if token == "" {
			return nil, fmt.Errorf("no GCP token available")
		}
	} else {
		if paramsNil {
			var err error
			req, err = http.NewRequest(r.Method, url, nil)
			if err != nil {
				return nil, err
			}
		} else {
			bodyJSON, err := json.Marshal(r.Params)
			if err != nil {
				return nil, err
			}
			req, err = http.NewRequest(r.Method, url, bytes.NewReader(bodyJSON))
			if err != nil {
				return nil, err
			}
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "Terraform_NetApp")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if accountID != "" {
		log.Print(" inside account TOKEN")
		if token != "" {
			req.Header.Set("X-User-Token", "Bearer "+token)
		}
		req.Header.Set("X-Tenancy-Account-Id", accountID)
	}
	if clientID != "" {
		if strings.HasSuffix(clientID, "clients") {
			req.Header.Set("X-Agent-Id", clientID)
		} else {
			req.Header.Set("X-Agent-Id", clientID+"clients")
		}
	}
	// AWS FSx
	if cloudmanagerSimulator {
		req.Header.Set("x-simulator", "true")
		log.Println("Running on simulator")
	}

	return req, nil
}
