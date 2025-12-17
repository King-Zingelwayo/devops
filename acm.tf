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

# ACM Certificate Validation (skip DNS record creation)
resource "aws_acm_certificate_validation" "main" {
  certificate_arn = aws_acm_certificate.main.arn
  
  timeouts {
    create = "10m"
  }
}
