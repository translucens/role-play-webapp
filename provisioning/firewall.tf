resource "google_compute_firewall" "allow_ingress_from_internet" {
  name    = "allow-ingress-from-internet"
  network = google_compute_network.default.name

  allow {
    protocol = "tcp"
    ports = ["22","80"]
  }
  direction = "INGRESS"
  source_ranges = ["0.0.0.0/0"]
}