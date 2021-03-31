resource "netapp-cloudmanager_nss_account" "nss-account" {
		provider = netapp-cloudmanager
		client_id = "Rw5Q2O1kdnLtvhwegGalFnodEHUfPJWh"
		username = "user"
		password = "password"
	}

data "netapp-cloudmanager_nss_account" "nss-account-2" {
		provider = netapp-cloudmanager
		client_id = "Rw5Q2O1kdnLtvhwegGalFnodEHUfPJWh"
		username = "user"
	}
