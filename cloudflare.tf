# Cloudflare record for portfolio app
resource "cloudflare_record" "portfolio_app" {
  zone_id = var.cloudflare_zone_id
  name    = "sihle"
  content = module.ecs.alb_dns_name
  type    = "CNAME"
  proxied = true
  ttl     = 1
}

# Cloudflare record for monitoring (Grafana)
resource "cloudflare_record" "monitoring" {
  zone_id = var.cloudflare_zone_id
  name    = "monitoring"
  content = module.ecs.monitoring_alb_dns_name
  type    = "CNAME"
  proxied = true
  ttl     = 1
}
