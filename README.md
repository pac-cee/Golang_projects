# Golang Projects

This repository contains various Go projects and learning materials.

## Expense Tracker

A full-stack expense tracking application built with Go and vanilla JavaScript.

### Features

- User authentication (email/password and social logins)
- Expense management (CRUD operations)
- Category-based expense organization
- Spending limits and notifications
- Date range filtering
- Responsive design

### Tech Stack

#### Backend
- Go (Gin framework)
- MongoDB
- JWT authentication
- OAuth2 for social logins

#### Frontend
- Vanilla JavaScript
- HTML5
- CSS3 (with CSS variables for theming)
- Font Awesome icons
- Chart.js for visualizations

### Project Structure

```
expense-tracker/
├── backend-vanilla/       # Go backend
│   ├── config/           # Database configuration
│   ├── controllers/      # Request handlers
│   ├── middleware/       # Custom middleware
│   ├── models/          # Data models
│   └── main.go          # Entry point
└── frontend-vanilla/     # JavaScript frontend
    ├── css/             # Stylesheets
    ├── js/              # JavaScript modules
    └── index.html       # Main HTML file
```

### Setup Instructions

1. Clone the repository:
```bash
git clone https://github.com/your-username/Golang.git
cd Golang/expense-tracker
```

2. Backend Setup:
```bash
cd backend-vanilla
# Copy example env file and update with your values
cp .env.example .env
# Install dependencies
go mod tidy
# Run the server
go run main.go
```

3. Frontend Setup:
```bash
cd frontend-vanilla
# Serve with any static file server
# Example using Python
python -m http.server 5500
```

4. Open http://localhost:5500 in your browser

### Environment Variables

Create a `.env` file in the backend-vanilla directory with the following variables:

```env
MONGODB_URI=mongodb://localhost:27017
DB_NAME=expense_tracker
JWT_SECRET=your-secret-key

# Social OAuth Credentials
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
TWITTER_CLIENT_ID=your-twitter-client-id
TWITTER_CLIENT_SECRET=your-twitter-client-secret
```

### API Endpoints

#### Authentication
- POST `/api/auth/register` - Register new user
- POST `/api/auth/login` - Login user
- GET `/api/auth/:provider/callback` - Social login callback

#### Expenses
- GET `/api/expenses` - Get all expenses
- POST `/api/expenses` - Create new expense
- PUT `/api/expenses/:id` - Update expense
- DELETE `/api/expenses/:id` - Delete expense

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### License

This project is licensed under the MIT License - see the LICENSE file for details.
