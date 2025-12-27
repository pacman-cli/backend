# Portfolio Backend

Spring Boot backend for Personal Portfolio Blog.

## Tech Stack
- **Java 17**
- **Spring Boot 3.4.1**
- **PostgreSQL** (Database)
- **Flyway** (Migrations)
- **Spring Security + JWT** (Authentication)
- **Lombok**

## Setup

1.  **Database**: Ensure PostgreSQL is running on `localhost:5432`.
    - Create database: `portfolio_db`
    - Username: `postgres`, Password: `password` (Update in `application.yml`)
2.  **Run**:
    ```bash
    ./mvnw spring-boot:run
    ```

## API Documentation

### Public Endpoints
- `GET /api/blogs`: List blogs (pagination: `page`, `size`, `sort`). Query params: `search`, `publicOnly=true`.
- `GET /api/blogs/{slug}`: Get single blog by slug.
- `GET /uploads/{filename}`: Serve static images.

### Admin Endpoints (Requires `Authorization: Bearer <token>`)
- `POST /api/auth/register`: Register admin (Disable in prod).
- `POST /api/auth/login`: Login to get JWT token.
- `POST /api/blogs`: Create blog.
- `PUT /api/blogs/{id}`: Update blog.
- `DELETE /api/blogs/{id}`: Delete blog.
- `POST /api/upload`: Upload image (returns URL).

## Folder Structure
- `controller`: API Endpoints
- `service`: Business Logic
- `repository`: Database Access
- `model`: JPA Entities
- `dto`: Data Transfer Objects
- `security`: JWT Config
- `config`: App Config
- `db/migration`: SQL Scripts
