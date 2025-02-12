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

### **2️⃣ Using Air (Hot Reloading - Recommended for Development)**

Install Air if you haven't already:
```sh
go install github.com/air-verse/air@latest
```

Run the service with Air:
```sh
air
```

### **3️⃣ Not using Air, running with `go run` (Manual Restart Required)**
```sh
go run cmd/*.go
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
| Method | Endpoint              | Description                 |
|--------|----------------------|-----------------------------|
| GET    | `/recipes/{id}`       | Fetch a recipe by ID       |
| POST   | `/recipes`            | Create a new recipe        |
| PUT    | `/recipes/{id}`       | Update an existing recipe  |
| DELETE | `/recipes/{id}`       | Delete a recipe            |

## **Contributing**
Pull requests are welcome! For major changes, please open an issue first to discuss the proposed changes.

## **License**
This project is licensed under the MIT License.

