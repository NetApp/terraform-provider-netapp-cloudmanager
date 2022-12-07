package restapi

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

// Client represents a client for interaction with a CloudManager API
type Client struct {
	CloudManagerHost     string
	AuthHost             string
	SaAuthHost           string
	CVOHostName          string
	RefreshToken         string
	SaClientID           string
	SaSecretKey          string
	Audience             string
	GCPDeploymentManager string
	CVSHostName          string
	GCPCompute           string

	httpClient http.Client
}

// Do sends the API Request, parses the response as JSON, and returns the HTTP status code as int, onCloudRequestID from header as string, the "result" value as byte
func (c *Client) Do(baseURL string, hostType string, token string, paramsNil bool, accountID string, clientID string, req *Request, simulator bool) (int, []byte, string, error) {

	var host string
	var res []byte
	statusCode := 0
	res = nil
	onCloudRequestID := ""
	gcpType := false

	if hostType == "CloudManagerHost" {
		host = c.CloudManagerHost
	} else if hostType == "AuthHost" {
		host = c.AuthHost
		accountID = ""
	} else if hostType == "SaAuthHost" {
		host = c.SaAuthHost
		accountID = ""
	} else if hostType == "GCPDeploymentManager" {
		host = c.GCPDeploymentManager
		gcpType = true
	} else if hostType == "CVSHost" {
		host = c.CVSHostName
	} else if hostType == "GCPCompute" {
		host = c.GCPCompute
		gcpType = true
	}

	httpReq, err := req.BuildHTTPReq(host, token, c.Audience, baseURL, paramsNil, accountID, clientID, gcpType, simulator)
	if err != nil {
		return statusCode, res, onCloudRequestID, err
	}
	httpRes, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Print("HTTP req failed")
		return statusCode, res, onCloudRequestID, err
	}

	if httpRes.Header.Get("OnCloud-Request-Id") != "" {
		log.Print("OnCloud-Request-Id ", httpRes.Header.Get("OnCloud-Request-Id"))
		onCloudRequestID = httpRes.Header.Get("OnCloud-Request-Id")
	}

	defer httpRes.Body.Close()

	res, err = ioutil.ReadAll(httpRes.Body)
	if err != nil {
		log.Print("HTTP decoder failed")
		return statusCode, res, onCloudRequestID, err
	}

	if res == nil {
		return statusCode, res, onCloudRequestID, errors.New("No result returned in REST response")
	}

	statusCode = httpRes.StatusCode
	return statusCode, res, onCloudRequestID, nil
}
