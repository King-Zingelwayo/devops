# ACM Certificate
resource "aws_acm_certificate" "main" {
  domain_name               = data.cloudflare_zone.this.name
  subject_alternative_names = ["*.${data.cloudflare_zone.this.name}"]
  validation_method         = "DNS"

  lifecycle {
    create_before_destroy = true
  }

  tags = {
    Name        = "${var.environment}-${var.project_name}-cert"
    Environment = var.environment
  }
}

# Cloudflare DNS records for ACM validation
resource "cloudflare_record" "acm_validation" {
  for_each = {
    for dvo in aws_acm_certificate.main.domain_validation_options : dvo.domain_name => {
      name  = dvo.resource_record_name
      type  = dvo.resource_record_type
      value = dvo.resource_record_value
    }
  }

  zone_id         = var.cloudflare_zone_id
  name            = each.value.name
  type            = each.value.type
  content         = each.value.value
  ttl             = 60
  proxied         = false
  allow_overwrite = true
}

# ACM Certificate Validation
resource "aws_acm_certificate_validation" "main" {
  certificate_arn         = aws_acm_certificate.main.arn
  validation_record_fqdns = [for record in cloudflare_record.acm_validation : record.hostname]
}
