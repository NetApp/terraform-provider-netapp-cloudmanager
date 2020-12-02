package cloudmanager

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fatih/structs"
)

// createUserData the users input for creating a occm
type createUserData struct {
	ClientID      string                `json:"clientId"`
	ClientSecret  string                `json:"clientSecret"`
	UUID          string                `json:"systemId"`
	AccountID     string                `json:"tenancyAccountId"`
	Company       string                `json:"company"`
	Name          string                `json:"instanceName"`
	ProxySettings proxySettingsResponse `json:"proxySettings"`
}

type proxySettingsResponse struct {
	ProxyURL      string `json:"proxyUrl"`
	ProxyUserName string `json:"proxyUserName"`
	ProxyPassword string `json:"proxyPassword"`
}

// createOCCMDetails the users input for creating a occm
type createOCCMDetails struct {
	Name                         string
	GCPProject                   string
	Company                      string
	InstanceID                   string
	Region                       string
	Location                     string
	Zone                         string
	AMI                          string
	KeyName                      string
	InstanceType                 string
	IamInstanceProfileName       string
	SecurityGroupID              string
	SubnetID                     string
	NetworkProjectID             string
	ProxyURL                     string
	ProxyUserName                string
	ProxyPassword                string
	ResourceGroup                string
	SubscriptionID               string
	MachineType                  string
	ServiceAccountEmail          string
	GCPCommonSuffixName          string
	VnetID                       string
	VnetResourceGroup            string
	AdminUsername                string
	AdminPassword                string
	VirtualMachineSize           string
	NetworkSecurityGroupName     string
	NetworkSecurityResourceGroup string
	AssociatePublicIPAddress     *bool
	AssociatePublicIP            bool
	FirewallTags                 bool
}

// deleteOCCMDetails the users input for deleting a occm
type deleteOCCMDetails struct {
	InstanceID          string
	Name                string
	SubscriptionID      string
	ResourceGroup       string
	Location            string
	Region              string
	Project             string
	GCPCommonSuffixName string
}

// OCCMMResult the users input for creating a occm
type OCCMMResult struct {
	ClientID   string
	AccountID  string
	InstanceID string
}

// accesTokenRequest the input for requesting a token
type accesTokenRequest struct {
	Audience     string `structs:"audience"`
	GrantType    string `structs:"grant_type"`
	RefreshToken string `structs:"refresh_token"`
	ClientID     string `structs:"client_id"`
}

// accesTokenResult to get token for the AUTH
type accesTokenResult struct {
	Token string `json:"access_token"`
}

// registerAgentTOServiceRequest input to register agent
type registerAgentTOServiceRequest struct {
	AccountID string           `structs:"accountId"`
	Name      string           `structs:"name"`
	Company   string           `structs:"company"`
	Placement placementRequest `structs:"placement"`
	Extra     extraRequest     `structs:"extra"`
}

// placementRequest input to register agent
type placementRequest struct {
	Subnet   string `structs:"subnet"`
	Provider string `structs:"provider"`
	Region   string `structs:"region"`
	Network  string `structs:"network"`
}

// extraRequest structure for the proxy credentials
type extraRequest struct {
	Proxy proxyRequest `structs:"proxy,omitempty"`
}

// proxyRequest the user input for using proxy credentials
type proxyRequest struct {
	ProxyURL      string `structs:"proxyUrl,omitempty"`
	ProxyUserName string `structs:"proxyUserName,omitempty"`
	ProxyPassword string `structs:"proxyPassword,omitempty"`
}

// accountResult lists account to get the account ID
type accountResult struct {
	Account []accountIDResult `json:"AnyValue"`
}

// accountIDResult to get the account ID
type accountIDResult struct {
	AccountID string `json:"accountPublicId"`
}

// listOCCMResult lists the details for given Client ID
type listOCCMResult struct {
	Agent occmAgent `json:"agent"`
}

// accountIDCreate to create the account ID
type accountIDCreate struct {
	Name string `structs:"name"`
}

// occmAgent lists the listOCCMResult details for given Client ID
type occmAgent struct {
	Status  string `json:"status"`
	AgentID string `json:"agentId"`
}

