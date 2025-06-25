#!/bin/bash

# Start only the frontend development server

set -e

echo "🎨 Starting Domain Check Frontend Development Server..."

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check Node.js and npm
if ! command_exists node; then
    echo -e "${RED}❌ Node.js is not installed${NC}"
    exit 1
fi

if ! command_exists npm; then
    echo -e "${RED}❌ npm is not installed${NC}"
    exit 1
fi

# Install frontend dependencies
echo -e "${BLUE}📦 Installing frontend dependencies...${NC}"
cd frontend

if [ ! -d "node_modules" ]; then
    npm install
    echo -e "${GREEN}✅ Frontend dependencies installed${NC}"
else
    echo -e "${YELLOW}   node_modules exists, skipping npm install${NC}"
fi

# Kill existing process on port 3000
echo -e "${BLUE}🧹 Cleaning up existing processes on port 3000...${NC}"
lsof -ti:3000 | xargs kill -9 2>/dev/null || true

# Start frontend development server
echo -e "${BLUE}🎨 Starting frontend development server on :3000...${NC}"
echo -e "${GREEN}Frontend will be available at: http://localhost:3000${NC}"
npm run serve
