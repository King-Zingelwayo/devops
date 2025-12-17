# Portfolio Wordle Platform - Architecture Documentation

## System Overview

This is a production-grade, cost-optimized Wordle-like game platform built as cloud-native microservices on AWS.

## Architecture Decisions

### 1. Monitoring Strategy: Hybrid CloudWatch + Prometheus

**Decision**: Use CloudWatch for infrastructure metrics and Prometheus for application metrics.

**Rationale**:
- **CloudWatch**: Native AWS integration, no additional infrastructure cost
- **Prometheus**: Better for custom application metrics, more cost-effective for high-frequency data
- **Grafana**: Single pane of glass for both data sources

**Cost Impact**: ~$20-25/month vs $100+ with full CloudWatch monitoring

### 2. Cost Optimization Strategies

#### CloudWatch Optimizations
- **Container Insights**: DISABLED (saves ~$30/month)
- **Log Retention**: 7 days instead of 30 days default
- **Metric Resolution**: Standard 60s intervals only
- **Selective Metrics**: Only essential ECS and ALB metrics

#### ECS Optimizations
- **Task Size**: 256 CPU / 512 MB memory (minimal viable)
- **Desired Count**: 1 instance per service (cost over availability)
- **Fargate**: No EC2 management overhead

#### Application Optimizations
- **Low Cardinality Metrics**: No user IDs or game IDs in Prometheus labels
- **Efficient Queries**: 30s+ refresh intervals in Grafana
- **Minimal Dependencies**: Lean Go services with distroless containers

### 3. Security Architecture

```
Internet → ALB → ECS Tasks (Private Subnets) → Service Discovery
```

- **Private Subnets**: All application containers run in private subnets
- **Security Groups**: Least privilege access (ALB → ECS only)
- **IAM Roles**: Task-specific permissions, no broad access
- **Service Discovery**: Internal DNS for service-to-service communication

### 4. Microservices Design

#### Game Service
- **Responsibility**: Core Wordle logic, game state management
- **API**: REST endpoints for game operations
- **Metrics**: Business metrics (games started, win rate, etc.)
- **Storage**: In-memory (stateless design)

#### Frontend Service
- **Responsibility**: UI serving and API proxying
- **Features**: Embedded HTML/CSS/JS, API gateway pattern
- **Caching**: Static asset serving with proper headers

#### Monitoring Service
- **Grafana**: Dashboards and alerting
- **Prometheus**: Application metrics collection (optional)
- **CloudWatch**: Infrastructure metrics

## Infrastructure Components

### Network Layer
- **VPC**: Isolated network environment
- **Subnets**: Public (ALB) + Private (ECS tasks)
- **NAT Gateways**: Outbound internet access for private subnets
- **Security Groups**: Network-level access control

### Compute Layer
- **ECS Fargate**: Serverless container platform
- **Application Load Balancer**: HTTP traffic distribution
- **Service Discovery**: DNS-based service location

### Monitoring Layer
- **CloudWatch Logs**: Centralized logging with retention policies
- **CloudWatch Metrics**: Infrastructure monitoring
- **Prometheus**: Application metrics (optional)
- **Grafana**: Visualization and dashboards

## Deployment Pipeline

### CI/CD Flow
1. **Code Push** → GitHub Actions triggered
2. **Build & Test** → Go tests and compilation
3. **Container Build** → Multi-stage Docker builds
4. **ECR Push** → Container registry storage
5. **Terraform Plan** → Infrastructure change preview
6. **Manual Approval** → Production deployment gate
7. **Terraform Apply** → Infrastructure deployment

### Environments
- **Local**: Docker Compose for development
- **Development**: AWS with cost-optimized settings
- **Production**: AWS with enhanced monitoring and redundancy

## Cost Analysis

### Monthly Cost Breakdown (Estimated)
- **ECS Fargate**: ~$15-20 (2 tasks, minimal CPU/memory)
- **ALB**: ~$16 (base cost)
- **NAT Gateways**: ~$32 (2 AZs)
- **CloudWatch**: ~$5-10 (optimized settings)
- **Monitoring**: ~$5 (Grafana container)
- **Total**: ~$73-83/month

### Cost vs Full Monitoring
- **This Architecture**: ~$80/month
- **Full CloudWatch**: ~$150+/month
- **Savings**: ~$70/month (47% reduction)

## Production Readiness

### Reliability
- **Health Checks**: All services have health endpoints
- **Graceful Shutdown**: Proper signal handling
- **Circuit Breakers**: Retry logic with exponential backoff
- **Service Discovery**: Automatic service location

### Observability
- **Structured Logging**: JSON format for parsing
- **Metrics**: Business and technical metrics
- **Dashboards**: Real-time monitoring
- **Alerting**: Critical threshold notifications

### Security
- **Least Privilege**: Minimal IAM permissions
- **Network Isolation**: Private subnets
- **Container Security**: Distroless base images
- **Secrets Management**: Environment variables (extend with AWS Secrets Manager)

## Trade-offs Made

### Cost vs Availability
- **Single Instance**: Cost savings over high availability
- **Minimal Resources**: Right-sized for demo workload
- **Shared Infrastructure**: Multi-tenant approach

### Simplicity vs Features
- **In-Memory Storage**: No database costs, but no persistence
- **Embedded UI**: No separate frontend deployment
- **Basic Auth**: Simple Grafana authentication

### Monitoring vs Cost
- **Selective Metrics**: Essential monitoring only
- **Short Retention**: 7-day log retention
- **Standard Resolution**: 60s metric intervals

## Future Enhancements

### Scalability
- **Auto Scaling**: CPU/memory-based scaling
- **Database**: RDS or DynamoDB for persistence
- **CDN**: CloudFront for static assets

### Reliability
- **Multi-AZ**: Cross-availability zone deployment
- **Backup**: Database backup strategies
- **Disaster Recovery**: Cross-region replication

### Security
- **WAF**: Web Application Firewall
- **Secrets Manager**: Centralized secret management
- **VPC Endpoints**: Private AWS service access