terraform {
    required_providers {
        docker = {
            source = "kreuzwerker/docker"
            version = "2.7.2"
        }
    }
}

provider "docker" {}

resource "docker_network" "app_network" {
    name = "app_network"
}

resource "docker_image" "postgres" {
    name = "postgres:latest"
    keep_locally = false
}

resource "docker_container" "postgres" {
    name = "postgres"
    image  = docker_image.postgres.latest
    networks_advanced {
        name = docker_network.app_network.name
    }
    env = [
        "POSTGRES_USER=admin",
        "POSTGRES_PASSWORD=admin",
        "POSTGRES_DB=appdb"
    ]
    ports {
        internal = 5432
        external = 5432
    }
}

resource "docker_image" "nginx" {
    name = "nginx:latest"
    keep_locally = false
}

resource "docker_container" "nginx" {
    name = "nginx"
    image = docker_image.nginx.latest
    depends_on = [docker_container.postgres]
    networks_advanced {
        name = docker_network.app_network.name
    }
    ports {
        internal = 80
        external = 8080
    }
}