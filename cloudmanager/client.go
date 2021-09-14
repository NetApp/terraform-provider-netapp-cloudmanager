package cloudmanager

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-07-01/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/netapp/terraform-provider-netapp-cloudmanager/cloudmanager/cloudmanager/restapi"
	"github.com/sirupsen/logrus"
)

var ourlog = logrus.WithFields(logrus.Fields{
	"prefix": "main",
})

// A Client to interact with the CVO/AWS/GCP/ANF/OCCM API
type Client struct {
	CloudManagerHost        string
	AuthHost                string
	SaAuthHost              string
	SaSecretKey             string
	SaClientID              string
	CVOHostName             string
	HostType                string
	MaxConcurrentRequests   int
	UserData                string
	BaseURL                 string
	RefreshToken            string
	Project                 string
	CloudManagerDomain      string
	Audience                string
	Auth0Client             string
	ClientID                string
	AccountID               string
	Token                   string
	AMIFilter               string
	AWSAccount              string
	AzureEnvironmentForOCCM string
	GCPDeploymentManager    string
	GCPImageProject         string
	GCPImageFamily          string
	GCPDeploymentTemplate   string
	GCPServiceAccountPath   string
	CVSHostName             string

	initOnce      sync.Once
	instanceInput *restapi.Client
	restapiClient *restapi.Client
	requestSlots  chan int
	Simulator     bool
}

// CallAWSInstanceCreate can be used to make a request to create AWS Instance
func (c *Client) CallAWSInstanceCreate(occmDetails createOCCMDetails) (string, error) {

	sess := session.Must(session.NewSession(aws.NewConfig().WithRegion(occmDetails.Region)))

	// Create EC2 service client
	svc := ec2.New(sess)

	var securityGroupIds []*string
	split := strings.Split(occmDetails.SecurityGroupID, ",")
	for _, sgid := range split {
		securityGroupIds = append(securityGroupIds, aws.String(sgid))
	}

	tags := []*ec2.Tag{}
	tag := &ec2.Tag{
		Key:   aws.String("Name"),
		Value: aws.String(occmDetails.Name),
	}
	tags = append(tags, tag)

	tag = &ec2.Tag{
		Key:   aws.String("OCCMInstance"),
		Value: aws.String("true"),
	}
	tags = append(tags, tag)

	if len(occmDetails.AwsTags) > 0 {
		for _, awsTag := range occmDetails.AwsTags {
			tag := &ec2.Tag{
				Key:   aws.String(awsTag.TagKey),
				Value: aws.String(awsTag.TagValue),
			}
			tags = append(tags, tag)
		}
	}
	// Specify the details of the instance that you want to create.
	runInstancesInput := &ec2.RunInstancesInput{
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/sda1"),
				Ebs: &ec2.EbsBlockDevice{
					Encrypted:  aws.Bool(true),
					VolumeSize: aws.Int64(100),
					VolumeType: aws.String("gp2"),
				},
			},
		},
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Name: aws.String(occmDetails.IamInstanceProfileName),
		},
		ImageId:               aws.String(occmDetails.AMI),
		InstanceType:          aws.String(occmDetails.InstanceType),
		MinCount:              aws.Int64(1),
		MaxCount:              aws.Int64(1),
		KeyName:               aws.String(occmDetails.KeyName),
		DisableApiTermination: aws.Bool(*occmDetails.EnableTerminationProtection),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags:         tags,
			},
		},
		UserData: aws.String(base64.StdEncoding.EncodeToString([]byte(c.UserData))),
	}

	if occmDetails.AssociatePublicIPAddress != nil {
		networkInterface := ec2.InstanceNetworkInterfaceSpecification{
			AssociatePublicIpAddress: aws.Bool(*occmDetails.AssociatePublicIPAddress),
			DeviceIndex:              aws.Int64(0),
			SubnetId:                 aws.String(occmDetails.SubnetID),
		}
		networkInterface.Groups = securityGroupIds
		runInstancesInput.NetworkInterfaces = []*ec2.InstanceNetworkInterfaceSpecification{&networkInterface}
	} else {
		// Network interfaces and an instance-level subnet ID may not be specified on the same request
		runInstancesInput.SubnetId = aws.String(occmDetails.SubnetID)

		// Network interfaces and an instance-level security groups may not be specified on the same request
		runInstancesInput.SecurityGroupIds = securityGroupIds
	}

	runResult, err := svc.RunInstances(runInstancesInput)

	if err != nil {
		log.Print("Could not create instance ", err)
		return "", err
	}

	log.Printf("Created instance %s", *runResult.Instances[0].InstanceId)

	return *runResult.Instances[0].InstanceId, nil
}

