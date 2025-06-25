#!/bin/bash

# View logs and monitor Domain Check services

echo "📋 Domain Check Logs Viewer"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Function to show usage
show_usage() {
    echo -e "${BLUE}Usage:${NC}"
    echo "  ./scripts/logs.sh [option]"
    echo ""
    echo -e "${BLUE}Options:${NC}"
    echo "  api       - Show API server logs"
    echo "  frontend  - Show frontend logs"
    echo "  live      - Live tail of API logs"
    echo "  system    - Show system resource usage"
    echo "  ports     - Show port usage"
    echo "  processes - Show running processes"
    echo "  clear     - Clear log files"
    echo ""
}

# Function to show API logs
show_api_logs() {
    echo -e "${BLUE}📋 API Server Logs${NC}"
    echo "================================"

    # Check if API is running
    if lsof -ti:8080 >/dev/null 2>&1; then
        echo -e "${GREEN}✅ API server is running on port 8080${NC}"
        echo ""
        echo -e "${YELLOW}Recent API activity:${NC}"

        # Show recent requests from access logs if available
        if [ -f "access.log" ]; then
            tail -n 20 access.log
        else
            echo "No access.log file found. API logs are shown in the terminal running the server."
        fi
    else
        echo -e "${RED}❌ API server is not running${NC}"
    fi
}

# Function to show frontend logs
show_frontend_logs() {
    echo -e "${BLUE}🎨 Frontend Development Server Logs${NC}"
    echo "================================"

    # Check if frontend dev server is running
    if lsof -ti:3000 >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Frontend dev server is running on port 3000${NC}"
        echo ""
        echo -e "${YELLOW}Frontend build/serve logs are shown in the terminal running npm run serve${NC}"
    else
        echo -e "${RED}❌ Frontend dev server is not running${NC}"
    fi

    # Show build logs if available
    if [ -f "frontend/npm-debug.log" ]; then
        echo -e "${YELLOW}NPM Debug Log:${NC}"
        tail -n 20 frontend/npm-debug.log
    fi
}

# Function to live tail logs
live_tail() {
    echo -e "${BLUE}📊 Live Log Monitoring${NC}"
    echo "Press Ctrl+C to stop"
    echo "================================"

    # Monitor API access if log file exists
    if [ -f "access.log" ]; then
        echo -e "${GREEN}Monitoring access.log...${NC}"
        tail -f access.log
    else
        echo -e "${YELLOW}No access.log file found.${NC}"
        echo "API logs are shown in the terminal running the server."
        echo ""
        echo -e "${BLUE}Monitoring system resources instead...${NC}"

        while true; do
            clear
            echo -e "${BLUE}📊 System Monitor - $(date)${NC}"
            echo "================================"

            # Show running processes
            echo -e "${YELLOW}Domain Check Processes:${NC}"
            ps aux | grep -E "(domaincheck|cmd/server|8080|3000)" | grep -v grep || echo "No processes found"

            echo ""
            echo -e "${YELLOW}Port Usage:${NC}"
            lsof -i :8080 2>/dev/null || echo "Port 8080: Not in use"
            lsof -i :3000 2>/dev/null || echo "Port 3000: Not in use"

            echo ""
            echo -e "${YELLOW}Memory Usage:${NC}"
            ps aux | grep -E "(domaincheck|cmd/server)" | grep -v grep | awk '{print $2, $3, $4, $11}' | head -5

            sleep 5
        done
    fi
}

# Function to show system info
show_system_info() {
    echo -e "${BLUE}💻 System Information${NC}"
    echo "================================"

    echo -e "${YELLOW}Go Version:${NC}"
    go version 2>/dev/null || echo "Go not found"

    echo -e "${YELLOW}Node.js Version:${NC}"
    node --version 2>/dev/null || echo "Node.js not found"

    echo -e "${YELLOW}npm Version:${NC}"
    npm --version 2>/dev/null || echo "npm not found"

    echo ""
    echo -e "${YELLOW}System Resources:${NC}"
    echo "CPU: $(sysctl -n hw.ncpu 2>/dev/null || echo 'Unknown') cores"
    echo "Memory: $(sysctl -n hw.memsize 2>/dev/null | awk '{print $1/1024/1024/1024 " GB"}' || echo 'Unknown')"

    echo ""
    echo -e "${YELLOW}Disk Usage:${NC}"
    df -h . | tail -1 | awk '{print "Free: " $4 " / Total: " $2}'
}

# Function to show port usage
show_ports() {
    echo -e "${BLUE}🔌 Port Usage${NC}"
    echo "================================"

    echo -e "${YELLOW}Port 8080 (API Server):${NC}"
    lsof -i :8080 2>/dev/null || echo "Not in use"

    echo ""
    echo -e "${YELLOW}Port 3000 (Frontend Dev Server):${NC}"
    lsof -i :3000 2>/dev/null || echo "Not in use"

    echo ""
    echo -e "${YELLOW}All Node.js and Go processes:${NC}"
    lsof -i -P | grep -E "(node|go|:8080|:3000)" | head -10
}

# Function to show processes
show_processes() {
    echo -e "${BLUE}⚙️  Running Processes${NC}"
    echo "================================"

    echo -e "${YELLOW}Domain Check related processes:${NC}"
    ps aux | grep -E "(domaincheck|cmd/server|node.*serve)" | grep -v grep

    echo ""
    echo -e "${YELLOW}Process tree:${NC}"
    if command -v pstree >/dev/null 2>&1; then
        pstree -p $(pgrep -f "domaincheck\|cmd/server\|npm.*serve" 2>/dev/null | head -1) 2>/dev/null || echo "No process tree available"
    else
        echo "pstree not available"
    fi
}

# Function to clear logs
clear_logs() {
    echo -e "${BLUE}🧹 Clearing Log Files${NC}"
    echo "================================"

    # Clear access logs
    if [ -f "access.log" ]; then
        > access.log
        echo -e "${GREEN}✅ Cleared access.log${NC}"
    fi

    # Clear npm logs
    if [ -f "frontend/npm-debug.log" ]; then
        rm frontend/npm-debug.log
        echo -e "${GREEN}✅ Removed npm-debug.log${NC}"
    fi

    # Clear system logs if any
    rm -f *.log 2>/dev/null

    echo -e "${GREEN}✅ Log files cleared${NC}"
}

# Main script logic
case "$1" in
    "api")
        show_api_logs
        ;;
    "frontend")
        show_frontend_logs
        ;;
    "live")
        live_tail
        ;;
    "system")
        show_system_info
        ;;
    "ports")
        show_ports
        ;;
    "processes")
        show_processes
        ;;
    "clear")
        clear_logs
        ;;
    *)
        show_usage
        echo ""
        echo -e "${BLUE}Quick Status:${NC}"
        echo "================================"

        # Show current status
        if lsof -ti:8080 >/dev/null 2>&1; then
            echo -e "${GREEN}✅ API Server: Running (port 8080)${NC}"
        else
            echo -e "${RED}❌ API Server: Not running${NC}"
        fi

        if lsof -ti:3000 >/dev/null 2>&1; then
            echo -e "${GREEN}✅ Frontend Dev Server: Running (port 3000)${NC}"
        else
            echo -e "${RED}❌ Frontend Dev Server: Not running${NC}"
        fi

        echo ""
        echo -e "${YELLOW}💡 Use specific commands for detailed logs:${NC}"
        echo "   ./scripts/logs.sh api       - API logs"
        echo "   ./scripts/logs.sh live      - Live monitoring"
        echo "   ./scripts/logs.sh system    - System info"
        ;;
esac
