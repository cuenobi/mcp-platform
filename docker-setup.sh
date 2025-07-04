#!/bin/bash

set -e

echo "🚀 MCP Platform Docker Setup"
echo "=============================="

if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

setup_env() {
    if [ ! -f .env ]; then
        echo "📝 Creating .env file template..."
        cat > .env << EOF
JIRA_BASE_URL=https://your-domain.atlassian.net
JIRA_EMAIL=your-email@example.com
JIRA_API_TOKEN=your-api-token
JIRA_PROJECT_KEY=YOUR_PROJECT
JIRA_BOARD_ID=1
OPENAI_API_KEY=your-openai-api-key
OLLAMA_BASE_URL=http://ollama:11434
EOF
        echo "✅ .env file created. Please update it with your actual values."
        echo "⚠️  You need to set proper Jira credentials before continuing."
    else
        echo "✅ .env file already exists."
    fi
}

setup_ollama() {
    echo "🤖 Setting up Ollama..."
    docker-compose up -d ollama
    
    echo "⏳ Waiting for Ollama to be ready..."
    sleep 10
    
    while ! docker-compose exec ollama ollama list &>/dev/null; do
        echo "   Still waiting for Ollama..."
        sleep 5
    done
    
    echo "📥 Pulling llama3 model..."
    docker-compose exec ollama ollama pull llama3
    
    echo "✅ Ollama setup complete!"
}

build_services() {
    echo "🔨 Building all services..."
    docker-compose build --no-cache
    echo "✅ Build complete!"
}

start_services() {
    echo "🚀 Starting all services..."
    docker-compose up -d
    
    echo "⏳ Waiting for services to be healthy..."
    sleep 30
    
    docker-compose ps
    
    echo "✅ All services started!"
    echo ""
    echo "🌐 Available endpoints:"
    echo "   - API Gateway: http://localhost:8080"
    echo "   - Ollama: http://localhost:11434"
    echo "   - MCP Server (gRPC): localhost:50051"
}

show_logs() {
    echo "📋 Recent logs:"
    docker-compose logs --tail=50
}

stop_services() {
    echo "🛑 Stopping all services..."
    docker-compose down
    echo "✅ All services stopped!"
}

cleanup() {
    echo "🧹 Cleaning up..."
    docker-compose down -v
    docker system prune -f
    echo "✅ Cleanup complete!"
}

test_setup() {
    echo "🧪 Testing setup..."
    
    if docker-compose ps | grep -q "Up"; then
        echo "✅ Services are running"
    else
        echo "❌ Services are not running properly"
        return 1
    fi
    
    echo "✅ Basic setup test passed!"
}

while true; do
    echo ""
    echo "Choose an option:"
    echo "1) Setup environment file"
    echo "2) Setup Ollama"
    echo "3) Build services"
    echo "4) Start services"
    echo "5) Show logs"
    echo "6) Stop services"
    echo "7) Cleanup (removes volumes)"
    echo "8) Test setup"
    echo "9) Exit"
    echo ""
    read -p "Enter your choice (1-9): " choice

    case $choice in
        1)
            setup_env
            ;;
        2)
            setup_ollama
            ;;
        3)
            build_services
            ;;
        4)
            start_services
            ;;
        5)
            show_logs
            ;;
        6)
            stop_services
            ;;
        7)
            cleanup
            ;;
        8)
            test_setup
            ;;
        9)
            echo "👋 Goodbye!"
            exit 0
            ;;
        *)
            echo "❌ Invalid option. Please choose 1-9."
            ;;
    esac
done 