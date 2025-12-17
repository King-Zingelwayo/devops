#cloudflare hosted zone
data "cloudflare_zone" "this" {
  zone_id = var.cloudflare_zone_id
}
data "aws_caller_identity" "current" {}