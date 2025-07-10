terraform {
  required_providers {
    netapp-cloudmanager = {
      source = "NetApp/netapp-cloudmanager"
      version = "20.10.0"
    }
  }
}

resource "netapp-cloudmanager_connector_aws" "cm-aws" {
  provider = netapp-cloudmanager
  name = "TerraformnAWS"
  region = "us-east-1"
  company = "NetApp"
  key_name = "key1"
  subnet_id = "subnet-xxxxxxxx"
  security_group_id = "sg-xxxxxxxx"
  iam_instance_profile_name = "OCCM"
  account_id = "account-xxxxxxx"
}

resource "netapp-cloudmanager_cvo_aws" "cvo-aws" {
  provider = netapp-cloudmanager
  name = "TerraformCVO1"
  region = "us-east-1"
  subnet_id = "subnet-xxxxxxxx"
  vpc_id = "vpc-xxx"
  svm_password = "********"
  client_id = netapp-cloudmanager_connector_aws.cm-aws.client_id 
}

resource "netapp-cloudmanager_cvo_aws" "cvo-aws-2" {
  provider = netapp-cloudmanager
  name = "TerraformCVO2"
  region = "us-east-1"
  subnet_id = "subnet-xxxxxxxx"
  vpc_id = "vpc-xxx"
  svm_password = "********"
  client_id = netapp-cloudmanager_connector_aws.cm-aws.client_id 
}

resource "netapp-cloudmanager_volume" "cvo-volume2" {
  provider = netapp-cloudmanager
  name = "test_vo2"
  size = 10
  unit = "GB"
  provider_volume_type = "gp3"
  iops = 3000
  throughput = 1000
  capacity_tier = "none"
  enable_thin_provisioning = true
  enable_compression = true
  enable_deduplication = true
  snapshot_policy_name = "default"
  export_policy_type = "custom"
  export_policy_ip = ["10.30.0.1/16"]
  export_policy_nfs_version = ["nfs4"]
  working_environment_name = netapp-cloudmanager_cvo_aws.cvo-aws-2.name
  client_id = netapp-cloudmanager_connector_aws.cm-aws.client_id 
}

resource "netapp-cloudmanager_snapmirror" "cl-snapmirror" {
  client_id = netapp-cloudmanager_connector_aws.cm-aws.client_id 
  destination_volume_name = "tgt_vol1"
  destination_working_environment_id = netapp-cloudmanager_cvo_aws.cvo-aws.id
  max_transfer_rate = 102400
  policy = "MirrorAllSnapshots"
  schedule = "5min"
  source_volume_name = netapp-cloudmanager_volume.cvo-volume2.name
  source_working_environment_id = netapp-cloudmanager_cvo_aws.cvo-aws-2.id
  delete_destination_volume = true  # Enable automatic deletion of destination volume when snapmirror is destroyed
}
