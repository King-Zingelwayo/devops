# ECS Cluster
resource "aws_ecs_cluster" "main" {
  name = "${var.environment}-${var.project_name}-cluster"

  setting {
    name  = "containerInsights"
    value = var.container_insights_enabled ? "enabled" : "disabled"
  }

  tags = merge(var.tags, {
    Name        = "${var.environment}-${var.project_name}-cluster"
    Environment = var.environment
  })
}

# Application Load Balancer
resource "aws_lb" "main" {
  name               = "${var.environment}-${var.project_name}-alb"
  internal           = var.alb_internal
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = var.public_subnet_ids

  enable_deletion_protection = var.environment == "prod"
  drop_invalid_header_fields = true

  tags = merge(var.tags, {
    Name        = "${var.environment}-${var.project_name}-alb"
    Environment = var.environment
  })
}

resource "aws_lb_target_group" "services" {
  for_each = { for k, v in var.services : k => v if v.expose_alb }

  name     = "${var.environment}-${substr(each.key, 0, 20)}-tg"
  port     = each.value.port
  protocol = "HTTP"
  vpc_id   = var.vpc_id
  target_type = "ip"

  health_check {
    enabled             = true
    healthy_threshold   = var.health_check_healthy_threshold
    interval            = var.health_check_interval
    matcher             = "200"
    path                = var.health_check_path
    port                = "traffic-port"
    protocol            = "HTTP"
    timeout             = var.health_check_timeout
    unhealthy_threshold = var.health_check_unhealthy_threshold
  }

  tags = merge(var.tags, {
    Name        = "${var.environment}-${var.project_name}-${each.key}-tg"
    Environment = var.environment
  })
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.main.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

resource "aws_lb_listener" "https" {
  count = var.certificate_arn != null ? 1 : 0

  load_balancer_arn = aws_lb.main.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS-1-2-2017-01"
  certificate_arn   = var.certificate_arn

  default_action {
    type = "fixed-response"
    fixed_response {
      content_type = "text/plain"
      message_body = "Not Found"
      status_code  = "404"
    }
  }
}

resource "aws_lb_listener_rule" "services" {
  for_each = { for k, v in var.services : k => v if v.expose_alb }

  listener_arn = aws_lb_listener.https[0].arn
  priority     = 50 + index(keys({ for k, v in var.services : k => v if v.expose_alb }), each.key)

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.services[each.key].arn
  }

  condition {
    path_pattern {
      values = ["/*"]
    }
  }
}

# Security Groups
resource "aws_security_group" "alb" {
  name_prefix = "${var.environment}-${var.project_name}-alb-"
  vpc_id      = var.vpc_id

  ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "HTTPS"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, {
    Name        = "${var.environment}-${var.project_name}-alb-sg"
    Environment = var.environment
  })
}

resource "aws_security_group" "ecs_tasks" {
  name_prefix = "${var.environment}-${var.project_name}-ecs-tasks-"
  vpc_id      = var.vpc_id

  dynamic "ingress" {
    for_each = { for k, v in var.services : k => v if v.expose_alb }
    content {
      description     = "HTTP from ALB for ${ingress.key}"
      from_port       = ingress.value.port
      to_port         = ingress.value.port
      protocol        = "tcp"
      security_groups = [aws_security_group.alb.id]
    }
  }

  ingress {
    description     = "Grafana from ALB"
    from_port       = 3000
    to_port         = 3000
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  ingress {
    description = "Service to service"
    from_port   = 0
    to_port     = 65535
    protocol    = "tcp"
    self        = true
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, {
    Name        = "${var.environment}-${var.project_name}-ecs-tasks-sg"
    Environment = var.environment
  })
}

# IAM Roles
resource "aws_iam_role" "ecs_task_execution_role" {
  name = "${var.environment}-${var.project_name}-ecs-task-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })

  tags = merge(var.tags, {
    Name        = "${var.environment}-${var.project_name}-ecs-task-execution-role"
    Environment = var.environment
  })
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_role_policy" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role" "ecs_task_role" {
  name = "${var.environment}-${var.project_name}-ecs-task-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })

  tags = merge(var.tags, {
    Name        = "${var.environment}-${var.project_name}-ecs-task-role"
    Environment = var.environment
  })
}

# CloudWatch Log Groups
resource "aws_cloudwatch_log_group" "services" {
  for_each = var.services

  name              = "/ecs/${var.environment}-${var.project_name}-${each.key}"
  retention_in_days = var.log_retention_days
  kms_key_id        = var.kms_key_id

  tags = merge(var.tags, {
    Name        = "${var.environment}-${var.project_name}-${each.key}-logs"
    Environment = var.environment
  })
}

# ECS Task Definitions
resource "aws_ecs_task_definition" "services" {
  for_each = var.services

  family                   = "${var.environment}-${var.project_name}-${each.key}"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = each.value.cpu
  memory                   = each.value.memory
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn           = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([
    {
      name  = each.key
      image = each.value.image
      
      portMappings = [
        {
          containerPort = each.value.port
          protocol      = "tcp"
        }
      ]

      environment = [
        for k, v in merge({
          PORT = tostring(each.value.port)
          LOG_LEVEL = "info"
        }, each.value.environment) : {
          name  = k
          value = v
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.services[each.key].name
          awslogs-region        = data.aws_region.current.name
          awslogs-stream-prefix = "ecs"
        }
      }

      healthCheck = {
        command     = ["CMD-SHELL", "wget --quiet --tries=1 --spider http://localhost:${each.value.port}/health || exit 1"]
        interval    = 30
        timeout     = 5
        retries     = 3
        startPeriod = 60
      }
    }
  ])

  tags = merge(var.tags, {
    Name        = "${var.environment}-${var.project_name}-${each.key}-task"
    Environment = var.environment
  })
}

# Service Discovery
resource "aws_service_discovery_private_dns_namespace" "main" {
  name        = "${var.environment}-${var.project_name}.local"
  description = "Service discovery namespace for ${var.project_name} services"
  vpc         = var.vpc_id

  tags = merge(var.tags, {
    Name        = "${var.environment}-${var.project_name}-namespace"
    Environment = var.environment
  })
}

resource "aws_service_discovery_service" "services" {
  for_each = { for k, v in var.services : k => v if !v.expose_alb }

  name = each.key

  dns_config {
    namespace_id = aws_service_discovery_private_dns_namespace.main.id

    dns_records {
      ttl  = 10
      type = "A"
    }

    routing_policy = "MULTIVALUE"
  }

  tags = merge(var.tags, {
    Name        = "${var.environment}-${var.project_name}-${each.key}-discovery"
    Environment = var.environment
  })
}

# ECS Services
resource "aws_ecs_service" "services" {
  for_each = var.services

  name            = "${var.environment}-${var.project_name}-${each.key}"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.services[each.key].arn
  desired_count   = var.desired_count
  launch_type     = "FARGATE"

  network_configuration {
    security_groups = [aws_security_group.ecs_tasks.id]
    subnets         = var.private_subnet_ids
  }

  dynamic "service_registries" {
    for_each = each.value.expose_alb ? [] : [1]
    content {
      registry_arn = aws_service_discovery_service.services[each.key].arn
    }
  }

  dynamic "load_balancer" {
    for_each = each.value.expose_alb ? [1] : []
    content {
      target_group_arn = aws_lb_target_group.services[each.key].arn
      container_name   = each.key
      container_port   = each.value.port
    }
  }

  depends_on = [aws_lb_listener.http]

  tags = merge(var.tags, {
    Name        = "${var.environment}-${var.project_name}-${each.key}"
    Environment = var.environment
  })
}


