# BetterSpace API

A comprehensive REST API for office booking and coworking space management built with Go and Echo framework.

## Features

- **User Authentication**: JWT-based authentication system
- **Office Management**: Complete CRUD operations for office spaces
- **Facility Management**: Manage office amenities and facilities
- **Booking System**: Transaction management for office reservations
- **Review System**: User reviews and ratings for office spaces
- **Image Management**: Upload and manage office images via Cloudinary
- **Geolocation Support**: Location-based office search and validation

## Tech Stack

- **Language**: Go 1.18+
- **Framework**: Echo v4
- **Database**: MySQL with GORM ORM
- **Authentication**: JWT (golang-jwt/jwt)
- **Image Storage**: Cloudinary
- **Validation**: Go Playground Validator
- **Password Security**: Go Password Validator
- **Containerization**: Docker

## Project Structure

```
betterspace-api/
├── app/
│   ├── middlewares/          # Authentication, logging, CORS
│   └── routes/               # API route definitions
├── businesses/               # Business logic layer
│   ├── users/               # User domain logic
│   ├── offices/             # Office domain logic
│   ├── facilities/          # Facility domain logic
│   ├── transactions/        # Booking transaction logic
│   └── review/              # Review system logic
├── controllers/             # HTTP handlers and request/response
├── drivers/                 # Data access layer
│   └── mysql/              # MySQL repository implementations
├── utils/                   # Utility functions
├── main.go                 # Application entry point
├── go.mod                  # Go modules
└── Dockerfile              # Docker configuration
```

## Prerequisites

- Go 1.18 or higher
- MySQL 5.7+ or 8.0+
- Docker (optional)

## Installation

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/your-username/betterspace-api.git
   cd betterspace-api
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   Create a `.env` file in the root directory:
   ```env
   # Database Configuration
   DB_USERNAME=your_db_username
   DB_PASSWORD=your_db_password
   DB_HOST=localhost
   DB_PORT=3306
   DB_NAME=betterspace_db
   
   # JWT Configuration
   JWT_SECRET_KEY=your_jwt_secret_key
   
   # Server Configuration
   PORT=3000
   
   # Cloudinary Configuration (for image uploads)
   CLOUDINARY_CLOUD_NAME=your_cloud_name
   CLOUDINARY_API_KEY=your_api_key
   CLOUDINARY_API_SECRET=your_api_secret
   ```

4. **Set up MySQL database**
   ```sql
   CREATE DATABASE betterspace_db;
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

### Docker Deployment

1. **Build and run with Docker**
   ```bash
   docker build -t betterspace-api .
   docker run -p 3000:3000 --env-file .env betterspace-api
   ```

## API Endpoints

### Authentication
- `POST /auth/register` - User registration
- `POST /auth/login` - User login

### Users
- `GET /users/profile` - Get user profile
- `PUT /users/profile` - Update user profile

### Offices
- `GET /offices` - List all offices
- `GET /offices/:id` - Get office details
- `POST /offices` - Create new office
- `PUT /offices/:id` - Update office
- `DELETE /offices/:id` - Delete office

### Facilities
- `GET /facilities` - List all facilities
- `POST /facilities` - Create facility
- `PUT /facilities/:id` - Update facility
- `DELETE /facilities/:id` - Delete facility

### Office Facilities
- `POST /office-facilities` - Assign facility to office
- `DELETE /office-facilities/:id` - Remove facility from office

### Transactions
- `GET /transactions` - List user transactions
- `POST /transactions` - Create booking transaction
- `PUT /transactions/:id` - Update transaction status

### Reviews
- `GET /reviews/:office_id` - Get office reviews
- `POST /reviews` - Create review
- `PUT /reviews/:id` - Update review
- `DELETE /reviews/:id` - Delete review

### Office Images
- `POST /office-images` - Upload office image
- `DELETE /office-images/:id` - Delete office image

## Database Schema

The application automatically migrates the database schema on startup. Key entities include:

- **Users**: User authentication and profile data
- **Offices**: Office space information and details
- **Facilities**: Available amenities (WiFi, AC, Parking, etc.)
- **Office_Facilities**: Many-to-many relationship between offices and facilities
- **Transactions**: Booking records and payment information
- **Reviews**: User reviews and ratings for offices
- **Office_Images**: Image metadata for office photos

## Development

### Code Structure

The project follows Clean Architecture principles:

- **Domain Layer** (`businesses/`): Contains business logic and use cases
- **Infrastructure Layer** (`drivers/`): Database interactions and external services
- **Interface Layer** (`controllers/`): HTTP handlers and API contracts
- **Application Layer** (`main.go`): Dependency injection and app initialization

### Adding New Features

1. Define domain models in `businesses/[domain]/domain.go`
2. Implement use cases in `businesses/[domain]/usecase.go`
3. Create database models in `drivers/mysql/[domain]/record.go`
4. Implement repository in `drivers/mysql/[domain]/mysql.go`
5. Add HTTP handlers in `controllers/[domain]/http.go`
6. Define request/response models in `controllers/[domain]/request|response/`
7. Register routes in `app/routes/route.go`

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/new-feature`)
3. Commit your changes (`git commit -am 'Add new feature'`)
4. Push to the branch (`git push origin feature/new-feature`)
5. Create a Pull Request

## Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## CI/CD

The project includes GitHub Actions workflow for automated:
- Code quality checks
- Testing
- Building and deployment

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support and questions, please create an issue in the GitHub repository.