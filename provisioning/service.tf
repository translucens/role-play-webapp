locals {
    services = toset([
        "compute.googleapis.com",
        "sqladmin.googleapis.com",
        "iam.googleapis.com",
        "servicenetworking.googleapis.com",
   ])
}

resource "google_project_service" "service" {
    for_each = local.services
    project = var.project_id
    service = each.value
}