# Sihle Ndlovu - Portfolio & Cloud Architecture Demo

## Professional Portfolio & Interactive Demo

This platform showcases **Sihle Ndlovu's** cloud engineering expertise through:
- **Professional CV** with AWS certifications and experience
- **Interactive Space Invaders Game** demonstrating microservices architecture
- **Production-grade Infrastructure** with cost optimization

## Cloud Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Portfolio     │    │   Game Service   │    │   Monitoring    │
│   Frontend      │    │   (Go Backend)   │    │   (Grafana)     │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                        │                        │
         └────────┬───────────────┴────────────────────────┘
                  │
         ┌─────────────────┐
         │ Application     │
         │ Load Balancer   │
         └─────────────────┘
                  │
         ┌─────────────────┐
         │   ECS Cluster   │
         │   (Fargate)     │
         └─────────────────┘
```

## Monitoring Strategy & Cost Justification

### Hybrid Approach: CloudWatch + Prometheus + Grafana

**CloudWatch (Infrastructure)**
- ECS container insights: DISABLED (saves ~$30/month)
- Standard resolution metrics only (60s intervals)
- Essential ALB metrics only
- Log retention: 7 days (vs 30 days default)
- **Cost**: ~$5-10/month vs $50+ with full monitoring

**Prometheus (Application)**
- Custom business metrics (games played, success rate)
- Low-cardinality design (no user IDs in labels)
- Local storage with 7-day retention
- **Cost**: Container compute only (~$10/month)

**Grafana**
- Single dashboard for both data sources
- 30s+ refresh intervals
- Default time range: 1-6 hours
- **Cost**: Container compute only (~$5/month)

**Total Monitoring Cost**: ~$20-25/month vs $100+ with full CloudWatch

## Services

### 1. Game Service (Go)
- Core Wordle logic
- Daily word management
- Stateless design
- Metrics: `/metrics` endpoint

### 2. Frontend Service (Go)
- Serves static assets
- Health checks
- Minimal logging

### 3. Monitoring Stack
- Grafana container
- Prometheus container (optional)

## Quick Start

### Local Development
```bash
# Start services
docker-compose up -d

# Access
# Game: http://localhost:8080
# Grafana: http://localhost:3000 (admin/admin)
```

### AWS Deployment
```bash
# Initialize Terraform
cd terraform
terraform init

# Plan deployment
terraform plan -var="environment=prod"

# Deploy (with approval)
terraform apply
```

## Cost Controls Implemented

1. **CloudWatch Optimization**
   - Disabled container insights
   - 7-day log retention
   - Standard resolution metrics only
   - Selective metric collection

2. **ECS Optimization**
   - Fargate Spot instances where applicable
   - Right-sized task definitions
   - Auto-scaling based on CPU/memory

3. **Monitoring Efficiency**
   - Low-cardinality metrics
   - Efficient Grafana queries
   - Prometheus local storage

## Security

- Private subnets for ECS tasks
- Least privilege IAM roles
- Security groups with minimal access
- No hardcoded credentials

## Production Readiness

- Health checks on all services
- Structured JSON logging
- Graceful shutdown handling
- Circuit breaker patterns
- Retry logic with backoff