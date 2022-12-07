resource "netapp-cloudmanager_cifs_server" "cl-cifs" {
   provider = netapp-cloudmanager
   domain = "test.com"
   username = "admin"
   password = "abcde"
   dns_domain = "test.com"
   ip_addresses = ["1.0.0.2"]
   netbios = "justincluster"
   organizational_unit = "CN=Computers"
   client_id = "Nw4Q2O1kdnLtvhwegGalFnodEHUfPJWh"
   working_environment_name = "justincluster"
}