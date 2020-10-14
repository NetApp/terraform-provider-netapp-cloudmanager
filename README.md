# Terraform Provider for NetApp Cloud Volumes ONTAP for AWS, GCP and ANF

This is the repository for the Terraform Provider for NetApp Cloud Volumes ONTAP (CVO) for AWS, GCP and ANF.  The Provider can be used
with Terraform to work with Cloud Volumes ONTAP for AWS, GCP and ANF resources.

For general information about Terraform, visit the [official
website][tf-website] and the [GitHub project page][tf-github].

[tf-website]: https://terraform.io/
[tf-github]: https://github.com/hashicorp/terraform

The provider plugin was developed by NetApp.

# Naming Conventions

The APIs for NetApp Cloud Volumes ONTAP for AWS, GCP and ANF do not require resource names to be unique.  They are considered
as 'labels' and resources are uniquely identified by 'ids'.  However these ids are not
user friendly, and as they are generated on the fly, they make it difficult to track
resources and automate.

This provider assumes that resource names are unique, and enforces it within its scope.
This is not an issue if everything is managed through Terraform, but could raise
conflicts if the rule is not respected outside of Terraform.

# Using the Provider

The current version of this provider requires Terraform 0.13 or higher to
run.

Terraform 0.13 introduces a registry, and you can use directly the provider without
building it yourself.
See https://registry.terraform.io/providers/NetApp/netapp-cloumanager

