App Showcase

A full-stack application for browsing apps and their reviews, built with React frontend and Go backend.
🚀 Features

    App Discovery: Browse and search through a curated list of applications
    Review System: View detailed reviews and ratings for each app
    Time-based Filtering: Filter reviews by time periods (24h, 48h, 72h, 96h, or all)
    Responsive Design: Works seamlessly on desktop and mobile devices
    Real-time Data: Fresh app data and reviews fetched from backend API

🏗️ Architecture

    Frontend: React 18 with React Router for navigation
    Backend: Go REST API with JSON data storage
    Containerization: Docker and Docker Compose for easy deployment
    Networking: Container-based communication with CORS support

📋 Prerequisites

Make sure you have the following installed on your system:

    Docker (version 20.0 or higher)
    Docker Compose (version 2.0 or higher)
    Git

🛠️ Quick Start
1. Clone the Repository
   bash

`git clone <your-repository-url>`
cd <repository-name>

2. Project Structure

Ensure your project structure looks like this:

```project-root/
├── docker-compose.yml
├── README.md
├── back-end/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── .env
│   └── [other Go files]
└── front-end/
├── Dockerfile
├── package.json
├── .env
├── public/
└── src/
├── App.js
├── components/
└── [other React files]
```

3. Environment Configuration

Create the necessary environment files:

Backend (`.env` in `back-end/` directory):
env

# Add your backend environment variables here
`PORT=8080`
# Add any API keys, database URLs, etc.

Frontend (`.env` in `front-end/` directory):
env

`REACT_APP_API_URL=http://localhost:8080`

4. Build and Run
   bash

# Build and start all services
`docker-compose up --build`

# Or run in detached mode (background)
`docker-compose up --build -d`

5. Access the Application

   * Frontend: Open your browser and navigate to http://localhost:3000
   * Backend API: Available at http://localhost:8080

📖 Usage Guide
Getting Started

    Launch the application by following the Quick Start guide above
    Load app data by clicking the "Get Apps" button on the main page
    Browse apps by scrolling through the app cards
    View reviews by clicking "View Reviews" on any app card
    Filter reviews using the time-based dropdown filter
    Navigate using browser back/forward buttons or the "Back to App List" button

API Endpoints

The backend exposes the following endpoints:

    GET /app/list - Retrieve list of available apps
    GET /app/reviews?id={appId}&hours={hours} - Get reviews for a specific app

Frontend Routes

    / - Main app list page
    /app/{appId} - App reviews page with optional ?hours={hours} query parameter

🔧 Development
Running in Development Mode

For active development with hot reloading:
bash

# Stop the production containers
`docker-compose down`

# Start in development mode (if you have a docker-compose.dev.yml)
`docker-compose -f docker-compose.dev.yml up --build`

Backend Development
```
bash

cd back-end
go mod download
go run main.go

Frontend Development
bash

cd front-end
npm install
npm start

```

🐳 Docker Configuration
Services

    backend-service: Go API server running on port 8080
    frontend-service: React development server running on port 3000

Networks

The application uses a custom Docker network (my-custom-network) with static IP addresses for reliable container communication.
Volumes

    Frontend source code is mounted for hot reloading during development
    Backend data and logs are persisted in container directories

🔍 Troubleshooting
Common Issues
"NetworkError when attempting to fetch resource"

Solution: Ensure the backend is running and accessible:
bash

# Test backend connectivity
curl http://localhost:8080/app/list

# Check container status
`docker-compose ps`

# View container logs
```
docker-compose logs backend-service
docker-compose logs frontend-service

```

"Module not found: Error: Can't resolve 'react-router-dom'"

Solution: Rebuild the frontend container:
```
bash

docker-compose down
docker-compose build frontend-service
docker-compose up
```

Frontend not accessible from browser

Solution: Ensure the HOST environment variable is set correctly:
bash

# Check if containers are running on correct ports
`docker-compose ps`

# Verify environment variables
`docker exec frontend-service env | grep HOST`

Debug Commands
```
bash

# View all running containers
docker-compose ps

# Check container logs
docker-compose logs [service-name]

# Execute commands inside containers
docker exec -it frontend-service /bin/sh
docker exec -it backend-service /bin/sh

# Restart specific service
docker-compose restart [service-name]

# Rebuild and restart
docker-compose down
docker-compose up --build
```

🧪 Testing
Manual Testing

    Backend API: Test endpoints using curl or Postman
    bash

    curl http://localhost:8080/app/list
    curl "http://localhost:8080/app/reviews?id=APP_ID&hours=24"

    Frontend: Test navigation and functionality in browser
        App list loading
        App selection and navigation
        Review filtering
        Browser back/forward buttons

Automated Testing
bash

# Frontend tests
```
cd front-end
npm test
```
# Backend tests (if available)
```
cd back-end
go test ./...
```
📦 Production Deployment
Building for Production
```
bash

# Build optimized production images
docker-compose -f docker-compose.prod.yml build

# Run in production mode
docker-compose -f docker-compose.prod.yml up -d

```
Environment Variables for Production

Update your environment files with production values:

    Database URLs
    API keys
    CORS origins
    SSL certificates


🚀 What's Next?

    Add user authentication/authorization
    Add pagination (back-end side)
    Make use of some database
    Implement end-to-end tests
    Performance optimization
    CI/CD pipeline