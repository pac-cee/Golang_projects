# Student Activity Tracker - Technical Documentation

## 1. Project Overview and Features

### Core Purpose
The Student Activity Tracker is a comprehensive web application designed to help students manage their time effectively by tracking and monitoring their daily activities. It provides real-time progress tracking, time management features, and intuitive activity organization.

### Key Features

#### 1. Activity Management
- **Activity Creation**
  - Title and description
  - Duration setting (in minutes)
  - Priority levels
  - Categories/tags
  - Custom notes

#### 2. Time Tracking
- **Timer Functionality**
  - Start/pause/resume capabilities
  - Real-time countdown
  - Visual progress indicators
  - Overtime tracking
  - Break management

#### 3. Progress Monitoring
- **Status Tracking**
  - Planned activities
  - In-progress activities
  - Completed activities
  - Overdue activities

#### 4. Analytics
- **Time Analysis**
  - Time spent per activity
  - Daily/weekly/monthly summaries
  - Productivity patterns
  - Completion rates

#### 5. User Experience
- **Interface Features**
  - Drag-and-drop organization
  - Filter and search capabilities
  - Responsive design
  - Dark/light mode
  - Keyboard shortcuts

## 2. Detailed Tech Stack

### Frontend Technologies

#### 1. Astro (v4.x)
- **Purpose**: Static Site Generator
- **Key Features Used**:
  - Component Islands architecture
  - Zero-JS by default
  - Automatic page optimization
  - Built-in asset optimization
  - TypeScript integration

#### 2. React (v18.x)
- **Purpose**: Interactive UI Components
- **Key Features Used**:
  - Hooks (useState, useEffect, useContext)
  - Custom hooks for timer logic
  - Context API for state management
  - Error boundaries
  - Suspense for loading states

#### 3. TailwindCSS (v3.x)
- **Purpose**: Styling Framework
- **Key Features Used**:
  - JIT (Just-In-Time) compilation
  - Custom theme configuration
  - Responsive design utilities
  - Dark mode support
  - Animation classes

#### 4. TypeScript (v5.x)
- **Purpose**: Type Safety
- **Key Features Used**:
  - Strict type checking
  - Interface definitions
  - Generic types
  - Type guards
  - Utility types

### Backend Technologies

#### 1. Go (v1.21)
- **Purpose**: Server-side Language
- **Key Features Used**:
  - Goroutines for concurrency
  - Channels for communication
  - Context for cancellation
  - Error handling patterns
  - Interface implementation

#### 2. Gin Web Framework
- **Purpose**: HTTP Router and Middleware
- **Key Features Used**:
  - Middleware chain
  - Route grouping
  - Request validation
  - Error management
  - Static file serving

#### 3. MongoDB (v6.0)
- **Purpose**: Database
- **Key Features Used**:
  - Document model
  - Aggregation pipeline
  - Indexing strategies
  - Transaction support
  - Change streams

## 3. Complete Project Structure

### Directory Organization
```
student-activity-tracker/
├── backend/
│   ├── config/
│   │   ├── config.go           # Configuration management
│   │   └── config_test.go      # Configuration tests
│   │
│   ├── database/
│   │   ├── mongodb.go          # MongoDB connection
│   │   ├── migrations/         # Database migrations
│   │   └── seeds/             # Test data seeds
│   │
│   ├── handlers/
│   │   ├── activity_handler.go # HTTP handlers
│   │   ├── middleware/        # Custom middleware
│   │   └── validators/        # Request validators
│   │
│   ├── models/
│   │   ├── activity.go        # Data models
│   │   └── interfaces/        # Model interfaces
│   │
│   ├── services/
│   │   ├── activity_service.go # Business logic
│   │   └── interfaces/        # Service interfaces
│   │
│   └── utils/
│       ├── logger/            # Logging utilities
│       └── errors/            # Error handling
│
└── frontend/
    ├── src/
    │   ├── components/
    │   │   ├── activity/      # Activity-related components
    │   │   ├── common/        # Shared components
    │   │   └── layout/        # Layout components
    │   │
    │   ├── hooks/            # Custom React hooks
    │   │   ├── useTimer.ts
    │   │   └── useActivity.ts
    │   │
    │   ├── pages/
    │   │   ├── index.astro    # Main page
    │   │   └── analytics.astro # Analytics page
    │   │
    │   ├── styles/           # Global styles
    │   │   └── tailwind.css
    │   │
    │   └── utils/            # Utility functions
    │       ├── api.ts        # API client
    │       └── formatters.ts # Data formatters
    │
    └── public/              # Static assets
```

## 4. Architecture Details

### Backend Architecture

#### 1. Clean Architecture Layers
```
[HTTP Handlers] → [Services] → [Models] → [Database]
     ↑              ↑            ↑           ↑
     └──────────────┴────────────┴───────────┘
         Dependency Injection Flow
```

#### 2. Component Responsibilities

##### Config Layer
- Environment variable management
- Configuration validation
- Default values
- Type-safe configuration

##### Database Layer
- Connection management
- Query execution
- Transaction handling
- Error handling
- Connection pooling

##### Models Layer
- Data structure definitions
- Validation rules
- Type conversions
- Database mappings

##### Services Layer
- Business logic implementation
- Transaction coordination
- Error handling
- Event processing
- Data transformation

##### Handlers Layer
- Request parsing
- Response formatting
- Error handling
- Authentication
- Rate limiting

### Frontend Architecture

#### 1. Component Architecture
```
[Pages] → [Layouts] → [Feature Components] → [Base Components]
   ↑          ↑              ↑                     ↑
   └──────────┴──────────────┴─────────────────────┘
        Shared Hooks and Utilities
```

