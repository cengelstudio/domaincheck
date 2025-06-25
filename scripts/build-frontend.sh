#!/bin/bash

# Build frontend for production

set -e

echo "🎨 Building Domain Check Frontend for Production..."

# Colors
RED='\033[0;31m'
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

echo -e "${GREEN}✅ Node.js and npm are available${NC}"

# Go to frontend directory
cd frontend

# Install dependencies
echo -e "${BLUE}📦 Installing frontend dependencies...${NC}"
npm install
echo -e "${GREEN}✅ Frontend dependencies installed${NC}"

# Clean previous build
echo -e "${BLUE}🧹 Cleaning previous build...${NC}"
rm -rf dist/

# Build frontend
echo -e "${BLUE}🔨 Building frontend...${NC}"
npm run build

# Check if build was successful
if [ -d "dist" ]; then
    echo -e "${GREEN}✅ Frontend built successfully${NC}"
    echo -e "${GREEN}📁 Build output: frontend/dist/${NC}"

    # Show build info
    echo -e "${YELLOW}📊 Build statistics:${NC}"
    du -sh dist/ 2>/dev/null || echo "   Size calculation unavailable"

    echo -e "${YELLOW}📋 Build contents:${NC}"
    ls -la dist/ | head -10

    echo ""
    echo -e "${GREEN}🎉 Frontend build completed successfully!${NC}"
    echo -e "${BLUE}You can now run the API server to serve the built frontend${NC}"
else
    echo -e "${RED}❌ Frontend build failed${NC}"
    exit 1
fi
