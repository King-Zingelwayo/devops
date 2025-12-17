# ECR Repositories
resource "aws_ecr_repository" "game_service" {
  name                 = "${var.project_name}-game-service"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = var.tags
}

resource "aws_ecr_repository" "frontend_service" {
  name                 = "${var.project_name}-frontend-service"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = var.tags
}

# Build and push Docker images
resource "null_resource" "build_and_push_images" {
  depends_on = [
    aws_ecr_repository.game_service,
    aws_ecr_repository.frontend_service
  ]

  provisioner "local-exec" {
    command = <<-EOT
      # Login to ECR
      aws ecr get-login-password --region ${var.region} | docker login --username AWS --password-stdin ${data.aws_caller_identity.current.account_id}.dkr.ecr.${var.region}.amazonaws.com

      # Build and push game service
      cd  services/game-service/
      docker build -t ${aws_ecr_repository.game_service.repository_url}:latest .
      docker push ${aws_ecr_repository.game_service.repository_url}:latest

      # Build and push frontend service
      cd services/frontend-service/
      docker build -t ${aws_ecr_repository.frontend_service.repository_url}:latest .
      docker push ${aws_ecr_repository.frontend_service.repository_url}:latest
    EOT
  }

  triggers = {
    always_run = timestamp()
  }
}