// CallAWSInstanceTerminate can be used to make a request to terminate AWS Instance
func (c *Client) CallAWSInstanceTerminate(occmDetails deleteOCCMDetails) error {

	sess := session.Must(session.NewSession(aws.NewConfig().WithRegion(occmDetails.Region)))

	// Create EC2 service client
	svc := ec2.New(sess)

	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(occmDetails.InstanceID),
		},
	}

	// Specify the details of the instance that you want to terminate.
	runResult, err := svc.TerminateInstances(input)
	if err != nil {
		log.Print("Could not terminate instance ", err)
		return err
	}

	log.Printf("Terminated instance %s", *runResult)

	return nil
}

// CallDeployAzureVM can be used to make a request to deploy Azure VM
func (c *Client) CallDeployAzureVM(occmDetails createOCCMDetails) (string, error) {

	var template *map[string]interface{}
	var params *map[string]interface{}

	json.Unmarshal([]byte(c.callTemplate()), &template)
	json.Unmarshal([]byte(c.callParameters()), &params)

	(*params)["adminPassword"] = map[string]string{
		"value": occmDetails.AdminPassword,
	}

	(*params)["customData"] = map[string]string{
		"value": c.UserData,
	}

	(*params)["virtualMachineName"] = map[string]string{
		"value": occmDetails.Name,
	}

	(*params)["location"] = map[string]string{
		"value": occmDetails.Location,
	}

	(*params)["adminUsername"] = map[string]string{
		"value": occmDetails.AdminUsername,
	}

	if c.AzureEnvironmentForOCCM == "stage" {
		(*params)["environment"] = map[string]string{
			"value": c.AzureEnvironmentForOCCM,
		}
	}

	var vnetID string
	var networkSecurityGroupName string

	if occmDetails.VnetResourceGroup != "" {
		vnetID = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s", occmDetails.SubscriptionID, occmDetails.VnetResourceGroup, occmDetails.VnetID)
	} else {
		vnetID = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s", occmDetails.SubscriptionID, occmDetails.ResourceGroup, occmDetails.VnetID)
	}

	if occmDetails.NetworkSecurityResourceGroup != "" {
		networkSecurityGroupName = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkSecurityGroups/%s", occmDetails.SubscriptionID, occmDetails.NetworkSecurityResourceGroup, occmDetails.NetworkSecurityGroupName)
	} else {
		networkSecurityGroupName = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkSecurityGroups/%s", occmDetails.SubscriptionID, occmDetails.ResourceGroup, occmDetails.NetworkSecurityGroupName)
	}

	subnetID := fmt.Sprintf("%s/subnets/%s", vnetID, occmDetails.SubnetID)

	(*params)["virtualNetworkId"] = map[string]string{
		"value": vnetID,
	}

	(*params)["networkSecurityGroupName"] = map[string]string{
		"value": networkSecurityGroupName,
	}

	(*params)["virtualMachineSize"] = map[string]string{
		"value": occmDetails.VirtualMachineSize,
	}

	(*params)["subnetId"] = map[string]string{
		"value": subnetID,
	}

	deploymentsClient := resources.NewDeploymentsClient(occmDetails.SubscriptionID)
	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		log.Print("Could not authorize azure ", err)
		return "", err
	}
	deploymentsClient.Authorizer = authorizer

	deploymentFuture, err := deploymentsClient.CreateOrUpdate(
		context.Background(),
		occmDetails.ResourceGroup,
		occmDetails.Name,
		resources.Deployment{
			Properties: &resources.DeploymentProperties{
				Template:   template,
				Parameters: params,
				Mode:       resources.Incremental,
			},
		},
	)
	if err != nil {
		return "", err
	}
	err = deploymentFuture.Future.WaitForCompletionRef(context.Background(), deploymentsClient.BaseClient.Client)
	if err != nil {
		return "", err
	}

	return "", nil

}

