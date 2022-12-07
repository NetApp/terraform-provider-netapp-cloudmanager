# Specify CVO resources

terraform {
  required_providers {
    netapp-cloudmanager = {
      source = "NetApp/netapp-cloudmanager"
      version = "20.10.0"
    }
  }
}

resource "netapp-cloudmanager_connector_gcp" "cm-gcp" {
  provider = netapp-cloudmanager
  name = "occm-gcp"
  project_id = "project-id"
  zone = "us-east1-b"
  company = "NetApp"
  service_account_email = "cloudmanager-service-account@project-id.iam.gserviceaccount.com"
  service_account_path = "/Users/name/Downloads/project-id-terraform-user.json"
  account_id = "account-xxxxxxx"
}

resource "netapp-cloudmanager_cvo_gcp" "cvo-gcp" {
  provider = netapp-cloudmanager
  name = "terraformcvogcp"
  project_id = "project-id"
  zone = "us-east1-b"
  subnet_id = "default"
  gcp_service_account = "fabric-pool@project-id.iam.gserviceaccount.com"
  svm_password = "********"
  client_id = netapp-cloudmanager_connector_gcp.cm-gcp.client_id 
}

resource "netapp-cloudmanager_aggregate" "cvo-gcp-aggregate" {
  provider = netapp-cloudmanager
  name = "aggr2"
  working_environment_id = netapp-cloudmanager_cvo_gcp.cvo-gcp.id 
  client_id = netapp-cloudmanager_connector_gcp.cm-gcp.client_id 
  number_of_disks = 1
  provider_volume_type = "pd-standard"
  capacity_tier = "NONE"
}

resource "netapp-cloudmanager_cifs_server" "cvo-cifs-workgroup" {
   depends_on = [netapp-cloudmanager_aggregate.cvo-gcp-aggregate]
   provider = netapp-cloudmanager
   server_name = "server"
   workgroup_name  = "workgroup"
   client_id = netapp-cloudmanager_connector_gcp.cm-gcp.client_id 
   working_environment_name = "terraformcvogcp"
   is_workgroup = true
}

resource "netapp-cloudmanager_volume" "cvo-volume" {
  depends_on = [netapp-cloudmanager_volume.cifs-volume-1,netapp-cloudmanager_cifs_server.cvo-cifs-workgroup]
  provider = netapp-cloudmanager
  name = "vol1"
  size = 10
  unit = "GB"
  provider_volume_type = "pd-standard"
  export_policy_type = "custom"
  export_policy_ip = ["0.0.0.0/0"]
  export_policy_nfs_version = ["nfs4"]
  working_environment_id = netapp-cloudmanager_cvo_gcp.cvo-gcp.id 
  client_id = netapp-cloudmanager_connector_gcp.cm-gcp.client_id 
  capacity_tier = "cloudStorage"
  tiering_policy = "auto"
}

resource "netapp-cloudmanager_volume" "cifs-volume-1" {
  depends_on = [netapp-cloudmanager_cifs_server.cvo-cifs-workgroup]
  provider = netapp-cloudmanager
  name = "cifs_vol2"
  volume_protocol = "cifs"
  provider_volume_type = "pd-ssd"
  size = 10
  unit = "GB"
  share_name = "share_cifs"
  permission = "full_control"
  users = ["Everyone"]
  working_environment_id = netapp-cloudmanager_cvo_gcp.cvo-gcp.id
  client_id = netapp-cloudmanager_connector_gcp.cm-gcp.client_id
}

resource "netapp-cloudmanager_cvo_gcp" "cvo-gcp_ha" {
  provider = netapp-cloudmanager
  name = "terraformcvogcpHA"
  project_id = "default-project"
  zone = "us-east1-b"
  subnet_id = "default"
  vpc_id = "default"
  gcp_service_account = "gcpservice@project-id.iam.gserviceaccount.com"
  gcp_label {
        label_key = "abcdefg"
        label_value = "0222"
  }
  is_ha = true
  svm_password = "********"
  svm {
        svm_name = "svmx"
  }
  svm {
        svm_name = "svmy"
  }
  use_latest_version = true
  ontap_version = "latest"
  gcp_volume_type = "pd-ssd"
  instance_type = "n2-standard-4"
  mediator_zone = "us-east1-d"
  node1_zone = "us-east1-b"
  node2_zone =  "us-east1-c"
  subnet0_node_and_data_connectivity = "projects/default-project/regions/us-east1/subnetworks/default"
  subnet1_cluster_connectivity = "projects/default-project/regions/us-east1/subnetworks/subnet2"
  subnet2_ha_connectivity = "projects/default-project/regions/us-east1/subnetworks/subnet3"
  subnet3_data_replication = "projects/default-project/regions/us-east1/subnetworks/occm-us-east1-subnet-1"
  vpc0_node_and_data_connectivity = "projects/default-project/global/networks/default"
  vpc1_cluster_connectivity = "projects/default-project/global/networks/vpc2"
  vpc2_ha_connectivity = "projects/default-project/global/networks/vpc3"
  vpc3_data_replication = "projects/default-project/global/networks/occm-network-1"
  nss_account = "c3c3d4e5-123d-4012-83b3-1e4abcdefg"
  client_id = netapp-cloudmanager_connector_gcp.cm-gcp.client_id
}