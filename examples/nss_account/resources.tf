resource "netapp-cloudmanager_nss_account" "nss-account" {
		provider = netapp-cloudmanager
		client_id = "Nw4Q2O1kdnLtvhwegGalFnodEHUfPJWh"
		name = "accName"
		username = "user"
		password = "password"
	}

data "netapp-cloudmanager_nss_account" "nss-account-2" {
		provider = netapp-cloudmanager
		client_id = "Nw4Q2O1kdnLtvhwegGalFnodEHUfPJWh"
		name = "accName"
	}