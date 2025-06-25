#!/bin/bash

# Development environment startup script for Domain Check API

set -e  # Exit on any error

echo "🚀 Starting Domain Check Development Environment..."

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to wait for server to be ready
wait_for_server() {
    echo -e "${YELLOW}⏳ Waiting for API server to be ready...${NC}"
    local max_attempts=30
    local attempt=1

    while [ $attempt -le $max_attempts ]; do
        if curl -s http://localhost:8080/api/health >/dev/null 2>&1; then
            echo -e "${GREEN}✅ API server is ready!${NC}"
            return 0
        fi
        echo -e "${YELLOW}   Attempt $attempt/$max_attempts - waiting...${NC}"
        sleep 1
        ((attempt++))
    done

    echo -e "${RED}❌ API server failed to start after $max_attempts seconds${NC}"
    return 1
}

# Check dependencies
echo -e "${BLUE}🔍 Checking dependencies...${NC}"

if ! command_exists go; then
    echo -e "${RED}❌ Go is not installed${NC}"
    exit 1
fi

if ! command_exists node; then
    echo -e "${RED}❌ Node.js is not installed${NC}"
    exit 1
fi

if ! command_exists npm; then
    echo -e "${RED}❌ npm is not installed${NC}"
    exit 1
fi

echo -e "${GREEN}✅ All dependencies are available${NC}"

# Install Go dependencies
echo -e "${BLUE}📦 Installing Go dependencies...${NC}"
go mod tidy

# Install frontend dependencies
echo -e "${BLUE}📦 Installing frontend dependencies...${NC}"
cd frontend
if [ ! -d "node_modules" ]; then
    npm install
else
    echo -e "${YELLOW}   node_modules exists, skipping npm install${NC}"
fi
cd ..

# Kill any existing processes on ports 8080 and 3000
echo -e "${BLUE}🧹 Cleaning up existing processes...${NC}"
lsof -ti:8080 | xargs kill -9 2>/dev/null || true
lsof -ti:3000 | xargs kill -9 2>/dev/null || true

# Start API server in background
echo -e "${BLUE}🚀 Starting API server...${NC}"
go run cmd/server/main.go &
API_PID=$!

# Wait for API server to be ready
if wait_for_server; then
    echo -e "${GREEN}✅ API server started successfully (PID: $API_PID)${NC}"
else
    echo -e "${RED}❌ Failed to start API server${NC}"
    kill $API_PID 2>/dev/null || true
    exit 1
fi

# Start frontend development server
echo -e "${BLUE}🎨 Starting frontend development server...${NC}"
cd frontend
npm run serve &
FRONTEND_PID=$!
cd ..

# Display running services
echo -e "${GREEN}"
echo "=================================="
echo "🎉 Development Environment Ready!"
echo "=================================="
echo -e "${NC}"
echo -e "${BLUE}📡 API Server:${NC}      http://localhost:8080"
echo -e "${BLUE}🎨 Frontend:${NC}        http://localhost:3000"
echo -e "${BLUE}📋 API Health:${NC}      http://localhost:8080/api/health"
echo -e "${BLUE}📚 API Docs:${NC}        http://localhost:8080"
echo ""
echo -e "${YELLOW}💡 Tips:${NC}"
echo "   • Use 'scripts/test-api.sh' to test API endpoints"
echo "   • Use 'scripts/stop.sh' to stop all services"
echo "   • Use 'scripts/logs.sh' to view logs"
echo ""
echo -e "${YELLOW}🔧 Process IDs:${NC}"
echo "   • API Server PID: $API_PID"
echo "   • Frontend PID: $FRONTEND_PID"

# Save PIDs for cleanup
echo $API_PID > .api.pid
echo $FRONTEND_PID > .frontend.pid

# Wait for user interrupt
echo -e "${GREEN}Press Ctrl+C to stop all services...${NC}"
trap 'echo -e "\n${YELLOW}🛑 Stopping services...${NC}"; kill $API_PID $FRONTEND_PID 2>/dev/null || true; rm -f .api.pid .frontend.pid; echo -e "${GREEN}✅ All services stopped${NC}"; exit 0' INT

# Keep script running
wait
