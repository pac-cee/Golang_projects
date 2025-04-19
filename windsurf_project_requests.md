# Windsurf Project Requests: Beginner to Advanced

This document contains a curated set of project requests you can give to your Windsurf agent. These cover frontend, backend, DevOps, and fullstack development, progressing from beginner to advanced. Each project includes a brief, a stack suggestion, and best practices to follow. Copy-paste any of these to your agent for implementation and learning!

---

## 1. Beginner Projects

### Hello World CLI
**Brief:** Create a Hello World command-line app in Go, Python, and JavaScript.
**Stack:** Go, Python, Node.js
**Requirements:**
- Print "Hello, World!" to the console.
- Add a README with run instructions.

---

### Simple Calculator
**Brief:** Build a CLI calculator that supports addition, subtraction, multiplication, and division.
**Stack:** Go, Python, Node.js
**Requirements:**
- Accept user input for numbers and operation.
- Handle invalid input gracefully.
- Add tests for calculator functions.

---

### Static Portfolio Website
**Brief:** Build a personal portfolio website.
**Stack:** HTML, CSS, JavaScript
**Requirements:**
- Responsive design.
- About, Projects, Contact sections.
- Use semantic HTML and CSS best practices.

---

## 2. Intermediate Projects

### RESTful API
**Brief:** Create a REST API for managing tasks (CRUD).
**Stack:** Go (net/http), Node.js (Express), Python (Flask)
**Requirements:**
- CRUD endpoints for tasks.
- Use environment variables for config.
- Modular project structure.
- Add OpenAPI/Swagger docs.

---

### Todo Web App (Fullstack)
**Brief:** Build a Todo app with a frontend and backend.
**Stack:** React or Vue (frontend), Go/Node.js/Python (backend)
**Requirements:**
- REST API for todos.
- Frontend with add, edit, delete, mark done.
- Use Docker Compose to run both services.
- Add tests for backend and frontend.

---

### Authentication System
**Brief:** Implement user authentication with JWT.
**Stack:** Go, Node.js, Python
**Requirements:**
- Register, login, logout endpoints.
- Password hashing and JWT token issuance.
- Protect routes with middleware.
- Add tests for authentication.

---

### File Uploader
**Brief:** Allow users to upload files via a web interface and store them on the server.
**Stack:** Go, Node.js, Python
**Requirements:**
- Frontend for file selection.
- Backend for file handling and storage.
- Validate file types and size.

---

## 3. Advanced Projects

### Real-Time Chat App
**Brief:** Build a chat app with real-time messaging.
**Stack:** Go (Gorilla WebSocket), Node.js (Socket.io), Python (websockets)
**Requirements:**
- User can join rooms and send messages.
- Show online users.
- Store chat history in a database.

---

### E-commerce Platform
**Brief:** Create a simple e-commerce system.
**Stack:** React/Vue (frontend), Go/Node.js/Python (backend), PostgreSQL/MongoDB
**Requirements:**
- Product listing, cart, checkout.
- User authentication and roles.
- Payment integration (mock or Stripe sandbox).
- Admin dashboard.

---

### Microservices Architecture
**Brief:** Refactor a monolithic Todo API into microservices.
**Stack:** Go/Node.js, Docker, REST/gRPC
**Requirements:**
- Separate services for users, tasks, notifications.
- API Gateway for routing.
- Use Docker Compose for orchestration.

---

### DevOps Pipeline
**Brief:** Set up CI/CD for any of the above projects.
**Stack:** GitHub Actions, Docker, Kubernetes (optional)
**Requirements:**
- Automated build, test, and deploy.
- Use Docker for containerization.
- Add README with pipeline instructions.

---

### Data Visualization Dashboard
**Brief:** Build a dashboard to display data from an external API.
**Stack:** React/D3.js (frontend), Go/Node.js (backend)
**Requirements:**
- Fetch and display data in charts.
- Handle large datasets efficiently.
- Responsive design.

---

## 4. Expert/Scalable Projects

### Distributed Task Queue
**Brief:** Implement a background job system.
**Stack:** Go/Python, Redis/RabbitMQ
**Requirements:**
- Submit jobs via API.
- Workers process jobs asynchronously.
- Track job status.

---

### Blockchain Demo
**Brief:** Build a basic blockchain implementation.
**Stack:** Go, Python
**Requirements:**
- Create blocks, proof-of-work, and peer-to-peer networking.
- Simple web interface to view chain.

---

### SaaS Multi-Tenant Platform
**Brief:** Create a multi-tenant SaaS starter.
**Stack:** React (frontend), Go/Node.js/Python (backend), PostgreSQL
**Requirements:**
- User and tenant separation.
- Billing (mock or Stripe test).
- Admin dashboard.

---

### Machine Learning API
**Brief:** Deploy an ML model as an API.
**Stack:** Python (FastAPI/Flask + scikit-learn/TensorFlow)
**Requirements:**
- Expose prediction endpoint.
- Handle model versioning.
- Secure endpoints.

---

## 5. Best Practices & Modern Developer Rules

For every project, ask your agent to:
- Use version control (Git).
- Write a comprehensive README.md.
- Use .env for secrets/config.
- Write unit/integration tests.
- Use linters and code formatters.
- Apply SOLID, DRY, KISS principles.
- Use Docker for containerization.
- Set up CI/CD automation.
- Document APIs and architecture.
- Follow secure coding practices (input validation, XSS, CSRF, SQL injection, etc.).

---

## How to Use
- Copy any project brief above and give it to your Windsurf agent.
- Specify your preferred stack/language if you want.
- Ask for best practices, documentation, and enhancements.
- Use these projects as references for interviews, learning, or your portfolio.

---

**Good luck on your software engineering journey! If you want more detailed specs or starter repos for any project, just ask.**