// CallGetAzureVM can be used to make a request to get Azure VM
func (c *Client) CallGetAzureVM(occmDetails createOCCMDetails) (string, error) {

	deploymentsClient := resources.NewDeploymentsClient(occmDetails.SubscriptionID)
	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		log.Print("Could not authorize azure ", err)
		return "", err
	}
	deploymentsClient.Authorizer = authorizer

	deploymentFuture, err := deploymentsClient.Get(
		context.Background(),
		occmDetails.ResourceGroup,
		occmDetails.Name,
	)
	if err != nil {
		return "", err
	}

	id := *deploymentFuture.ID
	s := strings.Split(id, "/")
	id = s[len(s)-1]

	return id, nil

}

// CallDeleteAzureVM can be used to make a request to delete Azure VM
func (c *Client) CallDeleteAzureVM(occmDetails deleteOCCMDetails) error {

	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		return fmt.Errorf("cannot authorize: %v", err)
	}

	log.Print("deleting vm")

	vmClient := compute.NewVirtualMachinesClient(occmDetails.SubscriptionID)
	vmClient.Authorizer = authorizer
	vmFuture, err := vmClient.Delete(
		context.Background(),
		occmDetails.ResourceGroup,
		occmDetails.Name,
	)
	if err != nil {
		return fmt.Errorf("cannot delete vm: %v", err)
	}

	err = vmFuture.WaitForCompletionRef(context.Background(), vmClient.Client)
	if err != nil {
		return fmt.Errorf("cannot get the vm delete future response: %v", err)
	}

	log.Print("deleting nic")

	nicClient := network.NewInterfacesClient(occmDetails.SubscriptionID)
	nicClient.Authorizer = authorizer
	nicFuture, err := nicClient.Delete(
		context.Background(),
		occmDetails.ResourceGroup,
		occmDetails.Name+"-nic",
	)
	if err != nil {
		return fmt.Errorf("cannot delete nic: %v", err)
	}

	err = nicFuture.WaitForCompletionRef(context.Background(), nicClient.Client)
	if err != nil {
		return fmt.Errorf("cannot get the nic delete future response: %v", err)
	}

	log.Print("deleting nsg")

	nsgClient := network.NewSecurityGroupsClient(occmDetails.SubscriptionID)
	nsgClient.Authorizer = authorizer
	nsgFuture, err := nsgClient.Delete(
		context.Background(),
		occmDetails.ResourceGroup,
		occmDetails.Name+"-nsg",
	)
	if err != nil {
		return fmt.Errorf("cannot delete nsg: %v", err)
	}
	err = nsgFuture.WaitForCompletionRef(context.Background(), nsgClient.Client)
	if err != nil {
		return fmt.Errorf("cannot get the nsg delete future response: %v", err)
	}

	log.Print("deleting storage account")

	storageAccountsClient := storage.NewAccountsClient(occmDetails.SubscriptionID)
	storageAccountsClient.Authorizer = authorizer
	_, err = storageAccountsClient.Delete(
		context.Background(),
		occmDetails.ResourceGroup,
		occmDetails.Name+"sa",
	)
	if err != nil {
		return fmt.Errorf("cannot delete storage account: %v", err)
	}

	log.Print("deleting ipaddress")

	ipClient := network.NewPublicIPAddressesClient(occmDetails.SubscriptionID)
	ipClient.Authorizer = authorizer
	ipFuture, err := ipClient.Delete(
		context.Background(),
		occmDetails.ResourceGroup,
		occmDetails.Name+"-ip",
	)
	if err != nil {
		return fmt.Errorf("cannot delete ipaddress: %v", err)
	}

	err = ipFuture.WaitForCompletionRef(context.Background(), ipClient.Client)
	if err != nil {
		return fmt.Errorf("cannot get the ipaddress delete future response: %v", err)
	}

	log.Print("deleting deployment")

	deploymentsClient := resources.NewDeploymentsClient(occmDetails.SubscriptionID)

	if err != nil {
		log.Print("Could not authorize azure ", err)
		return err
	}
	deploymentsClient.Authorizer = authorizer

	deploymentFuture, err := deploymentsClient.Delete(
		context.Background(),
		occmDetails.ResourceGroup,
		occmDetails.Name,
	)
	if err != nil {
		return err
	}
	err = deploymentFuture.Future.WaitForCompletionRef(context.Background(), deploymentsClient.BaseClient.Client)
	if err != nil {
		return err
	}

	return nil

}

