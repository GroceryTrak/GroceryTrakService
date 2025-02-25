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
go run main.go
```

Choice 2: Using Air
```
air
```

## **Environment Variables**
Before running, ensure your `.env` file contains the correct values:
Change localhost to name of the container in `docker-compose.yaml` (db and redis) if running app with `docker-compose.yaml`
```ini
DB_HOST=localhost
DB_USER=admin
DB_PASSWORD=adminpassword
DB_DATABASE=grocerytrak

REDIS_HOST=localhost
REDIS_PASS=adminpassword

JWT_SECRET=8dd256ba6e1462d6e3439e51794cd5746455cbd1340af5eb15363181e7edc73a

ENV=development
FRONTEND_DOMAIN=
DETECT_DOMAIN=https://grocerytrak-devs-grocerytrakdetect.hf.space
```

## **API Endpoints**
Please run the app and check `/swagger/index.html`.
When updating API documentation, run `swag init`

## **Contributing**
Pull requests are welcome! For major changes, please open an issue first to discuss the proposed changes.

## **License**
This project is licensed under the MIT License.

