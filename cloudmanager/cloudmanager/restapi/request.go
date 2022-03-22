package restapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	GCPServiceAccountKey  string
}

func getGCPToken(url string, gcpServiceAccountKey string) (string, error) {
	var token string
	scopes := []string{
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/compute",
		"https://www.googleapis.com/auth/compute.readonly",
		"https://www.googleapis.com/auth/ndev.cloudman",
		"https://www.googleapis.com/auth/ndev.cloudman.readonly",
		"https://www.googleapis.com/auth/devstorage.full_control",
		"https://www.googleapis.com/auth/devstorage.read_write",
	}
	if gcpServiceAccountKey != "" {
		var c = struct {
			Email      string `json:"client_email"`
			PrivateKey string `json:"private_key"`
		}{}
		json.Unmarshal([]byte(gcpServiceAccountKey), &c)
		config := &jwt.Config{
			Email:      c.Email,
			PrivateKey: []byte(c.PrivateKey),
			Scopes:     scopes,
			TokenURL:   google.JWTTokenURL,
		}
		gcpToken, err := config.TokenSource(oauth2.NoContext).Token()
		if err != nil {
			return "", err
		}
		token = gcpToken.AccessToken
	} else {
		// find default application credential
		ctx := context.Background()
		credential, err := google.FindDefaultCredentials(ctx, scopes...)
		if err != nil {
			return "", fmt.Errorf("cannot get credentials: %v", err)
		}
		t, err := credential.TokenSource.Token()
		if err != nil {
			return "", fmt.Errorf("getGCPToken failed on get token from credential: %v", err)
		}
		token = t.AccessToken
	}

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
			if paramsNil {
				req, err = http.NewRequest(r.Method, url, bytes.NewReader([]byte(r.GCPDeploymentTemplate)))
				if err != nil {
					return nil, err
				}
			} else {
				bodyJSON, err := json.Marshal(r.Params)
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

		token, err = getGCPToken(url, r.GCPServiceAccountKey)
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
