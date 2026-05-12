data "checkpointsase_standard_networks" "all" {}

# Publishes an internal HTTPS application to the first available network.
resource "checkpointsase_application" "internal_app" {
  name    = "myInternalApp"
  type    = "https"
  network = data.checkpointsase_standard_networks.all.networks[0].id
  host    = "internal.example.com"
  port    = 443

  # Optional. Reference user / group IDs returned by your IdP integration.
  # users  = ["user-id-1"]
  # groups = ["group-id-1"]
}
