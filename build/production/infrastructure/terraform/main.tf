terraform {
  required_version = ">= 1.5.0"
  
  # Audit Protocol: Remote GCS backend configured natively ensuring shared-state synchronization
  backend "gcs" {
    bucket  = "stadium-tf-state-secure"
    prefix  = "terraform/state/stadium-api"
  }
  
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

variable "project_id" {
  type    = string
  default = "stadium-experience-loc"
}

variable "docker_image" {
  type = string
}

provider "google" {
  project = var.project_id
  region  = "us-central1"
}

# 1. VPC Network strictly natively optimized for Serverless boundaries
resource "google_compute_network" "stadium_vpc" {
  name                    = "stadium-vpc"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "stadium_subnet" {
  name          = "stadium-serverless-subnet"
  ip_cidr_range = "10.0.0.0/28"
  region        = "us-central1"
  network       = google_compute_network.stadium_vpc.id
}

# 2. Serverless VPC Access Connector natively bridging Cloud Run directly to MemoryStore limits
resource "google_vpc_access_connector" "redis_bridge" {
  name          = "stadium-redis-bridge"
  region        = "us-central1"
  network       = google_compute_network.stadium_vpc.name
  ip_cidr_range = "10.8.0.0/28"
  machine_type  = "e2-micro"
  min_instances = 2
  max_instances = 10
}

# 3. High-Availability Memorystore (Redis) executing the Ephemeral Cache block
resource "google_redis_instance" "ephemeral_buffer" {
  name               = "stadium-redis-ha"
  tier               = "STANDARD_HA"
  memory_size_gb     = 1
  region             = "us-central1"
  authorized_network = google_compute_network.stadium_vpc.id
  redis_version      = "REDIS_7_0"
  connect_mode       = "DIRECT_PEERING"
}

# 4. Native Cloud Run v2 API Platform natively integrating parameters from service.yaml
resource "google_cloud_run_v2_service" "stadium_backend" {
  name     = "stadium-experience-backend"
  location = "us-central1"
  ingress  = "INGRESS_TRAFFIC_INTERNAL_LOAD_BALANCER"

  template {
    scaling {
      min_instance_count = 2
      max_instance_count = 100
    }
    
    vpc_access {
      connector = google_vpc_access_connector.redis_bridge.id
      egress    = "PRIVATE_RANGES_ONLY"
    }

    containers {
      image = var.docker_image
      
      env {
        name  = "GOOGLE_CLOUD_PROJECT"
        value = var.project_id
      }
      
      # Implicit Zero-Trust Protocol: Reading Redis metrics directly via Secret Manager locally avoiding explicit TF exposures
      env {
        name = "REDIS_HOST"
        value_source {
          secret_key_ref {
            secret  = "stadium-redis-host"
            version = "latest"
          }
        }
      }
      env {
        name = "REDIS_PORT"
        value_source {
          secret_key_ref {
            secret  = "stadium-redis-port"
            version = "latest"
          }
        }
      }
      
      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }
    }
  }
}

output "redis_internal_ip" {
  value     = google_redis_instance.ephemeral_buffer.host
  sensitive = true
}

output "redis_port" {
  value     = google_redis_instance.ephemeral_buffer.port
  sensitive = true
}
