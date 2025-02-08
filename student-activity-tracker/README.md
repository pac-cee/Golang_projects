# Student Activity Tracker

A comprehensive web application that helps students manage and track their daily activities, set timers, and monitor progress. Built with modern technologies and following best practices in both frontend and backend development.

## Features

- Add activities with detailed information (title, description, duration)
- Track activity status (planned, in-progress, completed)
- Real-time activity timer
- Activity progress monitoring
- Clean and modern UI with responsive design
- RESTful API backend
- MongoDB for persistent storage
- Comprehensive test coverage

## Tech Stack

### Frontend
- **Astro**: Modern static site generator with excellent performance
  - Server-first architecture
  - Zero-JS by default for better performance
  - Supports multiple UI frameworks
- **React**: For interactive components
  - Activity form with real-time validation
  - Activity list with status management
  - Timer component
- **TailwindCSS**: Utility-first CSS framework
  - Responsive design
  - Modern UI components
  - Easy customization
- **TypeScript**: For type safety and better developer experience

### Backend
- **Go**: High-performance backend language
  - Strong typing and compilation
  - Excellent concurrency support
  - Great for building APIs
- **Gin**: Web framework for Go
  - Fast HTTP routing
  - Middleware support
  - Built-in validation
- **MongoDB**: NoSQL database
  - Flexible schema
  - Great for document-based data
  - Excellent Go driver support

## Project Structure

```
student-activity-tracker/
├── backend/
│   ├── config/
│   │   └── config.go           # Application configuration
│   ├── database/
│   │   └── mongodb.go          # Database connection management
│   ├── handlers/
│   │   ├── activity_handler.go # HTTP request handlers
│   │   └── activity_handler_test.go
│   ├── models/
│   │   └── activity.go         # Data models
│   ├── services/
│   │   ├── activity_service.go # Business logic
│   │   └── activity_service_test.go
│   ├── go.mod                  # Go module definition
│   └── main.go                 # Application entry point
└── frontend/
    ├── src/
    │   ├── components/
    │   │   ├── ActivityForm.tsx  # Form for creating activities
    │   │   └── ActivityList.tsx  # List of activities with controls
    │   └── pages/
    │       └── index.astro       # Main page
    ├── astro.config.mjs          # Astro configuration
    └── tailwind.config.mjs       # Tailwind configuration
```

## Architecture

### Backend Architecture

1. **Config Package**
   - Manages application configuration
   - Environment variable handling
   - Default configuration values
   - Type-safe configuration structure

2. **Database Package**
   - MongoDB connection management
   - Connection pooling
   - Database initialization
   - Collection management

3. **Models Package**
   - Defines data structures
   - JSON/BSON serialization
   - Type definitions
   - Data validation

4. **Services Package**
   - Business logic implementation
   - CRUD operations
   - Activity state management
   - Error handling

5. **Handlers Package**
   - HTTP request handling
   - Request validation
   - Response formatting
   - Route management

### Frontend Architecture

1. **Components**
   - Reusable React components
   - TypeScript for type safety
   - Tailwind for styling
   - Client-side interactivity

2. **Pages**
   - Astro pages for static generation
   - Hybrid static/dynamic content
   - SEO optimization
   - Performance optimization

## API Endpoints

### Activities

```
POST /api/activities
- Create a new activity
- Request: { title: string, description: string, duration: number }
- Response: Activity object

GET /api/activities
- Get all activities
- Response: Array of Activity objects

PUT /api/activities/:id
- Update an activity
- Request: { title?: string, description?: string, duration?: number }
- Response: Updated Activity object

DELETE /api/activities/:id
- Delete an activity
- Response: { message: string }

PUT /api/activities/:id/start
- Start an activity
- Response: Updated Activity object

PUT /api/activities/:id/complete
- Complete an activity
- Response: Updated Activity object
```

## Data Models

### Activity Model
```typescript
interface Activity {
  id: string;
  title: string;
  description: string;
  duration: number;      // in minutes
  status: string;        // planned, in-progress, completed
  startTime?: Date;      // when activity was started
  endTime?: Date;        // when activity was completed
  createdAt: Date;
  updatedAt: Date;
}
```

## Testing

### Backend Tests
- Unit tests for services
- Integration tests for handlers
- Database operation tests
- Mock interfaces for dependencies

### Frontend Tests
- Component testing with React Testing Library
- End-to-end testing capabilities
- Accessibility testing

## Installation and Setup

1. **Prerequisites**
   ```bash
   # Install Go
   go version # Should be 1.21 or later

   # Install Node.js
   node --version # Should be 18 or later

   # Install MongoDB
   mongod --version # Should be 6.0 or later
   ```

2. **Backend Setup**
   ```bash
   cd backend
   go mod tidy
   go run main.go
   ```

3. **Frontend Setup**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

4. **Environment Variables**
   ```bash
   # Backend (.env)
   MONGO_URI=mongodb://localhost:27017
   DB_NAME=student_tracker
   SERVER_ADDRESS=:8080
   ALLOWED_ORIGIN=http://localhost:4321
   ```

## Development Workflow

1. **Backend Development**
   - Write tests first (TDD approach)
   - Implement features
   - Run tests: `go test ./...`
   - Manual testing with API client

2. **Frontend Development**
   - Create/modify components
   - Test in development mode
   - Build for production
   - Test production build

## Best Practices

### Backend
- Clean Architecture principles
- Dependency Injection
- Error handling
- Input validation
- Proper logging
- Security best practices

### Frontend
- Component reusability
- Performance optimization
- Progressive enhancement
- Accessibility
- Responsive design

## Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Submit Pull Request

## License

This project is licensed under the MIT License.
