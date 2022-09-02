resource "google_compute_network" "default" {
  name                    = var.service_name
  auto_create_subnetworks = true
}

resource "google_compute_address" "default" {
  name         = "ipv4-address"
  region       = var.region
  address_type = "EXTERNAL"
}