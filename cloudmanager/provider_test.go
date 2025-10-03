package cloudmanager

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	if Provider() == nil {
		t.Fatal("Provider should not be nil")
	}
}

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"netapp-cloudmanager": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("CLOUDMANAGER_REFRESH_TOKEN"); v == "" {
		t.Fatal("CLOUDMANAGER_REFRESH_TOKEN must be set for acceptance tests")
	}
}
