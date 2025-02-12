# GroceryTrakService

## Overview
GroceryTrakService is the backend service for **GroceryTrak**, an application designed to help users efficiently manage grocery lists, track inventory, and optimize shopping experiences. This service is built using **Go**, **PostgreSQL**, and **Redis**, and utilizes **Docker** for containerized deployment.

## Installation

### **Prerequisites**
Ensure you have the following installed:
- [Go](https://go.dev/dl/) (1.20+ recommended)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)
- [Air](https://github.com/cosmtrek/air) (for hot reloading, optional)

### **Clone the Repository**
```sh
git clone https://github.com/yourusername/GroceryTrakService.git
cd GroceryTrakService
```

## **Running the Service**

### **1️⃣ Running with Docker Compose**
Ensure Docker is running, then execute:
```sh
docker compose up --build
```

To stop the containers:
```sh
docker compose down
```

### **2️⃣ Install dependencies**
```sh
go mod tidy
```

Install Air (optional):
```sh
go install github.com/air-verse/air@latest
```

### **3️⃣ Running**
Choice 1: Not using Air
```sh
go run cmd/*.go
```

Choice 2: Using Air
```
air
```

## **Environment Variables**
Before running, ensure your `.env` file contains the correct values:
```ini
DB_HOST="localhost"
DB_USER="admin"
DB_PASSWORD="adminpassword"
DB_DATABASE="grocerytrak"
DB_SSLMODE="disable"
DB_PORT=5432

REDIS_HOST="redis"
REDIS_PASS=""
```

## **API Endpoints**
| Method | Endpoint             | Description                |
|--------|----------------------|----------------------------|
| POST   | `/auth/register`     | Sign up user               |
| POST   | `/auth/login`        | Sign in user               |

| Method | Endpoint             | Description                |
|--------|----------------------|----------------------------|
| GET    | `/item/{id}`         | Fetch a item by ID         |
| POST   | `/item`              | Create a new item          |
| PUT    | `/item/{id}`         | Update an existing item    |
| DELETE | `/item/{id}`         | Delete a item              |

| Method | Endpoint             | Description                |
|--------|----------------------|----------------------------|
| GET    | `/recipe/{id}`       | Fetch a recipe by ID       |
| POST   | `/recipe`            | Create a new recipe        |
| PUT    | `/recipe/{id}`       | Update an existing recipe  |
| DELETE | `/recipe/{id}`       | Delete a recipe            |

## **Contributing**
Pull requests are welcome! For major changes, please open an issue first to discuss the proposed changes.

## **License**
This project is licensed under the MIT License.

