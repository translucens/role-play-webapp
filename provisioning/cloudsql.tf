resource "random_id" "db_instance_name_suffix" {
  byte_length = 4
}

resource "google_sql_database_instance" "scstore" {
  name             = "${var.service_name}-${random_id.db_instance_name_suffix.hex}"
  database_version = "POSTGRES_14"

  depends_on = [google_service_networking_connection.private_vpc_connection]

  settings {
    tier              = var.cloudsql_machine_type
    availability_type = "REGIONAL"

    ip_configuration {
      ipv4_enabled        = false
      private_network     = google_compute_network.default.id
      require_ssl         = false
      allocated_ip_range  = null
    }
  }
}

resource "google_sql_database" "scstore" {
  name     = var.service_name
  instance = google_sql_database_instance.scstore.name
}

resource "google_sql_user" "scstore" {
  name     = var.service_name
  password = var.database_scstore_password
  instance = google_sql_database_instance.scstore.name
}