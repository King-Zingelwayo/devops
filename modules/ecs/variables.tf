variable "environment" {
  description = "Environment name"
  type        = string
}

variable "project_name" {
  description = "Project name for resource naming"
  type        = string
  default     = "portfolio"
}

variable "vpc_id" {
  description = "ID of the VPC"
  type        = string
}

variable "private_subnet_ids" {
  description = "IDs of the private subnets"
  type        = list(string)
}

variable "public_subnet_ids" {
  description = "IDs of the public subnets"
  type        = list(string)
}

variable "services" {
  description = "Map of services to create"
  type = map(object({
    image       = string
    port        = number
    cpu         = optional(number, 256)
    memory      = optional(number, 512)
    environment = optional(map(string), {})
    expose_alb  = optional(bool, false)
    host_header = optional(string, null)
  }))
  default = {}
}

variable "container_insights_enabled" {
  description = "Enable ECS container insights"
  type        = bool
  default     = false
}

variable "alb_internal" {
  description = "Whether ALB is internal"
  type        = bool
  default     = false
}

variable "deletion_protection_enabled" {
  description = "Enable deletion protection for ALB"
  type        = bool
  default     = false
}

variable "frontend_port" {
  description = "Port for frontend service"
  type        = number
  default     = 8080
}

variable "game_service_port" {
  description = "Port for game service"
  type        = number
  default     = 8081
}

variable "health_check_path" {
  description = "Health check path"
  type        = string
  default     = "/health"
}

variable "health_check_healthy_threshold" {
  description = "Health check healthy threshold"
  type        = number
  default     = 2
}

variable "health_check_unhealthy_threshold" {
  description = "Health check unhealthy threshold"
  type        = number
  default     = 2
}

variable "health_check_timeout" {
  description = "Health check timeout"
  type        = number
  default     = 5
}

variable "health_check_interval" {
  description = "Health check interval"
  type        = number
  default     = 30
}

variable "task_cpu" {
  description = "CPU units for ECS tasks"
  type        = number
  default     = 256
}

variable "task_memory" {
  description = "Memory for ECS tasks"
  type        = number
  default     = 512
}

variable "desired_count" {
  description = "Desired number of tasks"
  type        = number
  default     = 1
}

variable "log_retention_days" {
  description = "CloudWatch log retention in days"
  type        = number
  default     = 7
}

variable "tags" {
  description = "Additional tags for resources"
  type        = map(string)
  default     = {}
}

variable "kms_key_id" {
  description = "KMS key ID for log encryption"
  type        = string
  default     = null
}

variable "certificate_arn" {
  description = "ACM certificate ARN for HTTPS"
  type        = string
}

variable "enable_monitoring" {
  description = "Enable monitoring stack (CloudWatch, Grafana, Prometheus)"
  type = object({
    cloudwatch = optional(bool, true)
    grafana    = optional(bool, false)
    grafana_config = optional(object({
      admin_password  = optional(string, "admin123")
      dashboard_files = optional(list(string), [])
      domain          = optional(string, null)
    }), null)
    prometheus_config_file = optional(string, null)
  })
  default = {
    cloudwatch = true
    grafana    = false
  }

  validation {
    condition     = var.enable_monitoring.grafana == false || (var.enable_monitoring.grafana == true && var.enable_monitoring.grafana_config != null && length(var.enable_monitoring.grafana_config.dashboard_files) > 0)
    error_message = "When enable_monitoring.grafana is true, enable_monitoring.grafana_config.dashboard_files must be provided with at least one dashboard JSON file."
  }
}

variable "domain_name" {
  description = "Domain name for host-based routing"
  type        = string
}