# Domain Check

A simple domain availability checking API built with Go and a modern Vue.js frontend interface.

## Features

- Domain availability checking
- DNS lookup for IP address resolution
- Modern Vue.js interface with real-time updates
- RESTful API with WebSocket support
- CORS enabled
- Check history viewing
- WHOIS information lookup
- Bulk domain extension checking with real-time progress
- Modular API architecture
- **CI/CD Pipeline** with GitHub Actions
- **Docker** support for easy deployment
- **Dependabot** for automated dependency updates

## 🚀 Quick Start

### 🐳 **Docker (Recommended)**
```bash
# Production
docker-compose up -d

# Development with hot reload
docker-compose --profile dev up -d
```

### ⚡ **Start with Single Command**
```bash
./scripts/dev.sh
```
This command automatically:
- Installs dependencies
- Starts the API server (port 8080)
- Starts the frontend dev server (port 3000)
- Provides monitoring

### 🔧 **Manual Setup**

#### Backend (Go API)
```bash
go mod tidy
./scripts/start-api.sh
```

#### Frontend (Vue.js)
```bash
./scripts/start-frontend.sh
```

#### Production Build
```bash
./scripts/build.sh
```

### 📋 **All Scripts**
- `./scripts/dev.sh` - Complete development environment
- `./scripts/build.sh` - Production build
- `./scripts/test-api.sh` - API tests
- `./scripts/stop.sh` - Stop services
- `./scripts/logs.sh` - Logs and monitoring

For details: [scripts/README.md](scripts/README.md)

## 🐳 Docker

### Production Build
```bash
docker build -t domaincheck .
docker run -p 8080:8080 domaincheck
```

### Development with Hot Reload
```bash
docker-compose --profile dev up -d
```

### Multi-stage Build
The Dockerfile uses multi-stage builds for optimized production images:
- Frontend build stage
- Backend build stage
- Production stage with minimal footprint

## 🔄 CI/CD Pipeline

This project includes a comprehensive GitHub Actions CI/CD pipeline:

### Workflow Features
- **Automated Testing**: Backend and frontend tests
- **Code Quality**: Linting and formatting checks
- **Security Scanning**: Trivy vulnerability scanner
- **Multi-stage Deployment**: Staging and production environments
- **Build Artifacts**: Optimized build outputs

### Pipeline Stages
1. **Test**: Run backend and frontend tests
2. **Build**: Create production artifacts
3. **Security**: Vulnerability scanning
4. **Deploy Staging**: Auto-deploy to staging (develop branch)
5. **Deploy Production**: Auto-deploy to production (main branch)

### Environment Protection
- Staging deployment requires develop branch
- Production deployment requires main branch
- Manual approval for production deployments

## 🤖 Dependabot

Automated dependency updates are configured for:
- **Go modules**: Weekly updates
- **Node.js packages**: Weekly updates
- **Security updates**: Immediate notifications

### Configuration
- Updates scheduled for Mondays at 9:00 AM
- Pull request limit: 10 concurrent updates
- Automatic labeling and assignment
- Major version updates ignored for critical packages

## API Endpoints

### Health Check
- `GET /api/v1/health` - API health check

### Domain Operations
- `POST /api/v1/domains/check` - Check single domain
- `POST /api/v1/domains/check-all-extensions` - Check domain with all extensions
- `POST /api/v1/domains/check-multiple` - Check multiple domains
- `GET /api/v1/domains/history` - Get check history
- `DELETE /api/v1/domains/history` - Clear history
- `GET /api/v1/domains/whois/:domain` - Get WHOIS information

### Extensions Management
- `GET /api/v1/extensions` - Get valid extensions
- `POST /api/v1/extensions/reload` - Reload extensions

### WebSocket
- `WS /ws` - WebSocket connection for real-time updates

### Backward Compatibility
- `GET /api/health` - Health check (v0)
- `POST /api/check-domain` - Check domain (v0)
- `POST /api/check-all-extensions` - Check all extensions (v0)
- `GET /api/domains` - Get history (v0)

## Usage

1. Start backend and frontend servers
2. Navigate to `http://localhost:3000`
3. Enter a domain name and check availability
4. View results in the history section
5. Use WHOIS button to get detailed domain information
6. Use WebSocket tab for real-time bulk checking

## Technologies

- **Backend**: Go, Gin framework, WebSocket
- **Frontend**: Vue.js 3, Axios, Tailwind CSS
- **Real-time**: WebSocket for live updates
- **CORS**: Frontend-backend communication
- **CI/CD**: GitHub Actions
- **Containerization**: Docker, Docker Compose
- **Security**: Trivy vulnerability scanner

## Architecture

The project follows a modular architecture:

```
DomainCheck/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/         # Configuration management
│   ├── handlers/       # HTTP request handlers
│   ├── models/         # Data models
│   ├── services/       # Business logic
│   └── utils/          # Utility functions
├── frontend/           # Vue.js frontend
├── configs/            # Configuration files
├── data/               # Data files (domain extensions)
├── scripts/            # Build and deployment scripts
├── .github/            # CI/CD and GitHub configurations
├── Dockerfile          # Production Docker image
├── Dockerfile.dev      # Development Docker image
└── docker-compose.yml  # Multi-service orchestration
```

## Development

### Prerequisites
- Go 1.19+
- Node.js 18+
- npm or yarn
- Docker (optional)

### Local Development
1. Clone the repository
2. Run `./scripts/dev.sh` for full development environment
3. Access the application at `http://localhost:3000`

### API Testing
```bash
./scripts/test-api.sh
```

### Hot Reload Development
```bash
# Install air for Go hot reload
go install github.com/cosmtrek/air@latest

# Start with hot reload
air
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Issue Templates
- Bug reports: Use the bug report template
- Feature requests: Use the feature request template
- Pull requests: Follow the PR template

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
