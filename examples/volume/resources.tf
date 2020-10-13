resource "netapp-cloudmanager_volume" "cvo-volume" {
  provider = netapp-cloudmanager
  name = "test_vol"
  size = 10
  unit = "GB"
  enable_thin_provisioning = true
  enable_compression = true
  enable_deduplication = true
  snapshot_policy_name = "default"
  export_policy_type = "custom"
  export_policy_ip = ["10.30.0.1/16"]
  export_policy_nfs_version = ["nfs4"]
  working_environment_name = "justincluster"
  client_id = "Nw4Q2O1kdnLtvhwegGalFnodEHUfPJWh"
}

resource "netapp-cloudmanager_volume" "cvo-volume-2" {
  provider = netapp-cloudmanager
  name = "test_vol_2"
  size = 2
  unit = "TB"
  enable_thin_provisioning = true
  enable_compression = true
  enable_deduplication = true
  snapshot_policy_name = "default"
  working_environment_name = "justincluster"
  export_policy_name = "export-svm_justincluster-test_vol_2"
  export_policy_type = "custom"
  export_policy_ip = ["10.30.1.1/16"]
  export_policy_nfs_version = ["nfs3", "nfs4"]
  capacity_tier = "S3"
  tiering_policy = "auto"
  client_id = "Nw4Q2O1kdnLtvhwegGalFnodEHUfPJWh"
}

resource "netapp-cloudmanager_volume" "cvo-volume-3" {
  provider = netapp-cloudmanager
  name = "test_vol_3"
  size = 10
  unit = "GB"
  enable_thin_provisioning = true
  enable_compression = true
  enable_deduplication = true
  snapshot_policy_name = "default"
  working_environment_name = "justincluster"
  export_policy_name = "export-svm_justincluster-test_vol_3"
  export_policy_type = "custom"
  export_policy_ip = ["10.30.1.1/16"]
  export_policy_nfs_version = ["nfs3"]
  capacity_tier = "S3"
  tiering_policy = "auto"
  client_id = "Nw4Q2O1kdnLtvhwegGalFnodEHUfPJWh"
}

resource "netapp-cloudmanager_volume" "cvo-volume-4" {
  provider = netapp-cloudmanager
  name = "test_vol_4"
  size = 10
  unit = "GB"
  enable_thin_provisioning = true
  enable_compression = true
  enable_deduplication = true
  snapshot_policy_name = "default"
  working_environment_name = "justincluster"
  provider_volume_type = "io1"
  iops = 100
  export_policy_name = "export-svm_justincluster-test_vol_4"
  export_policy_type = "custom"
  export_policy_ip = ["10.30.1.1/16"]
  export_policy_nfs_version = ["nfs3"]
  capacity_tier = "S3"
  tiering_policy = "auto"
  client_id = "Nw4Q2O1kdnLtvhwegGalFnodEHUfPJWh"
}

resource "netapp-cloudmanager_volume" "cvo-volume-ha-1" {
  provider = netapp-cloudmanager
  name = "ha_test_vol"
  size = 10
  unit = "GB"
  enable_thin_provisioning = true
  enable_compression = true
  enable_deduplication = true
  snapshot_policy_name = "default"
  export_policy_type = "custom"
  export_policy_ip = ["10.30.0.1/16"]
  export_policy_nfs_version = ["nfs4"]
  working_environment_name = "justinclusters"
  client_id = "Nw4Q2O1kdnLtvhwegGalFnodEHUfPJWh"
}

resource "netapp-cloudmanager_volume" "cvo-volume-ha-2" {
  provider = netapp-cloudmanager
  name = "ha_test_vol_2"
  size = 2
  unit = "TB"
  enable_thin_provisioning = true
  enable_compression = true
  enable_deduplication = true
  snapshot_policy_name = "default"
  working_environment_name = "justinclusters"
  export_policy_type = "custom"
  export_policy_ip = ["10.30.1.1/16"]
  export_policy_nfs_version = ["nfs3", "nfs4"]
  capacity_tier = "S3"
  tiering_policy = "auto"
  client_id = "Nw4Q2O1kdnLtvhwegGalFnodEHUfPJWh"
}

resource "netapp-cloudmanager_volume" "cifs-volume-1" {
  provider = netapp-cloudmanager
  name = "cifs_test_vol_1"
  volume_protocol = "cifs"
  size = 10
  unit = "GB"
  enable_thin_provisioning = true
  enable_compression = true
  enable_deduplication = true
  snapshot_policy_name = "default"
  share_name = "share_cifs"
  permission = "full_control"
  users = ["Everyone"]
  working_environment_name = "justincluster"
  client_id = "Nw4Q2O1kdnLtvhwegGalFnodEHUfPJWh"
  capacity_tier= "S3"
  tiering_policy = "auto"
}

	resource "netapp-cloudmanager_volume" "iscsi-volume-1" {
		provider = netapp-cloudmanager
		name = "iscsi_test_vol"
    volume_protocol = "iscsi"
		size = 10
		unit = "GB"
		enable_thin_provisioning = true
		enable_compression = false
		enable_deduplication = false
		snapshot_policy_name = "default"
		working_environment_name = "justincluster"
    igroups = ["test_igroup"]
    initiator {
      alias = "test_alias"
      iqn = "test_iqn"
    }
    os_name = "linux"
		client_id = "Nw4Q2O1kdnLtvhwegGalFnodEHUfPJWh"
  }