# Konex Local Setup Guide

Welcome to the Konex local setup guide! This document covers the prerequisites and step-by-step instructions for getting the Konex project running on your local machine.

## Prerequisites

Before you begin, ensure you have the following installed on your machine:

1. **Docker & Docker Compose**: The easiest way to run the entire application stack.
   - [Install Docker Desktop (Mac/Windows)](https://www.docker.com/products/docker-desktop)
   - [Install Docker Engine (Linux)](https://docs.docker.com/engine/install/)
2. **Git**: To clone and manage the repository.
   - [Install Git](https://git-scm.com/downloads)

*(Optional) If you want to run the Go backend or Node frontend natively without Docker:*
- **Go**: Version `1.24` or higher. [Install Go](https://go.dev/doc/install)
- **Node.js**: Version `20` or higher, along with `yarn` or `npm`. [Install Node.js](https://nodejs.org/)

---

## Step-by-Step Setup

### 1. Clone the Repository

Clone the project to your local machine and navigate into the root directory:

```bash
git clone https://github.com/ravitejas-tech/konex.git
cd konex
```

### 2. Configure Environment Variables

The project requires some environment variables to run properly. 

1. Navigate to the `api` folder and copy the example environment file:
   ```bash
   cd api
   cp .env.example .env
   ```
2. Navigate back to the root directory:
   ```bash
   cd ..
   ```

*Note: The default values in `.env.example` are already configured to work seamlessly with the Docker setup.*

### 3. Run the Services with Docker

The project uses Docker Compose to orchestrate both the Go/PocketBase API backend and the web frontend.

Run the following command from the root directory to build and start all services:

```bash
docker compose up --build
```

Docker will now build the Go API and the Vite React frontend. Once the build is complete, the services will start up.

### 4. Access the Application

Once the services are running, you can access them in your browser:

- **Web Frontend**: [http://localhost:3000](http://localhost:3000)
- **PocketBase Admin UI**: [http://localhost:8090/_/](http://localhost:8090/_/)
- **API Base URL**: `http://localhost:8090`

### 5. First-Time Setup: PocketBase Admin Login

When running the PocketBase API for the first time, you can log in using the predefined superuser credentials:

- **Email**: `admin@konex.local`
- **Password**: `Password123!`

*Note: You can change these credentials or create new admin accounts from within the Admin UI.*

---

## Native Development (Without Docker)

If you prefer to run the services directly on your local machine for faster development iteration, follow these steps:

### Running the API (Go + PocketBase)

1. Open a new terminal and navigate to the `api` directory:
   ```bash
   cd api
   ```
2. Ensure dependencies are installed and the `.env` file exists:
   ```bash
   go mod download
   ```
3. Run the API server:
   ```bash
   go run ./cmd/api serve --http=0.0.0.0:8090
   ```

### Running the Web Frontend (Node.js)

1. Open another terminal and navigate to the `web` directory:
   ```bash
   cd web
   ```
2. Install the dependencies:
   ```bash
   npm install  # or yarn install
   ```
3. Start the development server:
   ```bash
   npm run dev  # or yarn dev
   ```

## Stopping the Application

To gracefully stop the Docker containers, press `Ctrl + C` in the terminal where `docker compose` is running. 

If you started the containers in detached mode (`-d`), you can stop them using:

```bash
docker compose down
```
