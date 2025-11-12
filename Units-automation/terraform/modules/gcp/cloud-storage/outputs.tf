output "name" {
  value = google_storage_bucket.storage_bucket.name
}

output "google_backups_bucket" {
  value = google_storage_bucket.google_backups_bucket.name
}