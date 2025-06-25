# Go API with PostgreSQL Backend
This is health-portal project. it has apis written in golang along with support of all CRUD operations. It usses postgre sql as database. It as has JWT Authentication.


## Project Structure
```
├── main.go
├── go.mod
├── .env
├── config/
│   └── config.go
├── controllers/
│   ├── auth.go
│   └── patient.go
├── models/
│   └── models.go
├── repositories/
│   └── patient.go
├── routes/
│   └── routes.go
├── middlewares/
│   └── auth.go
├── utils/
│   └── utils.go
└── db/
    └── database.go
```

## Features
```
- RESTful APIs in Go
- JWT Authentication
- Full CRUD for Patient Data
- PostgreSQL Database Integration
- Swagger (via static `index.html`) for API Testing
- Clean MVC + Repository pattern
```
---

## Install dependencies and run:
```
go mod download
go run main.go
```

---

## To Test APIs - 
```
- Open - index.html and update api endpoint to localhost:8080
- Dependencies
- Gin — Web framework for Go
- pq — PostgreSQL driver
- godotenv — .env loader
- gorm — ORM for Go
- jwt-go — JWT support

```



## 🧪 How to Test APIs
```
Open the index.html file in your browser.

Replace the API_BASE_URL with:
http://localhost:8080

```

## Use available endpoints to test: Example:
```
POST   /signup
POST   /login
GET    /patients
POST   /patients
PUT    /patients/:id
DELETE /patients/:id

```
