# Task Management Application

A simple application to manage tasks, built using **Hexagonal Architecture**. The app leverages **Redis** as its primary storage for tasks, adhering to the principles of clean architecture for better maintainability and scalability.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
## Features

- **Add Tasks**: Create a new task with a name and completion status.
- **Retrieve Tasks**: Fetch all stored tasks.
- **Mark Task as Done**: Update the status of a task to mark it as done.
- **Delete Tasks**: Remove tasks from the storage.

## Architecture

This application is built following the **Hexagonal Architecture (Ports and Adapters)** pattern:

- **Domain Layer**: Contains business logic and the `Task` entity.
- **Application Layer**: Provides use cases for task operations (e.g., adding, updating, retrieving tasks).
- **Infrastructure Layer**: Contains Redis-specific implementations and APIs for external communication.
- **Adapters**:
    - **Input Adapters**: REST API handlers for user interactions.
    - **Output Adapters**: Redis adapter for task storage.

### Hexagonal Diagram

```plaintext
+---------------------------------------+
|                                       |
|            Application Logic          |
|                                       |
|      +---------------------------+    |
|      |       Domain Layer        |    |
|      +---------------------------+    |
|                                       |
|      Adapters            Adapters     |
|      (Input)              (Output)    |
|   +--------------+     +-----------+  |
|   | REST Handler |<--->|   Redis    |  |
|   +--------------+     +-----------+  |
|                                       |
+---------------------------------------+
```

### installation

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/Ehsan-Eghbali/hexagonal.git
   cd hexagonal 
   ```
2. **Set Up Redis: Ensure Redis is installed and running locally  or use Docker to run PostgreSQL.**

    Option 1: Run Redis locally:
    ```bash
    docker run --name redis-service -p 6379:6379 redis:latest
    ```
   **Set Up PostgreSQL: Ensure PostgreSQL is installed and running locally or use Docker to run PostgreSQL.**
   Option 1: Run Redis locally:
   ```bash
    docker run --name postgres-service -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -e POSTGRES_DB=hexa_db -p 5432:5432 postgres:latest
    ```
3. **Install Dependencies:**
   ```bash
    go mod tidy
    ```
4. **Run the Application:**
   run application from deployment/makeFile

## Usage

1. Start the application.
2. Use a REST client like Postman, cURL, or any API tool to interact with the application.

### Base URL
`http://localhost:8080`

## Example Requests

- **POST** `/tasks` - Add a new task:
    - **Body Example**:
      ```json
      { "title": "Task title", "description": "Task description" }
      ```

- **GET** `/tasks` - Retrieve all tasks:
    - **No request body** required.

- **PATCH** `/tasks/:id` - Update a task:
    - **Body Example**:
      ```json
      { "title": "Updated title", "description": "Updated description" }
      ```

- **DELETE** `/tasks/:id` - Delete a task:
    - **No request body** required.
    - **Path Parameter**: `id` (Task ID).

## API Endpoints

### Endpoints Overview

| Method | Endpoint          | Description                | Body Example                                |
|--------|-------------------|----------------------------|--------------------------------------------|
| POST   | `/tasks`          | Add a new task             | `{ "title": "Task title", "description": "Task description" }`    |
| GET    | `/tasks`          | Retrieve all tasks         | _N/A_                                      |
| PATCH  | `/tasks/:id`      | Update a task              | `{ "title": "Updated title", "description": "Updated description" }` |
| DELETE | `/tasks/:id`      | Delete a task              | _N/A_ (Path Parameter: `id`)              |
