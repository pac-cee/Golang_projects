# URL Shortener

A web application that shortens long URLs into easy-to-share links, demonstrating web development, database operations, and template rendering in Go.

## Concepts Covered

- Web server creation
- Database operations (SQLite)
- HTML templates
- Static file serving
- URL routing
- Form handling
- Error handling
- Random string generation
- CSS styling

## Features

- Shorten long URLs to unique codes
- Track click counts for each shortened URL
- Display list of recently shortened URLs
- Responsive web interface
- Persistent storage using SQLite
- Automatic redirection to original URLs

## Prerequisites

1. Go installed on your system
2. SQLite installed
3. Install required Go packages:
   ```bash
   go get github.com/mattn/go-sqlite3
   ```

## How to Run

1. Navigate to the project directory:
   ```bash
   cd level-2/04-url-shortener
   ```

2. Run the program:
   ```bash
   go run main.go
   ```

3. Open your browser and visit:
   ```
   http://localhost:8080
   ```

## Project Structure

```
04-url-shortener/
├── main.go              # Main application file
├── urls.db              # SQLite database (created on first run)
├── README.md            # Project documentation
├── static/
│   └── style.css        # CSS styles
└── templates/
    └── index.html       # HTML template
```

## Database Schema

```sql
CREATE TABLE urls (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    long_url TEXT NOT NULL,
    short_code TEXT UNIQUE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    clicks INTEGER DEFAULT 0
);
```

## API Endpoints

- `GET /` - Home page with URL submission form
- `POST /shorten` - Create a new short URL
- `GET /{shortCode}` - Redirect to original URL

## Learning Objectives

- Building web servers in Go
- Working with SQL databases
- HTML template parsing and execution
- Static file serving
- Form handling and validation
- URL routing and redirection
- Error handling in web applications
- Using prepared statements
- Basic frontend development

## Next Steps

To extend this project, you could:
1. Add user authentication
2. Implement custom short codes
3. Add URL validation
4. Create an API endpoint
5. Add URL expiration
6. Implement rate limiting
7. Add QR code generation
8. Create detailed analytics
9. Add unit tests
10. Deploy to a cloud platform
