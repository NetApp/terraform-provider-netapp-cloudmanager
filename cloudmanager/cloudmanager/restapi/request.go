package restapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

// Request represents a request to a REST API
type Request struct {
	Method                string      `json:"method"`
	Params                interface{} `json:"params"`
	GCPDeploymentTemplate string
	GCPServiceAccountPath string
}

func getGCPToken(url string, gcpServiceAccountPath string) (string, error) {

	keyBytes, err := ioutil.ReadFile(gcpServiceAccountPath)
	if err != nil {
		return "", fmt.Errorf("Unable to read service account key file  %v", err)
	}

	var c = struct {
		Email      string `json:"client_email"`
		PrivateKey string `json:"private_key"`
	}{}
	json.Unmarshal(keyBytes, &c)
	config := &jwt.Config{
		Email:      c.Email,
		PrivateKey: []byte(c.PrivateKey),
		Scopes: []string{
			"https://www.googleapis.com/auth/cloud-platform",
			"https://www.googleapis.com/auth/compute",
			"https://www.googleapis.com/auth/compute.readonly",
			"https://www.googleapis.com/auth/ndev.cloudman",
			"https://www.googleapis.com/auth/ndev.cloudman.readonly",
			"https://www.googleapis.com/auth/devstorage.full_control",
			"https://www.googleapis.com/auth/devstorage.read_write",
		},
		TokenURL: google.JWTTokenURL,
	}
	gcpToken, err := config.TokenSource(oauth2.NoContext).Token()
	if err != nil {
		return "", err
	}
	token := gcpToken.AccessToken

	return token, nil
}

// BuildHTTPReq builds an HTTP request to carry out the REST request
func (r *Request) BuildHTTPReq(host string, token string, audience string, baseURL string, paramsNil bool, accountID string, clientID string, gcpType bool, cloudmanagerSimulator bool) (*http.Request, error) {

	url := host + baseURL
	var req *http.Request
	var err error

	// authenticating separately for GCP calls
	if gcpType == true {
		if r.Method == "POST" {
			req, err = http.NewRequest(r.Method, url, bytes.NewReader([]byte(r.GCPDeploymentTemplate)))
			if err != nil {
				return nil, err
			}
		} else {
			var err error
			req, err = http.NewRequest(r.Method, url, nil)
			if err != nil {
				return nil, err
			}
		}

		token, err = getGCPToken(url, r.GCPServiceAccountPath)
		if err != nil {
			return nil, err
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
