# WealthList API

![Go](https://img.shields.io/badge/Go-1.23.4-blue) ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-blue) ![Docker](https://img.shields.io/badge/Docker-✓-blue) ![Swagger](https://img.shields.io/badge/Swagger-✓-green)

## 📌 Description
The WealthList API is a RESTful API for storing, searching, and managing millionaire data. The project is implemented on **Go (Gen)**, uses **PostgreSQL** as a database, and is documented using **Swagger**.

#### Rich people will be called  *"millionaires"* for short.

## 🚀 Functionality
- **Adding** and **removing** millionaires
- **Search** filtered by name and country
- **Upload and receive** photos of millionaires
- **Documented API** via Swagger

## 📂 Project structure
```bash
millionairelist/
├── cmd/               # Initalize app
├── config/            # Configurations and environment variables
├── internal/
│   ├── handler/       # HTTP handlers
│   ├── logger/        # Custom log directory
│   ├── service/       # Business process logic
│   ├── repository/    # Working with the database
│   ├── router/        # Endpoints
│   └── models/        # Definitions of data structures
├── main.go            # Main application launch
├── migrations/        # SQL scripts for DATABASE migrations
├── docs/              # Swagger-documentation
├── Dockerfile         # Instructions for container assembly
├── docker-compose.yml # Docker Compose to launch the service
├── .env.example       # Environment variables template
└── README.md          # This file
```

## 🛠️ Installation and launch
### 🔹 1. Cloning a repository
```sh
git clone https://github.com/azhaxyly/wealthlist.git
cd wealthlist
```

### 🔹 2. Setting up the environment
Copy `.env.example` to `.env` and adjust the variables:
```sh
cp .env.example .env
```

Example `.env`:
```env
DB_HOST=db
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=MILLIONAIRE
DB_PORT=5432
```

### 🔹 3. Launching in Docker
```sh
docker-compose up --build
```

## 📖 API Documentation
The Swagger UI is available at:
```
http://your-host/swagger/index.html
```

Examples of API requests:
- `GET /millionaires` — Get a list of millionaires
- `POST /millionaires` — Add a millionaire
- `DELETE /millionaires/{id}` — Delete a millionaire
- `GET /millionaires/search?lastName=Jobs&country=USA` — Find by filter
- `POST /millionaires/{id}/photo` — Upload a photo
- `GET /millionaires/photos/{imageName}` — Get a photo

## 📦 Development
### 🔹 Local launch without Docker
1. Install Go and PostgreSQL.
2. Create a `MILLIONAIRE' database.
3. Configure the `.env'.
4. Start the server:
```sh
go run main.go
```

## 📜 License
MIT License © 2025

