# Specify CVO resources


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
  svm_password = "netapp1!"
  client_id = netapp-cloudmanager_connector_gcp.cm-gcp.client_id 
  workspace_id = "workspace-xxxxxx"
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
