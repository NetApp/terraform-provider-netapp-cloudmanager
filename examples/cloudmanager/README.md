# NetApp NetApp_CloudManager 0.1.0 Example

This repository is designed to demonstrate the capabilities of the [Terraform
NetApp netapp_cloudmanager Provider][ref-tf-netapp-cloudmanager] at the time of the 0.1.0 release.

[ref-tf-netapp-cloudmanager]: https://www.terraform.io/docs/providers/netapp/cloudmanager/index.html

This example performs the following:

* Creates a number of aws occm,
  using the [`netapp_cloudmanager_connector_aws` resource][ref-tf-netapp-cloudmanager-connector-aws].

[ref-tf-netapp-cloudmanager-connector-aws]: https://www.terraform.io/docs/providers/netapp/cloudmanager/r/occm_aws.html

## Requirements

* A working AWS, GCP and ANF account.

## Usage Details

You can either clone the entire
[terraform-provider-netapp_cloudmanager][ref-tf-netapp-cloudmanager-github] repository, or download the
`provider.tf`, `variables.tf`, `resources.tf`, and
`terraform.tfvars.example` files into a directory of your choice. Once done,
edit the `terraform.tfvars.example` file, populating the fields with the
relevant values, and then rename it to `terraform.tfvars`. Don't forget to
configure your endpoint and credentials by either adding them to the
`provider.tf` file, or by using enviornment variables. See
[here][ref-tf-netapp-cloudmanager-provider-settings] for a reference on provider-level
configuration values.

[ref-tf-netapp-cloudmanager-github]: https://github.com/terraform-providers/terraform-provider-netapp-cloudmanager
[ref-tf-netapp-cloudmanager-provider-settings]: https://www.terraform.io/docs/providers/netapp/cloudmanager/index.html#argument-reference

Once done, run `terraform init`, and `terraform plan` to review the plan, then
`terraform apply` to execute. If you use Terraform 0.11.0 or higher, you can
skip `terraform plan` as `terraform apply` will now perform the plan for you and
ask you confirm the changes.
