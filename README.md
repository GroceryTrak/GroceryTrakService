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
DB_HOST=localhost
DB_USER=admin
DB_PASSWORD=adminpassword
DB_DATABASE=grocerytrak

REDIS_HOST=localhost
REDIS_PASS=adminpassword

JWT_SECRET=8dd256ba6e1462d6e3439e51794cd5746455cbd1340af5eb15363181e7edc73a

ENV=development
FRONTEND_DOMAIN=
```

## **API Endpoints**

#### Without Bearer Token ####

| Method | Endpoint             | Description                |
|--------|----------------------|----------------------------|
| POST   | `/auth/register`     | Sign up user               |
| POST   | `/auth/login`        | Sign in user               |

| Method | Endpoint                       | Description                               |
|--------|--------------------------------|-------------------------------------------|
| GET    | `/item/{id}`                   | Fetch a item by ID                        |
| POST   | `/item`                        | Create a new item                         |
| PUT    | `/item/{id}`                   | Update an existing item                   |
| DELETE | `/item/{id}`                   | Delete a item                             |
| GET    | `/item/search?q=&ingredients=1`| Search an item by name, ingredients by ID |

| Method | Endpoint             | Description                |
|--------|----------------------|----------------------------|
| GET    | `/recipe/{id}`       | Fetch a recipe by ID       |
| POST   | `/recipe`            | Create a new recipe        |
| PUT    | `/recipe/{id}`       | Update an existing recipe  |
| DELETE | `/recipe/{id}`       | Delete a recipe            |
| GET    | `/recipe/search?q=`  | Search a recipe by name    |

#### With Bearer Token ####

| Method | Endpoint             | Description                |
|--------|----------------------|----------------------------|
| GET    | `/user_item/{id}`    | Fetch user item by ID      |
| POST   | `/user_item`         | Create new user item       |
| PUT    | `/user_item/{id}`    | Update existing user item  |
| DELETE | `/user_item/{id}`    | Delete user item           |

## **Contributing**
Pull requests are welcome! For major changes, please open an issue first to discuss the proposed changes.

## **License**
This project is licensed under the MIT License.

