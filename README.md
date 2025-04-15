# GroceryTrakService

## Overview
GroceryTrakService is the backend service for **GroceryTrak**, an application designed to help users efficiently manage grocery lists, track inventory, and optimize shopping experiences. This service is built using **Go**, **PostgreSQL**, and **Redis**, and utilizes **Docker** for containerized deployment.

## Installation

### **Prerequisites**
Ensure you have the following installed:
- [Go](https://go.dev/dl/) (1.20+ recommended)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)
- [Air](https://github.com/cosmtrek/air) (for hot reloading, optional)
```sh
go install github.com/air-verse/air@latest
```

### **Clone the Repository**
```sh
git clone https://github.com/GroceryTrak/GroceryTrakService.git
cd GroceryTrakService
```

## **Running the Service**

### **1️⃣ Install dependencies**
```sh
go mod tidy
```


### **2️⃣ Environment Variables**
Before running, ensure your `.env` file (located at root directory) contains the correct values:
Change **name of the container** in `docker-compose.yaml` (`db` and `redis`) to `localhost` if running app outside of `docker-compose.yaml`. Meaning if running app locally (with `go run main.go` or `air`, and not `docker compose up`), then keep `localhost`. Remember to change password in production.

```ini
DB_HOST=db
DB_USER=admin
DB_PORT=5432
DB_PASSWORD=adminpassword
DB_DATABASE=grocerytrak


REDIS_HOST=redis
REDIST_PORT=6379
REDIS_PASS=adminpassword

JWT_SECRET=8dd256ba6e1462d6e3439e51794cd5746455cbd1340af5eb15363181e7edc73a

SPOONACULAR_API_KEY=0123456789abcdef0123456789abcdef
SPOONACULAR_API_URL=https://api.spoonacular.com
SPOONACULAR_IMG_URL=https://img.spoonacular.com/ingredients_500x500

OPENAI_API_KEY=sk-proj-0123456789abcdef0123456789abcdef

ENV=development
FLUTTER_URL=http://localhost:53459
HUGGINGFACE_URL=https://grocerytrak-devs-grocerytrakdetect.hf.space
```


### **3️⃣ Running**
Choice 1: Using Docker:
```sh
docker compose up --build
```

To stop the containers:
```sh
docker compose down
```

Choice 2: Not using Docker and not using Air
```sh
go run main.go
```

Choice 3: Not using Docker but using Air
```
air
```

## **API Endpoints**
Please run the app and check `/swagger/index.html`.
When updating API documentation, run `swag init`

## **Contributing**
Pull requests are welcome! For major changes, please open an issue first to discuss the proposed changes.

## **License**
This project is licensed under the MIT License.