// CallAMIGet can be used to make a request to get AWS AMI
func (c *Client) CallAMIGet(occmDetails createOCCMDetails) (string, error) {
	sess := session.Must(session.NewSession(aws.NewConfig().WithRegion(occmDetails.Region)))
	svc := ec2.New(sess)
	input := &ec2.DescribeImagesInput{
		Owners: []*string{
			aws.String(c.AWSAccount),
		},
		Filters: []*ec2.Filter{
			{
				Name: aws.String("name"),
				Values: []*string{
					aws.String(c.AMIFilter),
				},
			},
		},
	}

	result, err := svc.DescribeImages(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return "", aerr
		}
		return "", err
	}

	latestDate := *result.Images[0].CreationDate
	latestAMI := *result.Images[0].ImageId
	for _, image := range result.Images {
		if *image.CreationDate > latestDate {
			latestDate = *image.CreationDate
			latestAMI = *image.ImageId
		}
	}

	return latestAMI, nil
}

// CallVPCGet can be used to make a request to get AWS AMI
func (c *Client) CallVPCGet(subnet string, region string) (string, error) {

	sess := session.Must(session.NewSession(aws.NewConfig().WithRegion(region)))
	svc := ec2.New(sess)
	input := &ec2.DescribeSubnetsInput{
		SubnetIds: []*string{
			aws.String(subnet),
		},
	}

	result, err := svc.DescribeSubnets(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return "", aerr
		}
		return "", err
	}

	log.Print("CallVPCGet ", *result.Subnets[0].VpcId)

	return *result.Subnets[0].VpcId, nil
}

// CallVNetGet can be used to make a request to get Azure virtual network
func (c *Client) CallVNetGet(subscriptionID string, resourceGroup string) (string, string, error) {

	vNet := network.NewVirtualNetworksClient(subscriptionID)
	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		log.Print("Could not authorize azure ", err)
		return "", "", err
	}
	vNet.Authorizer = authorizer

	for resList, err := vNet.ListComplete(context.Background(), resourceGroup); resList.NotDone(); err = resList.Next() {
		if err != nil {
			log.Print("Could not get vNet ", err)
			return "", "", err
		}
		name := *resList.Value().Name
		cidr := *resList.Value().VirtualNetworkPropertiesFormat.AddressSpace.AddressPrefixes
		log.Print("CallVNetGet ", cidr)
		log.Print("CallVNetGet ", name)
		return name, cidr[0], nil
	}

	return "", "", fmt.Errorf("vNet not found")
}

// CallVNetGetCidr can be used to make a request to get Azure virtual network
func (c *Client) CallVNetGetCidr(subscriptionID string, resourceGroup string, vnet string) (string, error) {

	vNet := network.NewVirtualNetworksClient(subscriptionID)
	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		log.Print("Could not authorize azure ", err)
		return "", err
	}
	vNet.Authorizer = authorizer

	resList, err := vNet.Get(context.Background(), resourceGroup, vnet, "")
	if err != nil {
		log.Print("Could not get cidr ", err)
		return "", err
	}
	name := *resList.Name
	cidr := *resList.VirtualNetworkPropertiesFormat.AddressSpace.AddressPrefixes
	log.Print("CallVNetGetCidr ", cidr)
	log.Print("CallVNetGetCidr ", name)
	return cidr[0], nil
}

