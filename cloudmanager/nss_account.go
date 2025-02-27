package cloudmanager

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/fatih/structs"
)

type nssAccountRequest struct {
	AccountCredentials nssAccountCredentials `structs:"providerKeys"`
	VsaList            []string              `structs:"vsaList"`
}

type nssAccountCredentials struct {
	Username string `structs:"nssUserName"`
	Password string `structs:"nssPassword"`
}

func (c *Client) createNssAccount(acc nssAccountRequest, clientID string) (map[string]interface{}, error) {
	hostType := "CloudManagerHost"
	baseURL := fmt.Sprint("/occm/api/accounts/nss")
	param := structs.Map(acc)
	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			return nil, err
		}
		c.Token = accesTokenResult.Token
	}
	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, param, c.Token, hostType, clientID)
	if err != nil {
		log.Print("createNssAccount request failed ", statusCode)
		return nil, err
	}
	responseError := apiResponseChecker(statusCode, response, "createNssAccount")
	if responseError != nil {
		return nil, responseError
	}
	var res map[string]interface{}
	if err := json.Unmarshal(response, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) getNssAccount(nssUserName string, clientID string) (map[string]interface{}, error) {
	hostType := "CloudManagerHost"
	baseURL := fmt.Sprint("/occm/api/accounts")
	if c.Token == "" {
		accessTokenResult, err := c.getAccessToken()
		if err != nil {
			return nil, err
		}
		c.Token = accessTokenResult.Token
	}
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("getNssAccount request failed ", statusCode)
		return nil, err
	}
	responseError := apiResponseChecker(statusCode, response, "getNssAccount")

	if responseError != nil {
		return nil, responseError
	}
	var allAccounts map[string]interface{}
	if err := json.Unmarshal(response, &allAccounts); err != nil {
		return nil, err
	}
	var res []interface{} = allAccounts["nssAccounts"].([]interface{})
	for _, acc := range res {
		info := acc.(map[string]interface{})
		if info["nssUserName"] == nssUserName {
			return info, nil
		}
	}
	return nil, nil
}

func (c *Client) deleteNssAccount(id string, clientID string) error {
	hostType := "CloudManagerHost"
	baseURL := fmt.Sprintf("/occm/api/accounts/%s", id)
	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			return err
		}
		c.Token = accesTokenResult.Token
	}
	statusCode, response, _, err := c.CallAPIMethod("DELETE", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("deleteNssAccount request failed ", statusCode)
		return err
	}
	responseError := apiResponseChecker(statusCode, response, "deleteNssAccount")
	if responseError != nil {
		return responseError
	}
	return nil
}
