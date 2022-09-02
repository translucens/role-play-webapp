resource "google_compute_instance" "scstore" {
  name         = var.service_name
  machine_type = var.gce_machine_type
  zone         = var.zone
  tags         = ["http-server"]
  boot_disk {
    initialize_params {
      image = var.gce_image_name
    }
  }
  network_interface {
    network = google_compute_network.default.self_link
    access_config {
      nat_ip = google_compute_address.default.address
    }
  }

  metadata = {
    enable-oslogin = "TRUE"
  }

  metadata_startup_script = file("${path.module}/startup.sh")
}