func (c *Client) getUserData(registerAgentTOService registerAgentTOServiceRequest) (string, error) {
	accesTokenResult, err := c.getAccessToken()
	if err != nil {
		return "", err
	}
	log.Print("getAccessToken  ", accesTokenResult.Token)
	c.Token = accesTokenResult.Token

	if c.AccountID == "" {
		accountID, err := c.getAccount()
		if err != nil {
			return "", err
		}
		log.Print("getAccount ", accountID)
		registerAgentTOService.AccountID = accountID
	} else {
		registerAgentTOService.AccountID = c.AccountID
	}

	userDataRespone, err := c.registerAgentTOService(registerAgentTOService)
	if err != nil {
		return "", err
	}

	c.ClientID = userDataRespone.ClientID
	c.AccountID = userDataRespone.AccountID

	userData := "{\n\t\"instanceName\": \"" + userDataRespone.Name + "\",\n\t\"company\": \"" + userDataRespone.Company + "\",\n\t\"clientId\": \"" + userDataRespone.ClientID + "\",\n\t\"clientSecret\": \"" + userDataRespone.ClientSecret + "\",\n\t\"systemId\": \"" + userDataRespone.UUID + "\",\n\t\"tenancyAccountId\": \"" + userDataRespone.AccountID + "\",\n\t\"proxySettings\": {\n\t\"proxyPassword\": \"" + userDataRespone.ProxySettings.ProxyPassword + "\",\n\t\"proxyUserName\": \"" + userDataRespone.ProxySettings.ProxyUserName + "\",\n\t\"proxyUrl\": \"" + userDataRespone.ProxySettings.ProxyURL + "\"\n}\n}"
	log.Print("userData ", userData)

	return userData, nil
}

func (c *Client) getAccessToken() (accesTokenResult, error) {

	log.Print("getAccessToken")

	hostType := "AuthHost"
	var accesTokenRequest accesTokenRequest

	accesTokenRequest.Audience = c.Audience
	accesTokenRequest.GrantType = "refresh_token"
	accesTokenRequest.RefreshToken = c.RefreshToken
	accesTokenRequest.ClientID = c.Auth0Client

	params := structs.Map(accesTokenRequest)
	statusCode, response, _, err := c.CallAPIMethod("POST", "", params, "", hostType)
	if err != nil {
		log.Print("getAccessToken request failed ", statusCode)
		return accesTokenResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "getAccessToken")
	if responseError != nil {
		return accesTokenResult{}, responseError
	}

	var result accesTokenResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getAccessToken ", err)
		return accesTokenResult{}, err
	}

	return result, nil
}

func (c *Client) registerAgentTOService(registerAgentTOServiceRequest registerAgentTOServiceRequest) (createUserData, error) {

	baseURL := "/agents-mgmt/connector-setup"
	hostType := "CloudManagerHost"

	vpcID, err := c.CallVPCGet(registerAgentTOServiceRequest.Placement.Subnet, registerAgentTOServiceRequest.Placement.Region)
	if err != nil {
		log.Print("CallVPCGet request failed")
		return createUserData{}, err
	}

	registerAgentTOServiceRequest.Placement.Network = vpcID
	registerAgentTOServiceRequest.Placement.Provider = "AWS"

	params := structs.Map(registerAgentTOServiceRequest)
	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, params, c.Token, hostType)
	if err != nil {
		log.Print("registerAgentTOService request failed ", statusCode)
		return createUserData{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "registerAgentTOService")
	if responseError != nil {
		return createUserData{}, responseError
	}

	var result createUserData
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from registerAgentTOService ", err)
		return createUserData{}, err
	}

	return result, nil
}

func (c *Client) getAccount() (string, error) {

	log.Print("getAccount")

	baseURL := "/tenancy/account"

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Print("getAccount request failed ", statusCode)
		return "", err
	}

	responseError := apiResponseChecker(statusCode, response, "getAccount")
	if responseError != nil {
		return "", responseError
	}

	var result []accountIDResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getAccount ", err)
		return "", err
	}

	// when no account exists, create
	if len(result) == 0 {
		accountID, err := c.createAccount()
		if err != nil {
			log.Print("createAccount request failed")
			return "", err
		}
		return accountID, nil
	}

	return result[0].AccountID, nil
}

func (c *Client) createAccount() (string, error) {

	log.Print("createAccount")

	baseURL := "/tenancy/account/MyAccount"

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Print("createAccount request failed ", statusCode)
		return "", err
	}

	responseError := apiResponseChecker(statusCode, response, "createAccount")
	if responseError != nil {
		return "", responseError
	}

	var result accountIDResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from createAccount ", err)
		return "", err
	}

	return result.AccountID, nil
}

