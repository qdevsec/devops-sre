provider "aws" {
  region = "us-east-1"
}

module "eks" {
  source          = "terraform-aws-modules/eks/aws"
  cluster_name    = "ha-logviewer-cluster"
  cluster_version = "1.29"
  subnets         = ["subnet-123456", "subnet-789012"]
  vpc_id          = "vpc-abcdef"
}

resource "kubernetes_namespace" "web_ha" {
  metadata {
    name = "web-ha"
  }
}

resource "kubernetes_config_map" "logviewer_html" {
  metadata {
    name      = "config-logviewer-html"
    namespace = kubernetes_namespace.web_ha.metadata[0].name
  }
  data = {
    "index.html" = file("${path.module}/index.html")
    "cloud.json" = file("${path.module}/cloud.json")
    "sysmon.json" = file("${path.module}/sysmon.json")
  }
}

resource "kubernetes_deployment" "ha_logviewer" {
  metadata {
    name      = "ha-logviewer"
    namespace = kubernetes_namespace.web_ha.metadata[0].name
    labels = {
      app = "ha-logviewer"
    }
  }
  spec {
    replicas = 1
    selector {
      match_labels = {
        app = "ha-logviewer"
      }
    }
    template {
      metadata {
        labels = {
          app = "ha-logviewer"
        }
      }
      spec {
        container {
          name  = "logviewer"
          image = "logviewer:latest"
          port {
            container_port = 8080
          }
          volume_mount {
            name       = "logviewer-html"
            mount_path = "/placeholders"
          }
        }
        volume {
          name = "logviewer-html"
          config_map {
            name = kubernetes_config_map.logviewer_html.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "ha_logviewer" {
  metadata {
    name      = "ha-logviewer-service"
    namespace = kubernetes_namespace.web_ha.metadata[0].name
  }
  spec {
    selector = {
      app = "ha-logviewer"
    }
    port {
      port        = 8080
      target_port = 8080
      node_port   = 30080
    }
    type = "NodePort"
  }
}
