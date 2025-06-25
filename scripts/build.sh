#!/bin/bash

# Production build script for Domain Check API

set -e

echo "🏗️  Building Domain Check API for Production..."

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
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

# Clean previous builds
echo -e "${BLUE}🧹 Cleaning previous builds...${NC}"
rm -rf dist/
rm -rf frontend/dist/
rm -f domaincheck
rm -f domaincheck.exe

# Install Go dependencies
echo -e "${BLUE}📦 Installing Go dependencies...${NC}"
go mod tidy

# Install frontend dependencies
echo -e "${BLUE}📦 Installing frontend dependencies...${NC}"
cd frontend
npm install
echo -e "${GREEN}✅ Frontend dependencies installed${NC}"

# Build frontend
echo -e "${BLUE}🎨 Building frontend...${NC}"
npm run build
echo -e "${GREEN}✅ Frontend built successfully${NC}"
cd ..

# Build Go binary
echo -e "${BLUE}🔨 Building Go binary...${NC}"

# Build for current platform
go build -ldflags "-s -w" -o domaincheck cmd/server/main.go
echo -e "${GREEN}✅ Go binary built successfully${NC}"

# Create distribution directory
echo -e "${BLUE}📦 Creating distribution package...${NC}"
mkdir -p dist

# Copy binary
cp domaincheck dist/

# Copy configuration
cp -r configs dist/

# Copy data
cp -r data dist/

# Copy frontend build
cp -r frontend/dist dist/frontend

# Copy scripts
cp -r scripts dist/

# Create README for distribution
cat > dist/README.md << EOF
# Domain Check API - Production Distribution

## Quick Start

1. **Configure**: Edit \`configs/config.yaml\` if needed
2. **Run**: \`./domaincheck\`
3. **Access**: http://localhost:8080

## Files

- \`domaincheck\` - Main executable
- \`configs/\` - Configuration files
- \`data/\` - Domain extensions data
- \`frontend/\` - Built frontend files
- \`scripts/\` - Utility scripts

## Configuration

Edit \`configs/config.yaml\` to customize:
- Server port
- CORS settings
- Domain check timeout
- Extensions file path

## Usage

### Start Server
\`\`\`bash
./domaincheck
\`\`\`

### API Endpoints
- \`GET /api/health\` - Health check
- \`POST /api/check-domain\` - Check domain
- \`GET /api/domains\` - Get history
- \`GET /api/v1/extensions\` - Get valid extensions

### Environment Variables
- \`CONFIG_PATH\` - Path to config file (default: ./configs/config.yaml)
- \`GIN_MODE\` - Gin mode (debug/release)

## Support

For support, visit the project repository.
EOF

# Create systemd service file (Linux)
cat > dist/domaincheck.service << EOF
[Unit]
Description=Domain Check API Server
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt/domaincheck
ExecStart=/opt/domaincheck/domaincheck
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# Create startup script
cat > dist/start.sh << 'EOF'
#!/bin/bash
cd "$(dirname "$0")"
./domaincheck
EOF
chmod +x dist/start.sh

# Cross-platform builds (optional)
echo -e "${BLUE}🔨 Building cross-platform binaries...${NC}"

# Linux AMD64
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o dist/domaincheck-linux-amd64 cmd/server/main.go
echo -e "${GREEN}   ✅ Linux AMD64 binary built${NC}"

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o dist/domaincheck-windows-amd64.exe cmd/server/main.go
echo -e "${GREEN}   ✅ Windows AMD64 binary built${NC}"

# macOS AMD64
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o dist/domaincheck-darwin-amd64 cmd/server/main.go
echo -e "${GREEN}   ✅ macOS AMD64 binary built${NC}"

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o dist/domaincheck-darwin-arm64 cmd/server/main.go
echo -e "${GREEN}   ✅ macOS ARM64 binary built${NC}"

# Create archive
echo -e "${BLUE}📦 Creating distribution archive...${NC}"
tar -czf domaincheck-dist.tar.gz -C dist .
echo -e "${GREEN}✅ Distribution archive created: domaincheck-dist.tar.gz${NC}"

# Display build summary
echo -e "${GREEN}"
echo "==============================="
echo "🎉 Build Completed Successfully!"
echo "==============================="
echo -e "${NC}"
echo -e "${BLUE}📁 Distribution folder:${NC} ./dist/"
echo -e "${BLUE}📦 Archive:${NC}             ./domaincheck-dist.tar.gz"
echo ""
echo -e "${YELLOW}📋 Built binaries:${NC}"
echo "   • domaincheck                    (current platform)"
echo "   • domaincheck-linux-amd64        (Linux 64-bit)"
echo "   • domaincheck-windows-amd64.exe  (Windows 64-bit)"
echo "   • domaincheck-darwin-amd64       (macOS Intel)"
echo "   • domaincheck-darwin-arm64       (macOS Apple Silicon)"
echo ""
echo -e "${YELLOW}🚀 To run:${NC}"
echo "   cd dist && ./domaincheck"
echo ""
echo -e "${YELLOW}📚 Documentation:${NC}"
echo "   See dist/README.md for deployment instructions"
