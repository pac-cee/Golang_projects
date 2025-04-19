# Fullstack Todo App (React + Go)

This project is a simple fullstack Todo application with a React frontend and a Go backend. Both services are orchestrated with Docker Compose.

## Features
- Add, edit, delete, and toggle todos
- REST API backend with Go
- Modern React frontend
- Docker Compose for easy orchestration

## Getting Started

### Prerequisites
- Docker and Docker Compose installed

### Running the App

1. Clone the repo or copy the files.
2. In the project root, run:

```sh
docker-compose up --build
```

3. Access the frontend at [http://localhost:3000](http://localhost:3000)
4. The backend API runs at [http://localhost:8080](http://localhost:8080)

## Project Structure

- `backend/` - Go REST API
- `frontend/` - React app
- `docker-compose.yml` - Orchestration

## API Endpoints
- GET /todos
- POST /todos
- PUT /todos/{id}
- DELETE /todos/{id}

## Customization
- You can switch the backend to use a database (e.g., SQLite or PostgreSQL) for persistence.
- Add authentication, user accounts, etc.

## License
MIT
