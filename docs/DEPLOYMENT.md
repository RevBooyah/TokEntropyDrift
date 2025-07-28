# TokEntropyDrift Deployment Guide

This guide covers deploying TokEntropyDrift in various environments, from local development to production systems.

## Table of Contents

1. [Local Development](#local-development)
2. [Docker Deployment](#docker-deployment)
3. [Cloud Deployment](#cloud-deployment)
4. [Production Deployment](#production-deployment)
5. [Monitoring and Logging](#monitoring-and-logging)
6. [Security Considerations](#security-considerations)
7. [Troubleshooting](#troubleshooting)

## Local Development

### Prerequisites

- Go 1.22 or later
- Python 3.8+ (for some tokenizer adapters)
- Git

### Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/RevBooyah/TokEntropyDrift.git
   cd TokEntropyDrift
   ```

2. **Build the application**:
   ```bash
   go build -o ted cmd/ted/main.go
   ```

3. **Install dependencies** (if using Python tokenizers):
   ```bash
   pip install tiktoken transformers sentencepiece
   ```

4. **Configure the application**:
   ```bash
   cp ted.config.yaml.example ted.config.yaml
   # Edit ted.config.yaml with your preferences
   ```

5. **Test the installation**:
   ```bash
   ./ted --help
   ./ted analyze examples/english_quotes.txt --tokenizers=gpt2
   ```

### Development Workflow

1. **Run tests**:
   ```bash
   go test ./...
   ```

2. **Start development server**:
   ```bash
   ./ted serve --port=8080
   ```

3. **Run analysis**:
   ```bash
   ./ted analyze examples/english_quotes.txt --tokenizers=gpt2,bert --visualize
   ```

## Docker Deployment

### Single Container Deployment

1. **Create Dockerfile**:
   ```dockerfile
   FROM golang:1.22-alpine AS builder
   
   WORKDIR /app
   COPY . .
   RUN go build -o ted cmd/ted/main.go
   
   FROM alpine:latest
   RUN apk --no-cache add ca-certificates python3 py3-pip
   
   WORKDIR /root/
   COPY --from=builder /app/ted .
   COPY --from=builder /app/ted.config.yaml .
   COPY --from=builder /app/examples ./examples
   
   # Install Python dependencies
   RUN pip3 install tiktoken transformers sentencepiece
   
   EXPOSE 8080
   ENTRYPOINT ["./ted"]
   ```

2. **Build the image**:
   ```bash
   docker build -t tokentropydrift .
   ```

3. **Run the container**:
   ```bash
   # Run analysis
   docker run -v $(pwd)/data:/app/data -v $(pwd)/output:/app/output tokentropydrift analyze /app/data/input.txt --tokenizers=gpt2
   
   # Run web server
   docker run -p 8080:8080 -v $(pwd)/data:/app/data tokentropydrift serve --host=0.0.0.0 --port=8080
   ```

### Multi-Container Deployment

1. **Create docker-compose.yml**:
   ```yaml
   version: '3.8'
   
   services:
     ted-analysis:
       build: .
       volumes:
         - ./data:/app/data
         - ./output:/app/output
         - ./logs:/app/logs
       environment:
         - TED_LOGGING_LEVEL=info
         - TED_CACHE_ENABLED=true
       command: analyze /app/data/input.txt --tokenizers=gpt2,bert --output=/app/output
       restart: unless-stopped
   
     ted-dashboard:
       build: .
       ports:
         - "8080:8080"
       volumes:
         - ./data:/app/data
         - ./output:/app/output
         - ./logs:/app/logs
       environment:
         - TED_LOGGING_LEVEL=info
         - TED_CACHE_ENABLED=true
       command: serve --host=0.0.0.0 --port=8080
       restart: unless-stopped
       depends_on:
         - ted-analysis
   
     redis:
       image: redis:alpine
       ports:
         - "6379:6379"
       volumes:
         - redis_data:/data
       restart: unless-stopped
   
   volumes:
     redis_data:
   ```

2. **Deploy with docker-compose**:
   ```bash
   docker-compose up -d
   ```

3. **Monitor the deployment**:
   ```bash
   docker-compose logs -f
   ```

## Cloud Deployment

### AWS Deployment

#### EC2 Deployment

1. **Launch EC2 instance**:
   ```bash
   # Using AWS CLI
   aws ec2 run-instances \
     --image-id ami-0c02fb55956c7d316 \
     --instance-type t3.medium \
     --key-name your-key-pair \
     --security-group-ids sg-xxxxxxxxx \
     --subnet-id subnet-xxxxxxxxx \
     --user-data file://user-data.sh
   ```

2. **Create user-data.sh**:
   ```bash
   #!/bin/bash
   yum update -y
   yum install -y git go python3 python3-pip
   
   # Clone repository
   git clone https://github.com/RevBooyah/TokEntropyDrift.git
   cd TokEntropyDrift
   
   # Build application
   go build -o ted cmd/ted/main.go
   
   # Install Python dependencies
   pip3 install tiktoken transformers sentencepiece
   
   # Configure application
   cp ted.config.yaml.example ted.config.yaml
   
   # Start application
   nohup ./ted serve --host=0.0.0.0 --port=8080 > app.log 2>&1 &
   ```

3. **Configure security groups**:
   - Allow inbound traffic on port 8080
   - Allow SSH access on port 22

#### ECS Deployment

1. **Create task definition**:
   ```json
   {
     "family": "tokentropydrift",
     "networkMode": "awsvpc",
     "requiresCompatibilities": ["FARGATE"],
     "cpu": "512",
     "memory": "1024",
     "executionRoleArn": "arn:aws:iam::account:role/ecsTaskExecutionRole",
     "containerDefinitions": [
       {
         "name": "ted",
         "image": "your-account.dkr.ecr.region.amazonaws.com/tokentropydrift:latest",
         "portMappings": [
           {
             "containerPort": 8080,
             "protocol": "tcp"
           }
         ],
         "environment": [
           {
             "name": "TED_LOGGING_LEVEL",
             "value": "info"
           }
         ],
         "logConfiguration": {
           "logDriver": "awslogs",
           "options": {
             "awslogs-group": "/ecs/tokentropydrift",
             "awslogs-region": "us-west-2",
             "awslogs-stream-prefix": "ecs"
           }
         }
       }
     ]
   }
   ```

2. **Deploy to ECS**:
   ```bash
   aws ecs create-service \
     --cluster your-cluster \
     --service-name tokentropydrift \
     --task-definition tokentropydrift:1 \
     --desired-count 2 \
     --launch-type FARGATE \
     --network-configuration "awsvpcConfiguration={subnets=[subnet-xxxxxxxxx],securityGroups=[sg-xxxxxxxxx],assignPublicIp=ENABLED}"
   ```

### Google Cloud Platform

#### Compute Engine Deployment

1. **Create instance**:
   ```bash
   gcloud compute instances create tokentropydrift \
     --zone=us-central1-a \
     --machine-type=e2-medium \
     --image-family=debian-11 \
     --image-project=debian-cloud \
     --metadata-from-file startup-script=startup.sh
   ```

2. **Create startup.sh**:
   ```bash
   #!/bin/bash
   apt-get update
   apt-get install -y git golang-go python3 python3-pip
   
   # Clone and build
   git clone https://github.com/RevBooyah/TokEntropyDrift.git
   cd TokEntropyDrift
   go build -o ted cmd/ted/main.go
   
   # Install dependencies
   pip3 install tiktoken transformers sentencepiece
   
   # Start application
   nohup ./ted serve --host=0.0.0.0 --port=8080 > app.log 2>&1 &
   ```

#### Cloud Run Deployment

1. **Create Dockerfile for Cloud Run**:
   ```dockerfile
   FROM golang:1.22-alpine AS builder
   WORKDIR /app
   COPY . .
   RUN go build -o ted cmd/ted/main.go
   
   FROM alpine:latest
   RUN apk --no-cache add ca-certificates python3 py3-pip
   WORKDIR /app
   COPY --from=builder /app/ted .
   COPY --from=builder /app/ted.config.yaml .
   RUN pip3 install tiktoken transformers sentencepiece
   
   EXPOSE 8080
   ENTRYPOINT ["./ted", "serve", "--host=0.0.0.0", "--port=8080"]
   ```

2. **Deploy to Cloud Run**:
   ```bash
   gcloud builds submit --tag gcr.io/PROJECT_ID/tokentropydrift
   gcloud run deploy tokentropydrift \
     --image gcr.io/PROJECT_ID/tokentropydrift \
     --platform managed \
     --region us-central1 \
     --allow-unauthenticated
   ```

### Azure Deployment

#### Azure Container Instances

1. **Create container group**:
   ```bash
   az container create \
     --resource-group myResourceGroup \
     --name tokentropydrift \
     --image your-registry.azurecr.io/tokentropydrift:latest \
     --dns-name-label tokentropydrift \
     --ports 8080 \
     --environment-variables TED_LOGGING_LEVEL=info
   ```

#### Azure Kubernetes Service

1. **Create deployment.yaml**:
   ```yaml
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: tokentropydrift
   spec:
     replicas: 3
     selector:
       matchLabels:
         app: tokentropydrift
     template:
       metadata:
         labels:
           app: tokentropydrift
       spec:
         containers:
         - name: ted
           image: your-registry.azurecr.io/tokentropydrift:latest
           ports:
           - containerPort: 8080
           env:
           - name: TED_LOGGING_LEVEL
             value: "info"
           resources:
             requests:
               memory: "512Mi"
               cpu: "250m"
             limits:
               memory: "1Gi"
               cpu: "500m"
   ---
   apiVersion: v1
   kind: Service
   metadata:
     name: tokentropydrift-service
   spec:
     selector:
       app: tokentropydrift
     ports:
     - protocol: TCP
       port: 80
       targetPort: 8080
     type: LoadBalancer
   ```

2. **Deploy to AKS**:
   ```bash
   kubectl apply -f deployment.yaml
   ```

## Production Deployment

### High Availability Setup

1. **Load Balancer Configuration**:
   ```nginx
   # nginx.conf
   upstream tokentropydrift {
       server 10.0.1.10:8080;
       server 10.0.1.11:8080;
       server 10.0.1.12:8080;
   }
   
   server {
       listen 80;
       server_name your-domain.com;
       
       location / {
           proxy_pass http://tokentropydrift;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
           proxy_set_header X-Forwarded-Proto $scheme;
       }
   }
   ```

2. **Database Configuration** (for persistent storage):
   ```yaml
   # PostgreSQL configuration
   database:
     host: your-db-host
     port: 5432
     name: tokentropydrift
     user: ted_user
     password: your-secure-password
     ssl_mode: require
   ```

3. **Redis Configuration** (for caching):
   ```yaml
   cache:
     enabled: true
     type: redis
     redis:
       host: your-redis-host
       port: 6379
       password: your-redis-password
       db: 0
   ```

### Scaling Configuration

1. **Horizontal Scaling**:
   ```yaml
   # Kubernetes HPA
   apiVersion: autoscaling/v2
   kind: HorizontalPodAutoscaler
   metadata:
     name: tokentropydrift-hpa
   spec:
     scaleTargetRef:
       apiVersion: apps/v1
       kind: Deployment
       name: tokentropydrift
     minReplicas: 3
     maxReplicas: 10
     metrics:
     - type: Resource
       resource:
         name: cpu
         target:
           type: Utilization
           averageUtilization: 70
     - type: Resource
       resource:
         name: memory
         target:
           type: Utilization
           averageUtilization: 80
   ```

2. **Resource Limits**:
   ```yaml
   resources:
     requests:
       memory: "512Mi"
       cpu: "250m"
     limits:
       memory: "2Gi"
       cpu: "1000m"
   ```

## Monitoring and Logging

### Application Monitoring

1. **Prometheus Configuration**:
   ```yaml
   # prometheus.yml
   global:
     scrape_interval: 15s
   
   scrape_configs:
     - job_name: 'tokentropydrift'
       static_configs:
         - targets: ['localhost:8080']
       metrics_path: '/metrics'
   ```

2. **Grafana Dashboard**:
   ```json
   {
     "dashboard": {
       "title": "TokEntropyDrift Metrics",
       "panels": [
         {
           "title": "Request Rate",
           "type": "graph",
           "targets": [
             {
               "expr": "rate(http_requests_total[5m])",
               "legendFormat": "{{method}} {{endpoint}}"
             }
           ]
         },
         {
           "title": "Response Time",
           "type": "graph",
           "targets": [
             {
               "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
               "legendFormat": "95th percentile"
             }
           ]
         }
       ]
     }
   }
   ```

### Logging Configuration

1. **Structured Logging**:
   ```yaml
   logging:
     level: info
     format: json
     file: /var/log/tokentropydrift/app.log
     max_size: 100MB
     max_age: 30d
     max_backups: 10
     compress: true
   ```

2. **Log Aggregation** (ELK Stack):
   ```yaml
   # Filebeat configuration
   filebeat.inputs:
   - type: log
     paths:
       - /var/log/tokentropydrift/*.log
     json.keys_under_root: true
     json.add_error_key: true
   
   output.elasticsearch:
     hosts: ["your-elasticsearch-host:9200"]
     index: "tokentropydrift-%{+yyyy.MM.dd}"
   ```

### Health Checks

1. **Application Health Endpoint**:
   ```go
   func healthHandler(w http.ResponseWriter, r *http.Request) {
       response := map[string]interface{}{
           "status": "healthy",
           "timestamp": time.Now().UTC(),
           "version": "1.0.0",
           "uptime": time.Since(startTime).String(),
       }
       
       w.Header().Set("Content-Type", "application/json")
       json.NewEncoder(w).Encode(response)
   }
   ```

2. **Kubernetes Health Checks**:
   ```yaml
   livenessProbe:
     httpGet:
       path: /health
       port: 8080
     initialDelaySeconds: 30
     periodSeconds: 10
   
   readinessProbe:
     httpGet:
       path: /ready
       port: 8080
     initialDelaySeconds: 5
     periodSeconds: 5
   ```

## Security Considerations

### Network Security

1. **Firewall Configuration**:
   ```bash
   # UFW configuration
   ufw allow 22/tcp
   ufw allow 8080/tcp
   ufw enable
   ```

2. **SSL/TLS Configuration**:
   ```nginx
   server {
       listen 443 ssl http2;
       server_name your-domain.com;
       
       ssl_certificate /path/to/cert.pem;
       ssl_certificate_key /path/to/key.pem;
       
       ssl_protocols TLSv1.2 TLSv1.3;
       ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;
       ssl_prefer_server_ciphers off;
       
       location / {
           proxy_pass http://tokentropydrift;
       }
   }
   ```

### Application Security

1. **Environment Variables**:
   ```bash
   # .env file
   TED_DB_PASSWORD=your-secure-password
   TED_REDIS_PASSWORD=your-redis-password
   TED_API_KEY=your-api-key
   ```

2. **Secrets Management**:
   ```yaml
   # Kubernetes secrets
   apiVersion: v1
   kind: Secret
   metadata:
     name: tokentropydrift-secrets
   type: Opaque
   data:
     db-password: <base64-encoded-password>
     redis-password: <base64-encoded-password>
     api-key: <base64-encoded-api-key>
   ```

3. **Input Validation**:
   ```go
   func validateInput(text string) error {
       if len(text) > MaxInputSize {
           return errors.New("input too large")
       }
       
       if !isValidText(text) {
           return errors.New("invalid input format")
       }
       
       return nil
   }
   ```

### Access Control

1. **API Authentication**:
   ```go
   func authMiddleware(next http.Handler) http.Handler {
       return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
           apiKey := r.Header.Get("X-API-Key")
           if apiKey != os.Getenv("TED_API_KEY") {
               http.Error(w, "Unauthorized", http.StatusUnauthorized)
               return
           }
           next.ServeHTTP(w, r)
       })
   }
   ```

2. **Rate Limiting**:
   ```go
   func rateLimitMiddleware(next http.Handler) http.Handler {
       limiter := rate.NewLimiter(rate.Every(time.Second), 10)
       
       return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
           if !limiter.Allow() {
               http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
               return
           }
           next.ServeHTTP(w, r)
       })
   }
   ```

## Troubleshooting

### Common Issues

1. **High Memory Usage**:
   ```bash
   # Check memory usage
   ps aux | grep ted
   
   # Enable streaming for large files
   ./ted advanced streaming large_file.txt
   ```

2. **Slow Performance**:
   ```bash
   # Enable caching
   ./ted advanced cache input.txt
   
   # Enable parallel processing
   ./ted advanced parallel input.txt
   ```

3. **Tokenizer Errors**:
   ```bash
   # Check tokenizer configuration
   cat ted.config.yaml | grep -A 10 tokenizers
   
   # Test individual tokenizer
   ./ted analyze input.txt --tokenizers=gpt2
   ```

### Debug Mode

1. **Enable Debug Logging**:
   ```bash
   export TED_LOGGING_LEVEL=debug
   ./ted serve --port=8080
   ```

2. **Profile Performance**:
   ```bash
   go build -o ted cmd/ted/main.go
   ./ted analyze input.txt &
   go tool pprof http://localhost:6060/debug/pprof/profile
   ```

### Recovery Procedures

1. **Application Restart**:
   ```bash
   # Graceful shutdown
   pkill -TERM ted
   
   # Force kill if needed
   pkill -KILL ted
   
   # Restart
   nohup ./ted serve --port=8080 > app.log 2>&1 &
   ```

2. **Data Recovery**:
   ```bash
   # Backup output directory
   tar -czf output_backup_$(date +%Y%m%d).tar.gz output/
   
   # Restore from backup
   tar -xzf output_backup_20231201.tar.gz
   ```

### Support and Resources

- **Documentation**: See `/docs/` for detailed guides
- **Issues**: Report problems on GitHub
- **Community**: Join discussions for help
- **Monitoring**: Use provided monitoring tools

## Conclusion

This deployment guide covers the essential aspects of deploying TokEntropyDrift in various environments. Choose the deployment method that best fits your requirements and infrastructure.

For additional help or custom deployment scenarios, please refer to the project documentation or community resources. 