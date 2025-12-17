output "cluster_name" {
  description = "Name of the ECS cluster"
  value       = aws_ecs_cluster.main.name
}

output "alb_dns_name" {
  description = "DNS name of the load balancer"
  value       = aws_lb.main.dns_name
}

output "monitoring_alb_dns_name" {
  description = "DNS name for monitoring (Grafana) access"
  value       = aws_lb.main.dns_name
}