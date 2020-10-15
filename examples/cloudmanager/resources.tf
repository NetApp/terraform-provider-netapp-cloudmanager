# Specify cloudmanager resources

resource "netapp-cloudmanager_connector_aws" "cl-occm-aws" {
  provider = netapp-cloudmanager
  name = "TerraformnoIp"
  region = "us-west-1"
  ami = "ami-02eee4e2dcf3ebf56"
  company = "NetApp"
  key_name = "dev_automation"
  instance_type = "t3.xlarge"
  subnet_id = "subnet-9d7590f8"
  security_group_id = "sg-02796afe76d97ef21"
  iam_instance_profile_name = "OCCM_AUTOMATION"
  proxy_user_name = "test"
  proxy_password = "test"
  associate_public_ip_address = true
  account_id = "account-moKEW1b5"
}

resource "netapp-cloudmanager_connector_azure" "cl-occm-azure" {
  provider = netapp-cloudmanager
  name = "testTFAzure"
  location = "westus"
  subscription_id = "d333af45-0d07-4154-943d-c25fbbce1b18"
  company = "NetApp"
  resource_group = "occm_group_westus"
  subnet_id = "Subnet1"
  vnet_id = "Vnet1"
  proxy_user_name = "test"
  proxy_password = "test"
  associate_public_ip_address = true
  account_id = "account-moKEW1b5"
  admin_password = "********"
  admin_username = "********"
}

resource "netapp-cloudmanager_connector_gcp" "cl-occm-gcp" {
  provider = netapp-cloudmanager
  name = "occm-gcp"
  project_id = "tlv-support"
  zone = "us-east4-b"
  company = "NetApp"
  service_account_email = "terraform-user@tlv-support.iam.gserviceaccount.com"
  service_account_path = "gcp_creds.json"
  proxy_user_name = "test"
  proxy_password = "test"
  account_id = "account-moKEW1b5"
}

resource "netapp-cloudmanager_cvo_aws" "cl-cvo-aws" {
  provider = netapp-cloudmanager
  name = "TerraformCVO"
  region = "us-east-1"
  subnet_id = "subnet-1"
  vpc_id = "vpc-1"
  data_encryption_type = "AWS"
  aws_tag {
              tag_key = "abcd"
              tag_value = "ABCD"
            }
  aws_tag {
              tag_key = "xxx"
              tag_value = "YYY"
            }
  ebs_volume_type = "gp2"
  svm_password = "netapp1!"
  ontap_version = "latest"
  use_latest_version = true
  license_type = "ha-cot-standard-paygo"
  instance_type = "m5.2xlarge"
  client_id = "Nw4Q2O1kdnLtvhwegGalFnodEHUfPJWh"
  workspace_id = "workspace-abaaFgcQ"
  capacity_tier = "S3"
  nss_account = "CloudProviderAccount-thlalnlg"
  writing_speed_state = "NORMAL"
  instance_tenancy = "default"
  cloud_provider_account =  "InstanceProfile"
  backup_volumes_to_cbs = false
  enable_compliance = false
  enable_monitoring = false
  is_ha = true
  failover_mode = "FloatingIP"
  node1_subnet_id = "subnet-1"
  node2_subnet_id = "subnet-1"
  mediator_subnet_id = "subnet-1"
  mediator_key_pair_name = "key1"
  cluster_floating_ip = "2.1.1.1"
  data_floating_ip = "2.1.1.2"
  data_floating_ip2 = "2.1.1.3"
  svm_floating_ip = "2.1.1.4"
  route_table_ids = ["rt-1","rt-2"]
}

resource "netapp-cloudmanager_cvo_azure" "cl-cvo-azure" {
  provider = netapp-cloudmanager
  name = "TerraformCVOAzure"
  location = "westus"
  subscription_id = "1"
  subnet_id = "subnet1"
  vnet_id = "Vnetid2"
  vnet_resource_group = "rg"
  cidr = "10.1.0.0/23"
  data_encryption_type = "AZURE"
  azure_tag {
              tag_key = "abcd"
              tag_value = "ABCD"
            }
  azure_tag {
              tag_key = "xxx"
              tag_value = "YYY"
            }
  storage_type = "Premium_LRS"
  svm_password = "Netapp1!"
  ontap_version = "latest"
  use_latest_version = true
  license_type = "azure-cot-standard-paygo"
  instance_type = "Standard_DS4_v2"
  client_id = "Nw4Q2O1kdnLtvhwegGalFnodEHUfPJWh"
  workspace_id = "workspace-abaaFgcQ"
  capacity_tier = "Blob"
  nss_account = "CloudProviderAccount-thlalnlg"
  writing_speed_state = "NORMAL"
  cloud_provider_account =  "InstanceProfile"
  backup_volumes_to_cbs = false
  enable_compliance = false
  enable_monitoring = false
  is_ha = false
}

resource "netapp-cloudmanager_cvo_gcp" "cvo-gcp" {
  provider = netapp-cloudmanager
  name = "terraformcvogcp"
  project_id = "tlv-support"
  zone = "us-east4-b"
  subnet_id = "default"
  gcp_volume_type = "pd-ssd"
  gcp_service_account = "occmservice@occm-dev.iam.gserviceaccount.com"
  data_encryption_type = "GCP"
  svm_password = "netapp1!"
  ontap_version = "latest"
  use_latest_version = true
  license_type = "gcp-cot-standard-paygo"
  instance_type = "n1-standard-8"
  client_id = "Nw4Q2O1kdnLtvhwegGalFnodEHUfPJWh"
  workspace_id = "workspace-abaaFgcQ"
  capacity_tier = "cloudStorage"
  writing_speed_state = "NORMAL"
}

# create aggregate by working environment name
resource "netapp-cloudmanager_aggregate" "cl-aggregate-name" {
  provider = netapp-cloudmanager
  name = "aggrtestbyname"
  working_environment_name = "testAWS"
  client_id = "Nw4Q2O1kdnLtvhwegGalFnodEHUfPJWh"
  number_of_disks = 2
  disk_size_size = 100
  disk_size_unit = "GB"
}

# create aggregate by working environment id
resource "netapp-cloudmanager_aggregate" "cl-aggregate" {
  provider = netapp-cloudmanager
  name = "aggrtest"
  working_environment_id = "VsaWorkingEnvironment-GNqPuWST"
  client_id = "Nw4Q2O1kdnLtvhwegGalFnodEHUfPJWh"
  number_of_disks = 1
  disk_size_size = 500
  disk_size_unit = "GB"
}