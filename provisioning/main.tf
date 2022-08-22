terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "3.5.0"
    }
  }
}

provider "google" {
  #credentials = file(var.credentials_file)
  project = var.project
}

resource "google_compute_network" "default" {
  name = "scstore"
}

resource "google_compute_subnetwork" "default" {
  name          = "my-subnet"
  region        = "asia-northeast1"
  network       = google_compute_network.default.id
  ip_cidr_range = "10.0.0.0/16"
}

resource "google_compute_address" "default" {
  name         = "ipv4-address"
  region       = "asia-northeast1"
  address_type = "EXTERNAL"
}

resource "google_project_service" "enable_api" {
  service = "compute.googleapis.com"
}

resource "google_compute_instance" "vm_instance" {
  name         = "scstore"
  machine_type = "n1-standard-1"
  zone = "asia-northeast1-c"
  tags = ["http-server","https-server"]
  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-10"
    }
  }
  network_interface {
    network = google_compute_network.default.self_link
    access_config {
      nat_ip = google_compute_address.default.address
    }
  }

  metadata = {
    enable-oslogin="TRUE"
  }

 metadata_startup_script = "/startup.sh"
}



  