// CallAWSInstanceGet can be used to make a request to get AWS Instance
func (c *Client) CallAWSInstanceGet(occmDetails createOCCMDetails) ([]ec2.Instance, error) {
	if occmDetails.Region == "" {
		regions, err := c.CallAWSRegionGet(occmDetails)
		if err != nil {
			return nil, err
		}
		var res []ec2.Instance
		for _, region := range regions {
			regionReservation, err := c.CallAWSGetReservationsForRegion(region)
			if err != nil {
				return nil, err
			}
			res = append(res, regionReservation...)
		}
		return res, nil
	}
	sess := session.Must(session.NewSession(aws.NewConfig().WithRegion(occmDetails.Region)))
	svc := ec2.New(sess)
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-type"),
				Values: []*string{
					aws.String(occmDetails.InstanceType),
				},
			},
			{
				Name: aws.String("key-name"),
				Values: []*string{
					aws.String(occmDetails.KeyName),
				},
			},
			{
				Name: aws.String("subnet-id"),
				Values: []*string{
					aws.String(occmDetails.SubnetID),
				},
			},
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(occmDetails.Name),
				},
			},
		},
	}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return nil, aerr
			}
		}
		return nil, err
	}

	var res []ec2.Instance
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			res = append(res, *instance)
		}
	}

	return res, nil
}

// CallAWSRegionGet describe all regions.
func (c *Client) CallAWSRegionGet(occmDetails createOCCMDetails) ([]string, error) {
	sess := session.Must(session.NewSession())
	svc := ec2.New(sess)

	result, err := svc.DescribeRegions(nil)
	if err != nil {
		log.Printf("CallAWSRegionGet error: %#v", err)
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return nil, aerr
			}
		}
		return nil, err
	}

	var res []string
	for _, region := range result.Regions {
		res = append(res, *region.RegionName)
	}
	return res, nil
}

// CallAPIMethod can be used to make a request to any CVO/OCCM API method, receiving results as byte
func (c *Client) CallAPIMethod(method string, baseURL string, params map[string]interface{}, token string, hostType string) (int, []byte, string, error) {
	c.initOnce.Do(c.init)

	c.waitForAvailableSlot()
	defer c.releaseSlot()

	ourlog.WithFields(logrus.Fields{
		"method": method,
		"params": params,
	}).Debug("Calling API")

	paramsNil := false
	if params == nil {
		paramsNil = true
	}
	statusCode, result, onCloudRequestID, err := c.restapiClient.Do(baseURL, hostType, token, paramsNil, c.AccountID, c.ClientID, &restapi.Request{
		Method:                method,
		Params:                params,
		GCPDeploymentTemplate: c.GCPDeploymentTemplate,
		GCPServiceAccountPath: c.GCPServiceAccountPath,
	})
	if err != nil {
		return statusCode, nil, "", err
	}
	ourlog.WithFields(logrus.Fields{
		"method": method,
	}).Debug("Received successful API response")
	return statusCode, result, onCloudRequestID, nil
}

func (c *Client) init() {
	if c.MaxConcurrentRequests == 0 {
		c.MaxConcurrentRequests = 6
	}
	c.requestSlots = make(chan int, c.MaxConcurrentRequests)
	c.restapiClient = &restapi.Client{
		CloudManagerHost:     c.CloudManagerHost,
		AuthHost:             c.AuthHost,
		SaAuthHost:           c.SaAuthHost,
		CVOHostName:          c.CVOHostName,
		RefreshToken:         c.RefreshToken,
		SaSecretKey:          c.SaSecretKey,
		SaClientID:           c.SaClientID,
		Audience:             c.Audience,
		GCPDeploymentManager: c.GCPDeploymentManager,
		CVSHostName:          c.CVSHostName,
	}
}

// SetRefreshToken for the client to use for requests to the CVO/OCCM API
func (c *Client) SetRefreshToken(refreshToken string) {
	c.RefreshToken = refreshToken
}

