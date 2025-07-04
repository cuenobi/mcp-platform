#!/bin/bash

# MCP Platform Docker Setup Script

set -e

echo "🐳 MCP Platform Docker Setup"
echo "============================"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Function to setup environment file
setup_env() {
    echo "📝 Setting up environment file..."
    
    if [ ! -f ".env" ]; then
        cp docker.env .env
        echo "⚠️  Please edit .env file with your Jira credentials:"
        echo "   - JIRA_BASE_URL"
        echo "   - JIRA_EMAIL"
        echo "   - JIRA_API_TOKEN"
        echo "   - JIRA_PROJECT_KEY"
        echo "   - JIRA_BOARD_ID"
        echo ""
        read -p "Press Enter after editing the .env file..."
    else
        echo "✅ .env file already exists"
    fi
}

# Function to pull Ollama model
setup_ollama() {
    echo "🤖 Setting up Ollama model..."
    
    # Start Ollama service if not running
    docker-compose up -d ollama
    
    # Wait for Ollama to be ready
    echo "⏳ Waiting for Ollama service to be ready..."
    until docker-compose exec ollama ollama list &> /dev/null; do
        echo "   Waiting for Ollama..."
        sleep 5
    done
    
    # Pull the llama3 model
    echo "📥 Pulling llama3 model (this may take a while)..."
    docker-compose exec ollama ollama pull llama3
    
    echo "✅ Ollama setup complete!"
}

# Function to build all services
build_services() {
    echo "🏗️  Building all services..."
    docker-compose build
    echo "✅ All services built successfully!"
}

# Function to start all services
start_services() {
    echo "🚀 Starting all services..."
    docker-compose up -d
    echo "✅ All services started!"
    
    echo ""
    echo "📋 Service Status:"
    echo "   - API Gateway: http://localhost:8080"
    echo "   - MCP Server Jira: gRPC on localhost:50051"
    echo "   - Ollama: http://localhost:11434"
    echo "   - MCP Host: Running as client service"
}

# Function to show logs
show_logs() {
    echo "📊 Showing service logs..."
    docker-compose logs -f
}

# Function to stop all services
stop_services() {
    echo "🛑 Stopping all services..."
    docker-compose down
    echo "✅ All services stopped!"
}

# Function to clean up
cleanup() {
    echo "🧹 Cleaning up..."
    docker-compose down -v
    docker system prune -f
    echo "✅ Cleanup complete!"
}

# Function to test the setup
test_setup() {
    echo "🧪 Testing the setup..."
    
    # Test if services are running
    if ! docker-compose ps | grep -q "Up"; then
        echo "❌ Services are not running. Please start them first."
        exit 1
    fi
    
    # Test Jira card creation
    echo "Creating a test Jira card..."
    docker-compose exec mcphost ./mcphost jira create-card --project AIT --prompt "Test card from Docker setup"
    
    echo "✅ Test complete!"
}

# Main menu
case "${1:-}" in
    "setup")
        setup_env
        build_services
        setup_ollama
        start_services
        ;;
    "build")
        build_services
        ;;
    "start")
        start_services
        ;;
    "stop")
        stop_services
        ;;
    "logs")
        show_logs
        ;;
    "test")
        test_setup
        ;;
    "cleanup")
        cleanup
        ;;
    "ollama")
        setup_ollama
        ;;
    *)
        echo "Usage: $0 {setup|build|start|stop|logs|test|cleanup|ollama}"
        echo ""
        echo "Commands:"
        echo "  setup   - Complete setup (env, build, ollama, start)"
        echo "  build   - Build all Docker images"
        echo "  start   - Start all services"
        echo "  stop    - Stop all services"
        echo "  logs    - Show service logs"
        echo "  test    - Test the setup by creating a Jira card"
        echo "  cleanup - Stop services and clean up"
        echo "  ollama  - Setup Ollama model only"
        exit 1
        ;;
esac 