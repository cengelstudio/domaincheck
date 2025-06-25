#!/bin/bash

# Stop all Domain Check services

echo "🛑 Stopping Domain Check Services..."

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Function to kill processes on a port
kill_port() {
    local port=$1
    local service_name=$2

    echo -e "${BLUE}🔍 Checking for processes on port $port ($service_name)...${NC}"

    local pids=$(lsof -ti:$port 2>/dev/null || true)

    if [ -n "$pids" ]; then
        echo -e "${YELLOW}   Found processes: $pids${NC}"
        echo -e "${BLUE}   Stopping $service_name...${NC}"
        echo $pids | xargs kill -TERM 2>/dev/null || true

        # Wait a bit for graceful shutdown
        sleep 2

        # Force kill if still running
        local remaining=$(lsof -ti:$port 2>/dev/null || true)
        if [ -n "$remaining" ]; then
            echo -e "${YELLOW}   Force killing remaining processes...${NC}"
            echo $remaining | xargs kill -9 2>/dev/null || true
        fi

        echo -e "${GREEN}   ✅ $service_name stopped${NC}"
    else
        echo -e "${GREEN}   ✅ No $service_name processes found${NC}"
    fi
}

# Stop services by port
kill_port 8080 "API Server"
kill_port 3000 "Frontend Dev Server"

# Stop processes by PID files if they exist
if [ -f ".api.pid" ]; then
    echo -e "${BLUE}🔍 Found API PID file...${NC}"
    api_pid=$(cat .api.pid)
    if kill -0 $api_pid 2>/dev/null; then
        echo -e "${BLUE}   Stopping API process (PID: $api_pid)...${NC}"
        kill -TERM $api_pid 2>/dev/null || true
        sleep 2
        kill -9 $api_pid 2>/dev/null || true
        echo -e "${GREEN}   ✅ API process stopped${NC}"
    fi
    rm -f .api.pid
fi

if [ -f ".frontend.pid" ]; then
    echo -e "${BLUE}🔍 Found frontend PID file...${NC}"
    frontend_pid=$(cat .frontend.pid)
    if kill -0 $frontend_pid 2>/dev/null; then
        echo -e "${BLUE}   Stopping frontend process (PID: $frontend_pid)...${NC}"
        kill -TERM $frontend_pid 2>/dev/null || true
        sleep 2
        kill -9 $frontend_pid 2>/dev/null || true
        echo -e "${GREEN}   ✅ Frontend process stopped${NC}"
    fi
    rm -f .frontend.pid
fi

# Clean up any remaining Go processes related to domaincheck
echo -e "${BLUE}🧹 Cleaning up any remaining domaincheck processes...${NC}"
pkill -f "domaincheck" 2>/dev/null || true
pkill -f "cmd/server/main.go" 2>/dev/null || true

# Clean up build artifacts if requested
if [ "$1" = "--clean" ]; then
    echo -e "${BLUE}🧹 Cleaning build artifacts...${NC}"
    rm -f domaincheck
    rm -f domaincheck.exe
    rm -rf dist/
    rm -rf frontend/dist/
    echo -e "${GREEN}   ✅ Build artifacts cleaned${NC}"
fi

echo -e "${GREEN}"
echo "================================"
echo "✅ All services stopped!"
echo "================================"
echo -e "${NC}"

# Show what was stopped
echo -e "${YELLOW}Stopped services:${NC}"
echo "   • API Server (port 8080)"
echo "   • Frontend Dev Server (port 3000)"
echo "   • Background processes"
echo ""

if [ "$1" = "--clean" ]; then
    echo -e "${YELLOW}Cleaned:${NC}"
    echo "   • Build artifacts"
    echo "   • Distribution files"
    echo ""
fi

echo -e "${BLUE}To start services again:${NC}"
echo "   • Development: ./scripts/dev.sh"
echo "   • API only: ./scripts/start-api.sh"
echo "   • Frontend only: ./scripts/start-frontend.sh"
