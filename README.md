# **Uriel**

A boilerplate for creating simple extensible backend services in Go

This project uses an **Ent** (Entity Framework) for database interactions, **net/http** for HTTP routing, and provides functionality for user authentication, background jobs, and service-layer abstractions.

---

## **Features**
- User authentication and token-based authorization.
- Database ORM support with **Ent** for schema management and query building.
- Background jobs and scheduled tasks.
- Simple separation of HTTP handlers, services, and repositories.
- Environment-specific configuration support.
- Bruno `.bru` API testing integration.

---

## **Project Structure**

The project directory is divided as follows:

### **Root Directory**
- `main.go`: The entry point for your application.
- `go.mod` / `go.sum`: Go module dependencies for the project.
- `uriel/`: A folder containing `.bru` files to define API testing collections using [Bruno](https://usebruno.com/), a REST API testing tool.

---

### **`cmd/`**
This directory contains various high-level commands and scripts for the application:
- **`httpd.go`**: Starts the HTTP server for the application.
- **`job.go`**: Entry point for running one time or background jobs.
- **`migrate.go`**: Manages database migrations.
- **`root.go`**: Main command executed by the CLI tool to initialize the app.
- **`sing.go`**: Example job subcommand.

---

### **`internal/`**
The internal directory is used for the project's core functionality. It is divided into logically separated packages, ensuring encapsulation.

#### **`config/`**
- **`config.go`**: Maintains global application configurations (e.g., environment variables, app secrets, etc.).

#### **`dto/`**
(Data Transfer Object)
- **`auth.go`**: DTOs for user authentication endpoints.
- **`viewer.go`**: DTOs for Viewer-related functionalities.

#### **`ent/`**
(Ent Framework for Database ORM)
- **Schema-related Definitions:**
    - **`schema/`**: Contains schema definitions like `user.go` for generating the database models and migrations.
- **Core Ent Components:**
    - Files like `client.go`, `mutation.go`, `tx.go`, and `user.go` provide core access to database queries and transactions.
- **Generated Code:**
    - The `ent.go`, `runtime/`, and related files contain Ent Framework-generated boilerplate code.

#### **`httpd/`**
(HTTP Handlers)
This folder contains all the HTTP endpoint handlers and their supporting logic:
- **`auth.go`**: Authentication-related request handlers (e.g., login, logout).
- **`helper.go`**: Utility functions for HTTP requests.
- **`router.go`**: Configures all app routes, including public and private routes.
- **`viewer.go`**: Viewer-related HTTP endpoints.

#### **`job/`**
Background jobs and task definitions:
- **`sing.go`**: Represents an example one time job

#### **`repository/`**
This layer interacts with the database and abstracts raw database queries:
- **`db.go`**: Initializes and manages database connections.
- **`repository.go`**: Contains the repository pattern implementation.

#### **`service/`**
This layer contains the business logic of the application:
- **`service.go`**: Base service implementation.
- **`user.go`**: Business logic for user/auth-related functionality.
- **`viewer.go`**: Viewer-specific service logic.

---

### **`uriel/`**
This folder contains API testing collections and environment files designed for the Bruno API testing tool:
- **`Login.bru`**, **`Viewer.bru`**: Define API requests for testing login and viewer-related endpoints.
- **`bruno.json`**: Configuration file for Bruno.
- **`collection.bru`**: A collection of the API endpoints used in testing.
- **`environments/local.bru`**: Environment-specific configuration for Bruno tests (e.g., local development).

---

## **Getting Started**

### **System Requirements**
- Go 1.23 or higher.
- A running MySQL database (configurable in `.config.yaml` or your ENV).

---

### **Setup Instructions**

1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/uriel.git
   cd uriel
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Generate Ent code (if you've made schema changes):
   ```bash
   go generate ./internal/ent
   ```

4. Run database migrations:
   ```bash
   go run main.go migrate
   ```

5. Start the HTTP server:
   ```bash
   go run main.go httpd
   ```

6. Test the API with Bruno:
    - Import the `.bru` collection files in the `uriel/` folder into your Bruno workspace.
    - Configure the `local.bru` environment for local testing.

---

## **Development**

### **Building the Application**
```bash
go build -o uriel main.go
```

### **Running Background Jobs**
```bash
go run main.go job <jobname>
```

### **Testing**
No testing, please

---
