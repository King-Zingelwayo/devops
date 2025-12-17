# Network Module
module "network" {
  source = "./modules/network"
  
  environment  = var.environment
  vpc_cidr     = var.vpc_cidr      
}

# ECS Module
module "ecs" {
  source = "./modules/ecs"
  
  environment         = var.environment
  vpc_id             = module.network.vpc_id
  private_subnet_ids = module.network.private_subnet_ids
  public_subnet_ids  = module.network.public_subnet_ids
  
  services = {
    game-service = {
      image       = "${aws_ecr_repository.game_service.repository_url}:latest"
      port        = 8081
      cpu         = 256
      memory      = 512
      environment = {}
      expose_alb  = false
    }
    frontend-service = {
      image       = "${aws_ecr_repository.frontend_service.repository_url}:latest"
      port        = 8080
      cpu         = 256
      memory      = 512
      environment = {
        GAME_SERVICE_URL = "http://game-service.${var.environment}-${var.project_name}.local:8081"
      }
      expose_alb  = true
    }
  }

  enable_monitoring = {
    cloudwatch = true
    grafana    = true
    grafana_config = {
      admin_password  = "admin123"
      dashboard_files = [
        "${path.module}/monitoring/grafana/dashboards/portfolio-dashboard.json",
        "${path.module}/monitoring/grafana/dashboards/portfolio-overview.json",
        "${path.module}/monitoring/grafana/dashboards/system-metrics.json"
      ]
      domain = "monitoring.${data.cloudflare_zone.this.name}"
    }
    prometheus_config_file = "${path.module}/monitoring/prometheus.yml"
  }
  certificate_arn = aws_acm_certificate.main.arn

  depends_on = [null_resource.build_and_push_images]
}

