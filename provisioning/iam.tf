resource "google_project_iam_member" "benchmark" {
  project = var.project_id
  role    = "roles/viewer"
  member  = "serviceAccount:${var.benchmark_service_account_name}"
}