func (c *Client) createAWSInstance(occmDetails createOCCMDetails) (string, error) {

	instanceID, err := c.CallAWSInstanceCreate(occmDetails)
	if err != nil {
		return "", err
	}

	log.Print("Sleep for 2 minutes")
	time.Sleep(time.Duration(120) * time.Second)

	retries := 16
	for {
		occmResp, err := c.checkOCCMStatus()
		if err != nil {
			return "", err
		}
		if occmResp.Status == "active" {
			break
		} else {
			if retries == 0 {
				log.Print("Taking too long for status to be active")
				return "", fmt.Errorf("Taking too long for OCCM agent to be active or not properly setup")
			}
			time.Sleep(time.Duration(30) * time.Second)
			retries--
		}
	}

	return instanceID, nil
}

func (c *Client) getAWSInstance(occmDetails createOCCMDetails, id string) (string, error) {

	log.Print("getAWSInstance")
	if occmDetails.AMI == "" {
		var err error
		occmDetails.AMI, err = c.CallAMIGet(occmDetails)
		if err != nil {
			return "", err
		}
	}

	res, err := c.CallAWSInstanceGet(occmDetails)
	if err != nil {
		return "", err
	}

	for _, instanceID := range res {
		if instanceID == id {
			return instanceID, nil
		}
	}

	return "", nil
}

func (c *Client) createOCCM(occmDetails createOCCMDetails) (OCCMMResult, error) {

	log.Print("createOCCM")
	if occmDetails.AMI == "" {

		ami, err := c.CallAMIGet(occmDetails)
		if err != nil {
			return OCCMMResult{}, err
		}
		occmDetails.AMI = ami
	}

	var registerAgentTOService registerAgentTOServiceRequest
	registerAgentTOService.Name = occmDetails.Name
	registerAgentTOService.Placement.Region = occmDetails.Region
	registerAgentTOService.Placement.Subnet = occmDetails.SubnetID
	registerAgentTOService.Company = occmDetails.Company
	if occmDetails.ProxyURL != "" {
		registerAgentTOService.Extra.Proxy.ProxyURL = occmDetails.ProxyURL
	}

	if occmDetails.ProxyUserName != "" {
		registerAgentTOService.Extra.Proxy.ProxyUserName = occmDetails.ProxyUserName
	}

	if occmDetails.ProxyPassword != "" {
		registerAgentTOService.Extra.Proxy.ProxyPassword = occmDetails.ProxyPassword
	}

	userData, err := c.getUserData(registerAgentTOService)
	if err != nil {
		return OCCMMResult{}, err
	}
	c.UserData = userData
	instanceID, err := c.createAWSInstance(occmDetails)
	if err != nil {
		return OCCMMResult{}, err
	}

	var result OCCMMResult
	result.InstanceID = instanceID
	result.ClientID = c.ClientID
	result.AccountID = c.AccountID

	return result, nil
}

func (c *Client) checkOCCMStatus() (occmAgent, error) {

	baseURL := fmt.Sprintf("/agents-mgmt/agent/%sclients", c.ClientID)

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Print("checkOCCMStatus request failed ", statusCode)
		return occmAgent{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "checkOCCMStatus")
	if responseError != nil {
		return occmAgent{}, responseError
	}

	var result listOCCMResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from checkOCCMStatus ", err)
		return occmAgent{}, err
	}

	return result.Agent, nil
}

func (c *Client) callOCCMDelete() error {

	baseURL := fmt.Sprintf("/agents-mgmt/agent/%sclients", c.ClientID)

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("DELETE", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Print("callOCCMDelete request failed ", statusCode)
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "callOCCMDelete")
	if responseError != nil {
		return responseError
	}

	return nil
}

func (c *Client) deleteOCCM(request deleteOCCMDetails) error {

	err := c.CallAWSInstanceTerminate(request)
	if err != nil {
		return err
	}

	log.Print("Sleep for 30 seconds")
	time.Sleep(time.Duration(30) * time.Second)

	accesTokenResult, err := c.getAccessToken()
	c.Token = accesTokenResult.Token
	if err != nil {
		return err
	}

	retries := 30
	for {
		occmResp, err := c.checkOCCMStatus()
		if err != nil {
			return err
		}
		if occmResp.Status != "active" {
			break
		} else {
			if retries == 0 {
				log.Print("Taking too long for instance to finish terminating")
				return fmt.Errorf("Taking too long for instance to finish terminating")
			}
			time.Sleep(time.Duration(10) * time.Second)
			retries--
		}
	}

	if err := c.callOCCMDelete(); err != nil {
		return err
	}

	return nil
}
