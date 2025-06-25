#!/bin/bash

# Start only the API server

set -e

echo "🚀 Starting Domain Check API Server..."

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

# Install dependencies
echo -e "${BLUE}📦 Installing Go dependencies...${NC}"
go mod tidy

# Kill existing process on port 8080
echo -e "${BLUE}🧹 Cleaning up existing processes on port 8080...${NC}"
lsof -ti:8080 | xargs kill -9 2>/dev/null || true

# Start API server
echo -e "${BLUE}🚀 Starting API server on :8080...${NC}"
go run cmd/server/main.go
