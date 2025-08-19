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

Arguments marked with “Forces new resource” will cause the resource to be recreated if their value is changed after creation.

The following arguments are supported:

* `name` - (Required, Forces new resource) The name of the volume.
* `size` - (Required, Forces new resource) The volume size, supported with decimal numbers.
* `size_unit` - (Required, Forces new resource) [ 'gb' ].
* `volume_path` - (Required, Forces new resource) The volume path.
* `protocol_types` (Required, Forces new resource) [ 'nfsv3', 'nfsv4', 'cifs' ].
* `region` - (Required, Forces new resource) The region where the volume is created.
* `service_level` - (Optional, Forces new resource) ['low' or 'medium' or 'high'].
* `network`  - (Required, Forces new resource) The network VPC of the volume.
* `account` - (Required, Forces new resource) The name of the account.
* `client_id` - (Required, Forces new resource) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `working_environment_name` - (Required, Forces new resource) The working environment name.
* `export_policy` - (Optional, Forces new resource) The rules of the export policy.
* `snapshot_policy` - (Optional, Forces new resource) The set of Snapshot Policy attributes for volume.

The `snapshot_policy` block supports:
* `enabled` - (Optional, Forces new resource) If enabled, make snapshots automatically according to the schedules. Default is false.
* `daily_schedule` - (Optional, Forces new resource) If enabled, make a snapshot every day. Defaults to midnight.
* `hourly_schedule` - (Optional, Forces new resource) If enabled, make a snapshot every hour e.g. at 04:00, 05:00, 06:00.
* `monthly_schedule` - (Optional, Forces new resource) If enabled, make a snapshot every month at a specific day or days, defaults to the first day of the month at midnight
* `weekly_schedule` - (Optional, Forces new resource) If enabled, make a snapshot every week at a specific day or days, defaults to Sunday at midnight.

The `daily_schedule` block supports:
* `hour` - (Optional, Forces new resource) Set the hour to start the snapshot (0-23), defaults to midnight (0).
* `minute` - (Optional, Forces new resource) Set the minute of the hour to start the snapshot (0-59), defaults to the top of the hour (0).
* `snapshots_to_keep` - (Optional, Forces new resource) The maximum number of Snapshots to keep for the daily schedule.

The `hourly_schedule` block supports:
* `minute` - (Optional, Forces new resource) Set the minute of the hour to start the snapshot (0-59), defaults to the top of the hour (0).
* `snapshots_to_keep` - (Optional, Forces new resource) The maximum number of Snapshots to keep for the hourly schedule.

The `monthly_schedule` block supports:
* `days_of_month` - (Optional, Forces new resource) Set the day or days of the month to make a snapshot (1-31). Accepts a comma delimited string of the day of the month e.g. '1,15,31'. Defaults to '1'.
* `hour` - (Optional, Forces new resource) Set the hour to start the snapshot (0-23), defaults to midnight (0).
* `minute` - (Optional, Forces new resource) Set the minute of the hour to start the snapshot (0-59), defaults to the top of the hour (0).
* `snapshots_to_keep` - (Optional, Forces new resource) The maximum number of Snapshots to keep for the monthly schedule.

The `weekly_schedule` block supports:
* `day` - (Optional, Forces new resource) Set the day or days of the week to make a snapshot. Accepts a comma delimited string of week day names in english. Defaults to 'Sunday'.
* `hour` - (Optional, Forces new resource) Set the hour to start the snapshot (0-23), defaults to midnight (0).
* `minute` - (Optional, Forces new resource) Set the minute of the hour to start the snapshot (0-59), defaults to the top of the hour (0).
* `snapshots_to_keep` - (Optional, Forces new resource) The maximum number of Snapshots to keep for the weekly schedule.

The `export_policy` block supports:
* `rule` - (Optional) Export Policy rule.

The `rule` block supports:
* `rule_index` -  (Optional, Forces new resource)
* `unix_read_only` - (Optional, Forces new resource)
* `unix_read_write`- (Optional, Forces new resource)
* `allowed_clients` - (Optional, Forces new resource) Defines the client ingress specification (allowed clients) as a comma seperated string with IPv4 CIDRs, IPv4 host addresses and host names.
* `nfsv3` - (Optional, Forces new resource) If enabled (true) the rule allows NFSv3 protocol for clients matching the 'allowedClients' specification.
* `nfsv4` - (Optional, Forces new resource) If enabled (true) the rule allows NFSv4 protocol for clients matching the 'allowedClients' specification.


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The name of the volume.