#### 2. State Management
- Local component state
- React Context for global state
- Custom hooks for complex logic
- Props for component communication

## 5. API Endpoints

### Activity Endpoints

#### Create Activity
```http
POST /api/activities
Content-Type: application/json

Request:
{
  "title": "Study Mathematics",
  "description": "Chapter 5: Linear Algebra",
  "duration": 60,
  "category": "study",
  "priority": "high"
}

Response: (201 Created)
{
  "id": "507f1f77bcf86cd799439011",
  "title": "Study Mathematics",
  "description": "Chapter 5: Linear Algebra",
  "duration": 60,
  "status": "planned",
  "category": "study",
  "priority": "high",
  "createdAt": "2025-02-08T16:34:22Z",
  "updatedAt": "2025-02-08T16:34:22Z"
}
```

#### Get Activities
```http
GET /api/activities
Query Parameters:
- status: string (planned|in-progress|completed)
- category: string
- priority: string
- startDate: ISO date
- endDate: ISO date

Response: (200 OK)
{
  "activities": [
    {
      "id": "507f1f77bcf86cd799439011",
      "title": "Study Mathematics",
      "description": "Chapter 5: Linear Algebra",
      "duration": 60,
      "status": "planned",
      "category": "study",
      "priority": "high",
      "createdAt": "2025-02-08T16:34:22Z",
      "updatedAt": "2025-02-08T16:34:22Z"
    }
  ],
  "total": 1,
  "page": 1,
  "pageSize": 10
}
```

#### Update Activity
```http
PUT /api/activities/:id
Content-Type: application/json

Request:
{
  "title": "Updated Title",
  "description": "Updated Description",
  "duration": 45
}

Response: (200 OK)
{
  "id": "507f1f77bcf86cd799439011",
  "title": "Updated Title",
  "description": "Updated Description",
  "duration": 45,
  "status": "planned",
  "updatedAt": "2025-02-08T16:35:22Z"
}
```

## 6. Data Models

### Activity Model
```go
type Activity struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Title       string            `bson:"title" json:"title" validate:"required,min=3,max=100"`
    Description string            `bson:"description" json:"description" validate:"required,max=500"`
    Duration    int               `bson:"duration" json:"duration" validate:"required,min=1"`
    Status      string            `bson:"status" json:"status" validate:"required,oneof=planned in-progress completed"`
    Category    string            `bson:"category" json:"category" validate:"required"`
    Priority    string            `bson:"priority" json:"priority" validate:"required,oneof=low medium high"`
    StartTime   *time.Time        `bson:"startTime,omitempty" json:"startTime,omitempty"`
    EndTime     *time.Time        `bson:"endTime,omitempty" json:"endTime,omitempty"`
    CreatedAt   time.Time         `bson:"createdAt" json:"createdAt"`
    UpdatedAt   time.Time         `bson:"updatedAt" json:"updatedAt"`
}
```

### Database Indexes
```javascript
// MongoDB Indexes
{
  "status": 1,
  "createdAt": -1
}
{
  "category": 1,
  "status": 1
}
{
  "priority": 1,
  "status": 1
}
```

## 7. Testing Strategies

### Backend Testing

#### 1. Unit Tests
```go
// Example service test
func TestActivityService_CreateActivity(t *testing.T) {
    // Test cases
    tests := []struct {
        name    string
        input   *models.Activity
        wantErr bool
    }{
        {
            name: "Valid activity",
            input: &models.Activity{
                Title:       "Test Activity",
                Description: "Test Description",
                Duration:    30,
            },
            wantErr: false,
        },
        // More test cases...
    }

    // Run tests
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

#### 2. Integration Tests
```go
// Example handler test
func TestActivityHandler_CreateActivity(t *testing.T) {
    // Setup test server
    router := setupTestRouter()
    
    // Test cases
    tests := []struct {
        name       string
        payload    map[string]interface{}
        wantStatus int
    }{
        {
            name: "Valid request",
            payload: map[string]interface{}{
                "title":       "Test Activity",
                "description": "Test Description",
                "duration":    30,
            },
            wantStatus: http.StatusCreated,
        },
        // More test cases...
    }

    // Run tests
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Frontend Testing

#### 1. Component Tests
```typescript
// Example React component test
describe('ActivityForm', () => {
    it('submits form with valid data', async () => {
        const onSubmit = jest.fn();
        render(<ActivityForm onSubmit={onSubmit} />);

        // Fill form
        await userEvent.type(
            screen.getByLabelText(/title/i),
            'Test Activity'
        );

        // Submit form
        await userEvent.click(screen.getByRole('button', { name: /submit/i }));

        // Assertions
        expect(onSubmit).toHaveBeenCalledWith({
            title: 'Test Activity',
            // Other fields...
        });
    });
});
```

#### 2. Integration Tests
```typescript
// Example page test
describe('ActivityPage', () => {
    it('loads and displays activities', async () => {
        // Mock API response
        server.use(
            rest.get('/api/activities', (req, res, ctx) => {
                return res(ctx.json({
                    activities: [
                        {
                            id: '1',
                            title: 'Test Activity',
                            // Other fields...
                        }
                    ]
                }));
            })
        );

        render(<ActivityPage />);

        // Wait for activities to load
        await screen.findByText('Test Activity');

        // Assertions
        expect(screen.getByText('Test Activity')).toBeInTheDocument();
    });
});
```

This documentation provides a comprehensive overview of the technical aspects of the Student Activity Tracker. Each section can be further expanded based on specific needs or questions.
