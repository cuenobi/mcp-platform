# MCP Platform

A comprehensive **Model Context Protocol (MCP) Platform** that enables seamless integration between LLM applications and various external services through a microservice architecture.

## 🎯 Project Objectives

- **Modular Architecture**: Build a scalable microservice platform for MCP server integration
- **LLM Integration**: Provide seamless connectivity between Large Language Models and external tools/services
- **Tool Orchestration**: Centralized management of tool calls and service interactions
- **Extensible Design**: Easy addition of new MCP servers and tools
- **Production Ready**: Include authentication, rate limiting, and session management
- **Multi-Service Support**: Support for multiple MCP servers (Jira, Calculator, etc.)

## 🏗️ Architecture

```text
┌──────────────────────────────┐
│          Web UI              │  <--- Frontend app (React, Vue, etc.)
└─────────────┬────────────────┘
              │ HTTP/REST/WS
┌─────────────▼────────────────┐
│       API Gateway / Backend  │  <--- Backend service (Golang, Node.js)
│  - รับ request จาก UI         │
│  - ติดต่อกับ mcphost process    │
│  - จัดการ session, auth,      │
│    rate limit ฯลฯ            │
└─────────────┬────────────────┘
              │ stdin/stdout or TCP or gRPC
┌─────────────▼────────────────┐
│          MCP Host            │  <--- mcphost CLI process (Microservice)
│  - เชื่อม LLM กับ MCP Server    │
│  - จัดการการเรียก tool         │
└─────────────┬────────────────┘
              │ gRPC
┌─────────────▼───────────────────────────┐
│       MCP Server(s)                     │  <--- Microservice MCP Server
│  - ให้บริการ tool ต่าง ๆ                   │
│  - อาจรันหลายตัว (calculator, Jira, etc.) │
└─────────────────────────────────────────┘
```

## 📦 Components

### 1. API Gateway (`api-gateway/`)
- **Purpose**: RESTful API gateway and backend service
- **Technology**: Go with Cobra CLI framework
- **Responsibilities**:
  - Handle HTTP/REST/WebSocket requests from frontend
  - Manage user sessions and authentication
  - Rate limiting and request throttling
  - Route requests to appropriate MCP Host processes

### 2. MCP Host (`mcphost/`)
- **Purpose**: Bridge between LLM and MCP Servers
- **Technology**: Go with gRPC communication
- **Responsibilities**:
  - Manage tool calls and orchestrate service interactions
  - Connect LLM requests to appropriate MCP servers
  - Handle protocol translation and message routing
  - Provide CLI interface for direct service interaction

### 3. MCP Server - Jira (`mcp-server-jira/`)
- **Purpose**: Jira integration MCP server
- **Technology**: Go with REST API integration
- **Features**:
  - Create Jira issues from prompts
  - Sync Jira project data
  - LLM-powered issue generation
  - Secure authentication with Jira API

### 4. Shared Protocol (`shared/proto/`)
- **Purpose**: Protocol definitions for inter-service communication
- **Technology**: Protocol Buffers (gRPC)
- **Services**:
  - `JiraService`: Issue creation, synchronization
  - Extensible for additional services

## ✨ Features

- **🔌 Modular MCP Server Architecture**: Easy integration of new tools and services
- **🤖 LLM Integration**: Seamless connectivity with Large Language Models
- **📊 Jira Integration**: Create and manage Jira issues through natural language
- **🔐 Secure Communication**: gRPC-based inter-service communication
- **🛠️ CLI Tools**: Direct command-line access to all services
- **📈 Scalable Design**: Microservice architecture for horizontal scaling
- **🔄 Real-time Updates**: Support for WebSocket connections
- **📋 Protocol Standardization**: Consistent API across all MCP servers

## 🚀 Prerequisites

- **Go**: Version 1.23.9 or later
- **Protocol Buffers**: For gRPC code generation
- **Make**: For build automation
- **Jira Account**: For Jira integration features

## 🔧 Installation & Setup

### 1. Clone the Repository
```bash
git clone https://github.com/cuenobi/mcp-platform.git
cd mcp-platform
```

### 2. Install Dependencies
```bash
# Install Go modules for all components
go mod download

# Install dependencies for individual components
cd api-gateway && go mod download && cd ..
cd mcphost && go mod download && cd ..
cd mcp-server-jira && go mod download && cd ..
cd shared/proto/gen && go mod download && cd ../../..
```

### 3. Generate Protocol Buffers
```bash
cd shared/proto
make generate
```

### 4. Environment Configuration
Create `.env` file for Jira integration:
```bash
# Jira Configuration
JIRA_BASE_URL=https://your-domain.atlassian.net
JIRA_EMAIL=your-email@example.com
JIRA_API_TOKEN=your-jira-api-token
```

### 5. Build Services
```bash
# Build API Gateway
cd api-gateway
go build -o api-gateway .

# Build MCP Host
cd ../mcphost
go build -o mcphost .

# Build MCP Server Jira
cd ../mcp-server-jira
go build -o mcp-server-jira .
```

## 📖 Usage

### API Gateway
```bash
cd api-gateway
./api-gateway apiGateway
```

### MCP Host - Jira Operations
```bash
cd mcphost

# Sync Jira issues
./mcphost jira sync --project YOUR_PROJECT_KEY

# Create Jira issue from prompt
./mcphost jira create-card --project YOUR_PROJECT_KEY --prompt "Create a bug report for login issue"
```

### MCP Server - Direct Usage
```bash
cd mcp-server-jira
./mcp-server-jira
```

## 🏗️ Development

### Adding a New MCP Server

1. **Create Service Directory**
```bash
mkdir mcp-server-newservice
cd mcp-server-newservice
```

2. **Initialize Go Module**
```bash
go mod init github.com/cuenobi/mcp-platform/mcp-server-newservice
```

3. **Define Protocol**
```proto
// shared/proto/protos/newservice.proto
syntax = "proto3";
package newservice;

service NewService {
  rpc ProcessRequest(ProcessRequest) returns (ProcessResponse);
}
```

4. **Generate Code**
```bash
cd shared/proto
make generate
```

5. **Implement Service**
```go
// Implement your MCP server logic
```

### Protocol Buffer Development
```bash
cd shared/proto

# Generate Go code
make generate

# Clean generated files
make clean
```

## 🔒 Security

- **Environment Variables**: Store sensitive credentials in environment variables
- **API Token Authentication**: Use secure API tokens for external service integration
- **gRPC Security**: Implement TLS for production gRPC communication
- **Rate Limiting**: Built-in rate limiting in API Gateway
- **Input Validation**: Comprehensive input validation across all services

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go best practices and conventions
- Write comprehensive tests for new features
- Update documentation for API changes
- Use meaningful commit messages
- Ensure all services build successfully

## 📄 License

This project is licensed under the MIT License - see the individual LICENSE files in each component directory for details.

## 🆘 Support

For questions and support:
- Create an issue in the GitHub repository
- Check existing documentation in each component
- Review the architecture diagram for system understanding

---

**Built with ❤️ using Go, gRPC, and Modern Microservice Architecture**