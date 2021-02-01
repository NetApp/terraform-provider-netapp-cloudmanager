---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_cvs_gcp_volume"
sidebar_current: "docs-netapp-cloudmanager-resource-cvs-gcp-volume"
description: |-
  Provides a netapp-cloudmanager_cvs_gcp_volume resource. This can be used to create, and delete volumes for Cloud Volume Service on GCP.
---

# netapp-cloudmanager_cvs_gcp_volume

Provides a netapp-cloudmanager_cvs_gcp_volume resource. This can be used to create, and delete volumes for Cloud Volume Service on GCP.

## Example Usages

**Create netapp-cloudmanager_cvs_gcp_volume:**

```
resource "netapp-cloudmanager_cvs_gcp_volume" "test-1" {
  provider = netapp-cloudmanager
  name = "test_vol"
  size = 105
  size_unit = "gb"
  volume_path = "test_vol"
  protocol_types = ["NFSv3"]
  region = "us-east4"
  service_level = "low"
  account = "Demo_SIM"
  client_id = "clientid"
  network = "mynetwork"
  working_environment_name = "GCP_environment"
  export_policy {
    rule {
      allowed_clients = "1.0.0.0/0"
      rule_index = 1 
      unix_read_only= true
      unix_read_write = false
      nfsv3 = true
      nfsv4 = true
    }
    rule {
      allowed_clients= "10.0.0.0"
      rule_index = 2
      unix_read_only= true
      unix_read_write = false
      nfsv3 = true
      nfsv4 = true
    }
  }
  snapshot_policy {
    enabled = true
    hourly_schedule {
      snapshots_to_keep = 48
      minute = 1
    }
    daily_schedule {
      snapshots_to_keep = 14
      hour = 23
      minute = 2
    }
    weekly_schedule {
      snapshots_to_keep = 4
      hour = 1
      minute = 3
      day = "Monday"
    }
    monthly_schedule {
      snapshots_to_keep = 6
      hour = 2
      minute = 4
      days_of_month = 6
    }    
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the volume.
* `size` - (Required) The volume size, supported with decimal numbers.
* `size_unit` - (Required) [ 'gb' ].
* `volume_path` - (Required) The volume path.
* `protocol_types` (Required) [ 'nfsv3', 'nfsv4', 'cifs' ].
* `region` - (Required) The region where the volume is created.
* `service_level` - (Required) ['low' or 'medium' or 'high'].
* `network`  - (Required) The network VPC of the volume.
* `account` - (Required) The name of the account.
* `client_id` - (Required) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://cloudmanager.netapp.com](https://cloudmanager.netapp.com).
* `working_environment_name` - (Required) The working environment name.
* `export_policy` - (Optional) The rules of the export policy.
* `snapshot_policy` - (Optional) The set of Snapshot Policy attributes for volume.

The `snapshot_policy` block supports:
* `enabled` - (Optional) If enabled, make snapshots automatically according to the schedules. Default is false.
* `daily_schedule` - (Optional) If enabled, make a snapshot every day. Defaults to midnight.
* `hourly_schedule` - (Optional) If enabled, make a snapshot every hour e.g. at 04:00, 05:00, 06:00.
* `monthly_schedule` - (Optional) If enabled, make a snapshot every month at a specific day or days, defaults to the first day of the month at midnight
* `weekly_schedule` - (Optional) If enabled, make a snapshot every week at a specific day or days, defaults to Sunday at midnight.

The `daily_schedule` block supports:
* `hour` - (Optional) Set the hour to start the snapshot (0-23), defaults to midnight (0).
* `minute` - (Optional) Set the minute of the hour to start the snapshot (0-59), defaults to the top of the hour (0).
* `snapshots_to_keep` - (Optional) The maximum number of Snapshots to keep for the daily schedule.

The `hourly_schedule` block supports:
* `minute` - (Optional) Set the minute of the hour to start the snapshot (0-59), defaults to the top of the hour (0).

The `monthly_schedule` block supports:
* `days_of_month` - (Optional) Set the day or days of the month to make a snapshot (1-31). Accepts a comma delimited string of the day of the month e.g. '1,15,31'. Defaults to '1'.
* `hour` - (Optional) Set the hour to start the snapshot (0-23), defaults to midnight (0).
* `minute` - (Optional) Set the minute of the hour to start the snapshot (0-59), defaults to the top of the hour (0).
* `snapshots_to_keep` - (Optional) The maximum number of Snapshots to keep for the daily schedule.

The `weekly_schedule` block supports:
* `day` - Set the day or days of the week to make a snapshot. Accepts a comma delimited string of week day names in english. Defaults to 'Sunday'.
* `hour` - (Optional) Set the hour to start the snapshot (0-23), defaults to midnight (0).
* `minute` - (Optional) Set the minute of the hour to start the snapshot (0-59), defaults to the top of the hour (0).
* `snapshots_to_keep` - (Optional) The maximum number of Snapshots to keep for the daily schedule.

The `export_policy` block supports:
* `rule` - (Optional) Export Policy rule.

The `rule` block supports:
* `access` - (Optional) Defines the access type for clients matching the 'allowedClients' specification.
* `allowed_clients` - (Optional) Defines the client ingress specification (allowed clients) as a comma seperated string with IPv4 CIDRs, IPv4 host addresses and host names.
* `nfsv3` - (Optional) If enabled (true) the rule allows NFSv3 protocol for clients matching the 'allowedClients' specification.
* `nfsv4` - (Optional) If enabled (true) the rule allows NFSv4 protocol for clients matching the 'allowedClients' specification.


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The name of the volume.

