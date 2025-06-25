#!/bin/bash

# Test Domain Check API endpoints

echo "🧪 Testing Domain Check API..."

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

API_BASE="http://localhost:8080/api"
API_V1_BASE="http://localhost:8080/api/v1"

# Function to make API call and show result
test_endpoint() {
    local method=$1
    local url=$2
    local description=$3
    local data=$4

    echo -e "${BLUE}🔗 Testing: $description${NC}"
    echo -e "${YELLOW}   $method $url${NC}"

    if [ -n "$data" ]; then
        echo -e "${YELLOW}   Data: $data${NC}"
        response=$(curl -s -w "\n%{http_code}" -X $method -H "Content-Type: application/json" -d "$data" "$url" 2>/dev/null || echo -e "\nERROR")
    else
        response=$(curl -s -w "\n%{http_code}" -X $method "$url" 2>/dev/null || echo -e "\nERROR")
    fi

    # Split response and status code
    body=$(echo "$response" | head -n -1)
    status=$(echo "$response" | tail -n 1)

    if [ "$status" = "ERROR" ]; then
        echo -e "${RED}   ❌ Connection failed${NC}"
        return 1
    fi

    # Check status code
    if [[ $status -ge 200 && $status -lt 300 ]]; then
        echo -e "${GREEN}   ✅ Status: $status${NC}"
        # Pretty print JSON if possible
        if command -v jq >/dev/null 2>&1; then
            echo "$body" | jq . 2>/dev/null || echo "$body"
        else
            echo "$body"
        fi
    else
        echo -e "${RED}   ❌ Status: $status${NC}"
        echo "$body"
        return 1
    fi

    echo ""
    return 0
}

# Check if API server is running
echo -e "${BLUE}🔍 Checking if API server is running...${NC}"
if ! curl -s http://localhost:8080 >/dev/null; then
    echo -e "${RED}❌ API server is not running on localhost:8080${NC}"
    echo -e "${YELLOW}💡 Start the server with: ./scripts/start-api.sh${NC}"
    exit 1
fi

echo -e "${GREEN}✅ API server is running${NC}"
echo ""

# Run tests
echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}🧪 Running API Tests${NC}"
echo -e "${BLUE}================================${NC}"
echo ""

# Test health endpoint
test_endpoint "GET" "$API_BASE/health" "Health Check"

# Test domain check
test_endpoint "POST" "$API_BASE/check-domain" "Domain Check (google.com)" '{"domain": "google.com"}'

# Test ALL EXTENSIONS check (main feature) - Legacy API
test_endpoint "POST" "$API_BASE/check-all-extensions" "Check All Extensions (testdomain88888)" '{"domain_name": "testdomain88888"}'

# Test invalid domain
test_endpoint "POST" "$API_BASE/check-domain" "Invalid Domain Check" '{"domain": "invalid..domain"}'

# Test empty domain
test_endpoint "POST" "$API_BASE/check-domain" "Empty Domain Check" '{"domain": ""}'

# Test domain history
test_endpoint "GET" "$API_BASE/domains" "Domain History"

# Test v1 endpoints
echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}🧪 Testing V1 API Endpoints${NC}"
echo -e "${BLUE}================================${NC}"
echo ""

# Test v1 health
test_endpoint "GET" "$API_V1_BASE/health" "V1 Health Check"

# Test v1 domain check
test_endpoint "POST" "$API_V1_BASE/domains/check" "V1 Domain Check (github.com)" '{"domain": "github.com"}'

# Test ALL EXTENSIONS check (main feature)
test_endpoint "POST" "$API_V1_BASE/domains/check-all-extensions" "V1 Check All Extensions (testdomain99999)" '{"domain_name": "testdomain99999"}'

# Test multiple domain check
test_endpoint "POST" "$API_V1_BASE/domains/check-multiple" "V1 Multiple Domain Check" '{"domains": ["apple.com", "microsoft.com", "nonexistentdomain12345.com"]}'

# Test domain history with pagination
test_endpoint "GET" "$API_V1_BASE/domains/history?page=1&per_page=5" "V1 Domain History (Paginated)"

# Test extensions endpoint
test_endpoint "GET" "$API_V1_BASE/extensions" "V1 Valid Extensions"

# Performance test
echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}⚡ Performance Test${NC}"
echo -e "${BLUE}================================${NC}"
echo ""

echo -e "${BLUE}🏃 Running performance test (5 concurrent requests)...${NC}"
start_time=$(date +%s%N)

# Run 5 concurrent requests
{
    curl -s -X POST -H "Content-Type: application/json" -d '{"domain": "test1.com"}' "$API_BASE/check-domain" >/dev/null &
    curl -s -X POST -H "Content-Type: application/json" -d '{"domain": "test2.com"}' "$API_BASE/check-domain" >/dev/null &
    curl -s -X POST -H "Content-Type: application/json" -d '{"domain": "test3.com"}' "$API_BASE/check-domain" >/dev/null &
    curl -s -X POST -H "Content-Type: application/json" -d '{"domain": "test4.com"}' "$API_BASE/check-domain" >/dev/null &
    curl -s -X POST -H "Content-Type: application/json" -d '{"domain": "test5.com"}' "$API_BASE/check-domain" >/dev/null &

    wait
}

end_time=$(date +%s%N)
duration=$(( (end_time - start_time) / 1000000 ))

echo -e "${GREEN}✅ Performance test completed in ${duration}ms${NC}"
echo ""

# Summary
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}🎉 API Test Summary${NC}"
echo -e "${GREEN}================================${NC}"
echo ""
echo -e "${YELLOW}Tested endpoints:${NC}"
echo "   • Health check"
echo "   • Domain validation"
echo "   • Domain history"
echo "   • V1 API endpoints"
echo "   • Extensions management"
echo "   • Error handling"
echo "   • Performance (concurrent requests)"
echo ""
echo -e "${BLUE}💡 Tips:${NC}"
echo "   • Use jq for better JSON formatting: brew install jq"
echo "   • Check logs in terminal running the API server"
echo "   • Monitor API performance with: ./scripts/monitor.sh"
