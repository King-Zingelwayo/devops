terraform {
  backend "s3" {
    bucket = "tf-state-king-zinge-devops-71583"
    key    = "main/terraform.tfstate"
    region = "af-south-1"
  }
}