// GetRefreshToken returns the API version that will be used for CVO/OCCM API requests
func (c *Client) GetRefreshToken() string {
	return c.RefreshToken
}

// SetServiceCredential for the client to use for requests to the CVO/OCCM API
func (c *Client) SetServiceCredential(SaSecretKey string, SaClientID string) {
	c.SaSecretKey = SaSecretKey
	c.SaClientID = SaClientID
}

// GetServiceCredential returns the service account secret key and secret client id that will be used for CVO/OCCM API requests
func (c *Client) GetServiceCredential() (string, string) {
	return c.SaSecretKey, c.SaClientID
}

func (c *Client) waitForAvailableSlot() {
	c.requestSlots <- 1
}

func (c *Client) releaseSlot() {
	<-c.requestSlots
}

// SetSimulator for the client to use for tests on simulator
func (c *Client) SetSimulator(simulator bool) {
	c.Simulator = simulator
}

// GetSimulator returns if it is set running on simulator
func (c *Client) GetSimulator() bool {
	return c.Simulator
}

// CallAWSTagCreate creates tag
func (c *Client) CallAWSTagCreate(occmDetails createOCCMDetails) error {
	sess := session.Must(session.NewSession(aws.NewConfig().WithRegion(occmDetails.Region)))

	svc := ec2.New(sess)

	tags := []*ec2.Tag{}
	if len(occmDetails.AwsTags) > 0 {
		for _, awsTag := range occmDetails.AwsTags {
			tag := &ec2.Tag{
				Key:   aws.String(awsTag.TagKey),
				Value: aws.String(awsTag.TagValue),
			}
			tags = append(tags, tag)
		}
	}

	input := &ec2.CreateTagsInput{
		Resources: []*string{
			aws.String(occmDetails.InstanceID),
		},
		Tags: tags,
	}

	result, err := svc.CreateTags(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return aerr
		}
		return err
	}

	fmt.Println(result)
	return nil
}

// CallAWSTagDelete deletes tag
func (c *Client) CallAWSTagDelete(occmDetails createOCCMDetails) error {
	sess := session.Must(session.NewSession(aws.NewConfig().WithRegion(occmDetails.Region)))

	svc := ec2.New(sess)

	tags := []*ec2.Tag{}
	if len(occmDetails.AwsTags) > 0 {
		for _, awsTag := range occmDetails.AwsTags {
			tag := &ec2.Tag{
				Key:   aws.String(awsTag.TagKey),
				Value: aws.String(awsTag.TagValue),
			}
			tags = append(tags, tag)
		}
	}

	input := &ec2.DeleteTagsInput{
		Resources: []*string{
			aws.String(occmDetails.InstanceID),
		},
		Tags: tags,
	}

	result, err := svc.DeleteTags(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return aerr
		}
		return err
	}

	fmt.Println(result)
	return nil
}

// CallAWSDescribeInstanceAttribute returns disableAPITermination.
func (c *Client) CallAWSDescribeInstanceAttribute(occmDetails createOCCMDetails) (bool, error) {
	sess := session.Must(session.NewSession(aws.NewConfig().WithRegion(occmDetails.Region)))
	svc := ec2.New(sess)
	input := &ec2.DescribeInstanceAttributeInput{
		Attribute:  aws.String("disableApiTermination"),
		InstanceId: aws.String(occmDetails.InstanceID),
	}

	result, err := svc.DescribeInstanceAttribute(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return false, aerr
			}
		}
		return false, err
	}
	disableAPITermination := *result.DisableApiTermination.Value
	if err != nil {
		return false, err
	}

	return disableAPITermination, nil
}

// CallAWSGetReservationsForRegion gets reservations for a region.
func (c *Client) CallAWSGetReservationsForRegion(region string) ([]ec2.Instance, error) {

	var res []ec2.Instance

	sess := session.Must(session.NewSession(aws.NewConfig().WithRegion(region)))
	svc := ec2.New(sess)
	input := &ec2.DescribeInstancesInput{}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return nil, aerr
			}
		}
		return nil, err
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			res = append(res, *instance)
		}
	}

	return res, err
}
