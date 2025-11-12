#############################
# Google Compute Engine Instance
#############################

resource "google_compute_instance" "this" {
  project             = var.project_id
  name                = var.vm_name
  zone                = var.vm_zone
  machine_type        = var.vm_machine_type
  deletion_protection = var.deletion_protection
  tags                = [var.vm_network_tag]

  boot_disk {
    initialize_params {
      image = var.vm_image
      size  = var.vm_disk_size
      type  = var.vm_disk_type
    }
    auto_delete = var.auto_delete_boot_disk
  }

  network_interface {
    network            = var.vm_network_self_link
    subnetwork         = var.vm_subnetwork_self_link

    dynamic "access_config" {
      for_each = var.enable_public_ip ? [1] : []
      content {
      }
    }
  }


  scheduling {
    preemptible         = var.preemptible
    automatic_restart   = var.automatic_restart
    on_host_maintenance = var.on_host_maintenance
  }

  lifecycle {
    prevent_destroy       = false
    create_before_destroy = false
  }
}
