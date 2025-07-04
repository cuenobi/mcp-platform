# MCP Platform Docker Setup 🐳

This guide will help you run the entire MCP Platform using Docker containers.

## Architecture

The Docker setup includes:
- **Ollama** - Local AI model server (llama3)
- **MCP Server Jira** - gRPC server for Jira operations
- **API Gateway** - HTTP gateway service
- **MCP Host** - Client service for interacting with Jira

## Prerequisites

1. **Docker** and **Docker Compose** installed
2. **Jira credentials** (API token, base URL, email)
3. At least **8GB RAM** (for Ollama llama3 model)
4. **10GB+ free disk space** (for Docker images and Ollama model)

## Quick Start

### 1. Setup Environment

```bash
# Make the setup script executable
chmod +x docker-setup.sh

# Run complete setup
./docker-setup.sh setup
```

This will:
- Create a `.env` file from the template
- Prompt you to edit Jira credentials
- Build all Docker images
- Download and setup Ollama llama3 model
- Start all services

### 2. Configure Jira Credentials

Edit the `.env` file with your Jira details:

```env
# Jira Configuration
JIRA_BASE_URL=https://your-domain.atlassian.net
JIRA_EMAIL=your-email@example.com
JIRA_API_TOKEN=your-api-token
JIRA_PROJECT_KEY=AIT
JIRA_BOARD_ID=your-board-id
```

### 3. Test the Setup

```bash
# Test by creating a Jira card
./docker-setup.sh test
```

## Manual Commands

### Build Services
```bash
docker-compose build
```

### Start Services
```bash
docker-compose up -d
```

### View Logs
```bash
docker-compose logs -f
```

### Stop Services
```bash
docker-compose down
```

### Clean Up Everything
```bash
docker-compose down -v
docker system prune -f
```

## Service Access

| Service | URL/Port | Description |
|---------|----------|-------------|
| API Gateway | `http://localhost:8080` | HTTP gateway |
| MCP Server Jira | `localhost:50051` | gRPC server |
| Ollama | `http://localhost:11434` | AI model server |
| MCP Host | Container only | Client service |

## Using the Services

### Create a Jira Card

```bash
# Enter the MCP Host container
docker-compose exec mcphost sh

# Create a card in Thai
./mcphost jira create-card --project AIT --prompt "สร้างการ์ดสำหรับปรับปรุงระบบสมัครสมาชิก"

# Create a card in English
./mcphost jira create-card --project AIT --prompt "Improve user registration system"
```

### Sync Jira Issues

```bash
docker-compose exec mcphost ./mcphost jira sync --project AIT
```

## Troubleshooting

### Services Not Starting
```bash
# Check service status
docker-compose ps

# View logs for specific service
docker-compose logs mcp-server-jira
docker-compose logs ollama
```

### Ollama Model Issues
```bash
# Restart Ollama and pull model again
docker-compose restart ollama
./docker-setup.sh ollama
```

### Port Conflicts
If ports are already in use, edit `docker-compose.yml`:
```yaml
ports:
  - "8081:8080"  # Change 8080 to 8081
  - "50052:50051"  # Change 50051 to 50052
```

### Memory Issues
Ollama requires significant RAM:
- Minimum: 8GB system RAM
- Recommended: 16GB+ system RAM
- For lower RAM systems, consider using a smaller model

### Reset Everything
```bash
./docker-setup.sh cleanup
./docker-setup.sh setup
```

## Development

### Rebuilding After Code Changes
```bash
# Rebuild specific service
docker-compose build mcp-server-jira

# Rebuild and restart
docker-compose up -d --build mcp-server-jira
```

### Viewing Service Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f mcp-server-jira
```

### Debugging Inside Containers
```bash
# Enter MCP Host container
docker-compose exec mcphost sh

# Enter MCP Server Jira container
docker-compose exec mcp-server-jira sh
```

## Architecture Details

### Network Configuration
- All services run on `mcp-network` Docker network
- Services communicate using Docker service names
- External ports are mapped for development access

### Environment Variables
- `MCP_SERVER_JIRA_ADDR=mcp-server-jira:50051` - gRPC address
- `OLLAMA_BASE_URL=http://ollama:11434` - Ollama API URL
- Jira credentials passed from `.env` file

### Data Persistence
- Ollama models stored in `ollama_data` volume
- No persistent data for other services (stateless)

## Performance Tips

1. **Ollama Model**: First run downloads ~4GB model
2. **Build Cache**: Use `docker-compose build --no-cache` if needed
3. **Resource Limits**: Adjust in `docker-compose.yml` if needed
4. **Network**: All services on same Docker network for fast communication

## Next Steps

1. **Add monitoring** with Prometheus/Grafana
2. **Add reverse proxy** with Nginx
3. **Add TLS certificates** for production
4. **Add backup scripts** for data
5. **Add CI/CD pipeline** for automated deployment

## Support

For issues or questions:
1. Check the logs: `docker-compose logs -f`
2. Verify .env configuration
3. Ensure Docker has sufficient resources
4. Check network connectivity between services 