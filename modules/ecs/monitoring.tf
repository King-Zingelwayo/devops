# Monitoring Stack - Grafana & Prometheus
locals {
  monitoring_enabled = var.enable_monitoring.grafana
}

# CloudWatch Log Group for Grafana
resource "aws_cloudwatch_log_group" "grafana" {
  count = local.monitoring_enabled ? 1 : 0

  name              = "/ecs/${var.environment}-${var.project_name}-grafana"
  retention_in_days = var.log_retention_days
  kms_key_id        = var.kms_key_id

  tags = merge(var.tags, {
    Name        = "${var.environment}-${var.project_name}-grafana-logs"
    Environment = var.environment
  })
}

# CloudWatch Log Group for Prometheus
resource "aws_cloudwatch_log_group" "prometheus" {
  count = local.monitoring_enabled ? 1 : 0

  name              = "/ecs/${var.environment}-${var.project_name}-prometheus"
  retention_in_days = var.log_retention_days
  kms_key_id        = var.kms_key_id

  tags = merge(var.tags, {
    Name        = "${var.environment}-${var.project_name}-prometheus-logs"
    Environment = var.environment
  })
}

# Monitoring Task Definition
resource "aws_ecs_task_definition" "monitoring" {
  count = local.monitoring_enabled ? 1 : 0

  family                   = "${var.environment}-${var.project_name}-monitoring"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = 512
  memory                   = 1024
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn           = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([
    {
      name  = "grafana"
      image = "grafana/grafana:10.0.0"
      
      portMappings = [
        {
          containerPort = 3000
          protocol      = "tcp"
        }
      ]

      environment = [
        {
          name  = "GF_SECURITY_ADMIN_PASSWORD"
          value = var.enable_monitoring.grafana_config.admin_password
        },
        {
          name  = "GF_USERS_ALLOW_SIGN_UP"
          value = "false"
        },
        {
          name  = "GF_DATASOURCES_DEFAULT_URL"
          value = "http://localhost:9090"
        }
      ]

      healthCheck = {
        command     = ["CMD-SHELL", "wget --quiet --tries=1 --spider http://localhost:3000/api/health || exit 1"]
        interval    = 30
        timeout     = 5
        retries     = 3
        startPeriod = 60
      }

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.grafana[0].name
          awslogs-region        = data.aws_region.current.name
          awslogs-stream-prefix = "ecs"
        }
      }
    },
    {
      name  = "prometheus"
      image = "prom/prometheus:v2.40.0"
      
      portMappings = [
        {
          containerPort = 9090
          protocol      = "tcp"
        }
      ]

      healthCheck = {
        command     = ["CMD-SHELL", "wget --quiet --tries=1 --spider http://localhost:9090/-/healthy || exit 1"]
        interval    = 30
        timeout     = 5
        retries     = 3
        startPeriod = 60
      }

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.prometheus[0].name
          awslogs-region        = data.aws_region.current.name
          awslogs-stream-prefix = "ecs"
        }
      }
    }
  ])

  tags = merge(var.tags, {
    Name        = "${var.environment}-${var.project_name}-monitoring-task"
    Environment = var.environment
  })
}

# Target Group for Grafana
resource "aws_lb_target_group" "grafana" {
  count = local.monitoring_enabled ? 1 : 0

  name     = "${var.environment}-grafana-tg"
  port     = 3000
  protocol = "HTTP"
  vpc_id   = var.vpc_id
  target_type = "ip"

  health_check {
    enabled             = true
    healthy_threshold   = var.health_check_healthy_threshold
    interval            = var.health_check_interval
    matcher             = "200"
    path                = "/api/health"
    port                = "traffic-port"
    protocol            = "HTTP"
    timeout             = var.health_check_timeout
    unhealthy_threshold = var.health_check_unhealthy_threshold
  }

  tags = merge(var.tags, {
    Name        = "${var.environment}-${var.project_name}-grafana-tg"
    Environment = var.environment
  })
}

# ALB Listener Rule for Grafana
resource "aws_lb_listener_rule" "grafana" {
  count = local.monitoring_enabled && var.certificate_arn != null && var.enable_monitoring.grafana_config.domain != null ? 1 : 0

  listener_arn = aws_lb_listener.https[0].arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.grafana[0].arn
  }

  condition {
    host_header {
      values = [var.enable_monitoring.grafana_config.domain]
    }
  }
}

# ECS Service for Monitoring
resource "aws_ecs_service" "monitoring" {
  count = local.monitoring_enabled ? 1 : 0

  name            = "${var.environment}-${var.project_name}-monitoring"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.monitoring[0].arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    security_groups = [aws_security_group.ecs_tasks.id]
    subnets         = var.private_subnet_ids
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.grafana[0].arn
    container_name   = "grafana"
    container_port   = 3000
  }

  depends_on = [aws_lb_listener.http]

  tags = merge(var.tags, {
    Name        = "${var.environment}-${var.project_name}-monitoring"
    Environment = var.environment
  })
}
