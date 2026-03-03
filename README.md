# 🏨 Hotel Bookings and Reservation System

**Full-Stack Clean Architecture Project (Go + Next.js)**

A production-ready **hotel booking and reservation platform** built with:

* **Backend:** Golang (Clean Architecture + Domain-Driven Design)
* **Frontend:** Next.js (modern React UI)
* **Architecture:** Modular, scalable, testable, enterprise-grade

This system provides complete room lifecycle management, authentication, and booking-ready infrastructure designed for real-world deployment.

---



# 🧭 Overview

This project demonstrates how to build a **scalable hotel reservation system** using:

✔ Clean Architecture
✔ Domain-Driven Design (DDD)
✔ Repository Pattern
✔ Dependency Injection
✔ Service Layer Testing
✔ Separation of Concerns

The backend exposes business logic via structured services and repositories, while the frontend provides a responsive user interface for managing rooms and bookings. 

---

# 🏗 System Architecture

```
Client (Next.js UI)
        │
        ▼
HTTP API Layer (Handlers / Controllers)
        │
        ▼
Application Layer (Usecases / Services)
        │
        ▼
Domain Layer (Entities / Business Rules)
        │
        ▼
Infrastructure Layer (Database / Repositories)
```

Key Principles:

* Domain logic is independent of frameworks
* Infrastructure is replaceable
* Business rules are testable
* Clear responsibility boundaries

---

# 🧰 Tech Stack

## Backend

* Golang
* Clean Architecture
* Repository Pattern
* bcrypt password hashing
* Unit testing with mocks

## Frontend

* Next.js (App Router or Pages Router)
* React Hooks
* Axios / Fetch API
* Tailwind CSS (optional styling)

## Testing

* Go testing package
* Mock repositories
* Race condition testing

---

# 🚀 Features

## 👤 User Management

* User registration
* Secure password hashing
* Authentication ready
* Role-based access
* Active / inactive user states
* Service layer validation

---

## 🛏 Room Management

* Create room
* Update room
* Delete room
* Change availability status
* Pagination support
* Status filtering
* Validation rules
* Fully unit tested

---

## 🧪 Testing Coverage

* Service logic tests
* Pagination tests
* Status transition tests
* Filtering tests
* Repository mocking
* Concurrency safety
* ~80%+ coverage

---

## 🔐 Security

* bcrypt hashing
* structured error handling
* validation guards
* safe state transitions

---


# ⚙ Backend Setup (Go)

### 1️⃣ Install dependencies

```
cd server
go mod tidy
```

---

### 2️⃣ Run server

```
go run cmd/main.go
```

Server runs on:

```
http://localhost:8080
```

---

# 💻 Frontend Setup (Next.js)

### 1️⃣ Install dependencies

```
cd client
npm install
```

---

### 2️⃣ Run development server

```
npm run dev
```

Frontend runs on:

```
http://localhost:3000
```

---

# 🔑 Environment Variables

## Backend (.env)



---

## Frontend (.env.local)

```
NEXT_PUBLIC_API_URL=http://localhost:8080
```

---

# ▶ Running the Full Application

Start backend:

```
cd server
go run cmd/main.go
```

Start frontend:

```
cd client
npm run dev
```

Open browser:

```
http://localhost:3000
```

---

# 📡 API Overview

## Rooms

### Create Room

```
POST /rooms
```

### Get Room

```
GET /rooms/{id}
```

### Update Room

```
PUT /rooms/{id}
```

### Delete Room

```
DELETE /rooms/{id}
```

### List Rooms (pagination)

```
GET /rooms?page=1&limit=10
```

### Filter by status

```
GET /rooms?status=available
```

---

# 🧪 Testing

Run all tests:

```
go test ./... -v
```

Run room service tests:

```
go test ./pkg/usecase/room -v
```

Race detection:

```
go test -race ./...
```

Coverage:

```
go test -cover ./...
```

---

# 📄 Pagination Logic

Supports:

* first page
* middle page
* last partial page
* overflow handling

Example:

```
GET /rooms?page=2&limit=5
```

---

# 🔄 Status Transitions

Allowed transitions:

```
available → booked
booked → available
```

Invalid transitions return validation errors.

---




# 📜 License

MIT License

---

# 👨‍💻 Author

Hotel Bookings and Reservation System
Clean Architecture Go + Next.js Implementation

---

If you found this project useful, feel free to star the repository ⭐