If you want to build it, see [the section below](#building-the-provider).

Note that you need to run `terraform init` to fetch the provider before
deploying.

## Provider Documentation

The documentation is available at:
https://registry.terraform.io/providers/NetApp/netapp-cloudmanager/latest/docs

The provider is also documented [here][tf-netapp-cloudmanager-docs].

Check the provider documentation for details on
entering your connection information and how to get started with writing
configuration for NetApp CVO resources.

[tf-netapp-cloudmanager-docs](website/docs/index.html.markdown)

### Controlling the provider version

Note that you can also control the provider version. This requires the use of a
`provider` block in your Terraform configuration if you have not added one
already.

The syntax is as follows:

```hcl
terraform {
  required_providers {
    netapp-gcp = {
      source = "NetApp/netapp-gcp"
      version = "20.10.0"
    }
  }
}
```

[Read more][provider-vc] on provider version control.

[provider-vc]: https://www.terraform.io/docs/configuration/provider-requirements.html#requiring-providers

# Building The Provider

## Prerequisites

If you wish to work on the provider, you'll first need [Go][go-website]
installed on your machine (version 1.11+ is **required**). You'll also need to
correctly setup a [GOPATH][gopath], as well as adding `$GOPATH/bin` to your
`$PATH`.

[go-website]: https://golang.org/
[gopath]: http://golang.org/doc/code.html#GOPATH

The following go packages are required to build the provider:
```
	github.com/Azure/azure-sdk-for-go v46.4.0+incompatible
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.3
	github.com/aws/aws-sdk-go v1.35.5
	github.com/fatih/structs v1.1.0
	github.com/hashicorp/terraform v0.13.4
	github.com/sirupsen/logrus v1.7.0
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	golang.org/x/tools v0.0.0-20201008025239-9df69603baec // indirect
```

Check go.mod for the latest list.

## Cloning the Project

First, you will want to clone the repository to
`$GOPATH/terraform-provider-netapp-cloudmanager`:

```sh
mkdir -p $GOPATH
cd $GOPATH
git clone https://github.com/NetApp/terraform-provider-netapp-cloudmanager.git
```

## Running the Build

After the clone has been completed, you can enter the provider directory and
build the provider.

```sh
cd $GOPATH/terraform-provider-netapp-cloudmanager
make build
```

## Installing the Local Plugin

With Terraform 0.13 or newer, see the [sanity check](#sanity-check) section under **Walkthrough example**.

With earlier versions of Terraform, after
the build is complete, copy the `terraform-provider-netapp-cloudmanager` binary into
the same path as your `terraform` binary, and re-run `terraform init`.

After this, your project-local `.terraform/plugins/ARCH/lock.json` (where `ARCH`
matches the architecture of your machine) file should contain a SHA256 sum that
matches the local plugin. Run `shasum -a 256` on the binary to verify the values
match.

# Developing the Provider

**NOTE:** Before you start work on a feature, please make sure to check the
[issue tracker][gh-issues] and existing [pull requests][gh-prs] to ensure that
work is not being duplicated. For further clarification, you can also ask in a
new issue.

[gh-issues]: https://github.com/netapp/terraform-provider-netapp-cloudmanager/issues
[gh-prs]: https://github.com/netapp/terraform-provider-netapp-cloudmanager/pulls

See [Building the Provider](#building-the-provider) for details on building the provider.

# Testing the Provider

**NOTE:** Testing the provider for NetApp Cloud Volumes ONTAP for AWS, GCP and ANF is currently a complex operation as it
requires having a NetApp CVO subscription in CVO to test against.
You can then use a .json file to expose your credentials.

## Configuring Environment Variables

Most of the tests in this provider require a comprehensive list of environment
variables to run. See the individual `*_test.go` files in the
[`cloudmanager/`](netapp_cloudmanager/) directory for more details. The next section also
describes how you can manage a configuration file of the test environment
variables.

### Using the `.tf-netapp-cloudmanager-devrc.mk` file

The [`tf-netapp-cloudmanager-devrc.mk.example`](tf-netapp-cloudmanager-devrc.mk.example) file contains
an up-to-date list of environment variables required to run the acceptance
tests. Copy this to `$HOME/.tf-netapp-cloudmanager-devrc.mk` and change the permissions to
something more secure (ie: `chmod 600 $HOME/.tf-netapp-cloudmanager-devrc.mk`), and
configure the variables accordingly.

## Running the Acceptance Tests

After this is done, you can run the acceptance tests by running:

```sh
$ make testacc
```

If you want to run against a specific set of tests, run `make testacc` with the
`TESTARGS` parameter containing the run mask as per below:

```sh
make testacc TESTARGS="-run=TestAccNetAppCVOOCCM"
```

This following example would run all of the acceptance tests matching
`TestAccNetAppCVOOCCM`. Change this for the specific tests you want to
run.


# Walkthrough example

### Installing go and terraform

```
bash
mkdir tf_na_netapp_cloudmanager
cd tf_na_netapp_cloudmanager

# if you want a private installation, use
export GO_INSTALL_DIR=`pwd`/go_install
mkdir $GO_INSTALL_DIR
# otherwise, go recommends to use
export GO_INSTALL_DIR=/usr/local
```

#### linux

```
curl -O https://dl.google.com/go/go1.15.2.linux-amd64.tar.gz
tar -C $GO_INSTALL_DIR -xvf go1.15.2.linux-amd64.tar.gz

export PATH=$PATH:$GO_INSTALL_DIR/go/bin

curl -O https://releases.hashicorp.com/terraform/0.13.4/terraform_0.13.4_linux_amd64.zip
unzip terraform_0.13.4_linux_amd64.zip
mv terraform $GO_INSTALL_DIR/go/bin
```

#### mac

```
curl -O https://dl.google.com/go/go1.15.2.darwin-amd64.tar.gz
tar -C $GO_INSTALL_DIR -xvf go1.15.2.darwin-amd64.tar.gz

export PATH=$PATH:$GO_INSTALL_DIR/go/bin

curl -O https://releases.hashicorp.com/terraform/0.13.4/terraform_0.13.4_darwin_amd64.zip
unzip terraform_0.13.4_darwin_amd64.zip
mv terraform $GO_INSTALL_DIR/go/bin
```

### Installing dependencies

We're using go.mod to manage dependencies, so there is not much to do.
```
# make sure git is installed
which git

export GOPATH=`pwd`
```

### Cloning the NetApp provider repository and building the provider


```
git clone https://github.com/NetApp/terraform-provider-netapp-cloudmanager.git
cd terraform-provider-netapp-cloudmanager
make build
# binary is in: $GOPATH/bin/terraform-provider-netapp-cloudmanager
```

The build step will install the provider in the $GOPATH/bin directory.

### Sanity check

#### Local installation - linux

```
mkdir -p /tmp/terraform/netapp.com/netapp/netapp-cloudmanager/20.10.0/linux_amd64
cp $GOPATH/bin/terraform-provider-netapp-cloudmanager /tmp/terraform/netapp.com/netapp/netapp-cloudmanager/20.10.0/linux_amd64
```

#### Local installation - mac

```
mkdir -p /tmp/terraform/netapp.com/netapp/netapp-cloudmanager/20.10.0/darwin_amd64
cp $GOPATH/bin/terraform-provider-netapp-cloudmanager /tmp/terraform/netapp.com/netapp/netapp-cloudmanager/20.10.0/darwin_amd64
```

#### Check the provider can be loaded
```
cd examples/cloudmanager/local
export TF_CLI_CONFIG_FILE=`pwd`/terraform.rc
terraform init
```

Should do nothing but indicate that `Terraform has been successfully initialized!`
