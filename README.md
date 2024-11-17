# Tender Bridge Project

The Tender Bridge Project is a backend application designed for managing tenders. It uses Golang for backend logic, PostgreSQL for database management, Redis for caching, and Docker for containerization.

---

## Setup and Usage

### 1. Start the Database

To start the database container, run:

```bash
make run_db
```

This will spin up a PostgreSQL container and automatically run migrations.

---

### 2. Start the Application

To start the entire application (database, backend, etc.), run:

```bash
make run
```

This command:
- Starts the database (if not already running)
- Builds the application Docker image
- Starts the application container

The backend will be available at:  
**[http://localhost:8888](http://localhost:8888)**

---

## Commands

### `make run_db`
- Starts the PostgreSQL database container.
- Automatically applies migrations.

### `make run`
- Starts the database, builds the application, and starts all services.

---

## Troubleshooting

1. **Database Connection Issues**  
   Ensure Docker is running and the database container is up using `docker ps`.

2. **Migrations Fail**  
   Verify the database environment variables in `docker-compose.yml`.

3. **Redis Not Working**  
   Check if the Redis service is running and verify the environment variables.

---

## Environment Variables

Environment variables are configured in the `docker-compose.yml` file. Key variables include:

| Variable                   | Default Value           | Description                       |
|----------------------------|-------------------------|-----------------------------------|
| `HOST`                     | `localhost`            | Application host.                |
| `PORT`                     | `8888`                 | Application port.                |
| `POSTGRES_HOST`            | `db`                   | PostgreSQL host.                 |
| `POSTGRES_PORT`            | `5432`                 | PostgreSQL port.                 |
| `POSTGRES_DB`              | `tender_bridge_db`     | PostgreSQL database name.        |
| `POSTGRES_USER`            | `postgres`             | PostgreSQL user.                 |
| `POSTGRES_PASSWORD`        | `password`             | PostgreSQL password.             |
| `REDIS_HOST`               | `redis`                | Redis host.                      |
| `REDIS_PORT`               | `6379`                 | Redis port.                      |
| `JWT_SECRET`               | `tender-bridge-forever` | JWT secret key.                  |

---

## All passed the test
![image_2024-11-17_13-28-44](https://github.com/user-attachments/assets/59f1c6ec-ca6d-4b54-903c-635dd64acd53)
