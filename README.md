# ms-authz

A production-grade **Authorization Microservice** written in **Golang**, designed to validate **JWT tokens**, handle **RBAC (Role-Based Access Control)** checks, and support **token blacklisting** across distributed systems. Built with **Fiber**, **GORM**, **RabbitMQ**, and **PostgreSQL**.

---

## 🚀 Features

* ✅ Stateless **JWT token** validation
* ✅ RBAC: roles ↔ permissions with many-to-many mappings
* ✅ JWT **blacklist caching** (in-memory, sync.Map based)
* ✅ **RabbitMQ-based** token blacklist and RBAC cache sync
* ✅ Clean Architecture with Unit of Work, Repositories, and Domain Models
* ✅ Swagger/OpenAPI 3.0 documentation via `swag`
* ✅ Auto migration of models on startup

---

## 🧱 Tech Stack

* **Golang 1.24+**
* **Fiber** (web framework)
* **GORM** (ORM)
* **RabbitMQ** (message broker)
* **PostgreSQL** (database)
* **Swagger (swaggo/fiber-swagger)** (API docs)
* **sync.Map** (blacklist cache)

---

## 📁 Project Structure (Simplified)

```
ms-authz/
├── cmd/main.go                 # Application entry point
├── internal/
│   ├── domain/
│   │   ├── model/              # GORM models: Role, Permission, User, etc.
│   │   └── repository/         # Interfaces for repositories
│   ├── infrastructure/
│   │   ├── db/                 # GORM-based repository implementations
│   │   ├── cache/              # In-memory token blacklist
│   │   └── mq/                 # RabbitMQ producer/consumer
│   ├── service/                # AuthService, RBACService
│   └── handler/                # Fiber HTTP handlers (RBAC + Auth)
├── pkg/jwtutil/               # Token parsing and public key management
├── keys/public/               # Public keys for token verification
├── docs/                      # Auto-generated Swagger files
└── go.mod / go.sum
```

---

## 📦 Environment Variables

| Variable         | Description                                        |
| ---------------- | -------------------------------------------------- |
| `PORT`           | Server port (default: 8000)                        |
| `DB_DSN`         | PostgreSQL DSN (e.g. `host=localhost...`)          |
| `RABBITMQ_URL`   | RabbitMQ URL (e.g. `amqp://guest:guest@...`)       |
| `PUBLIC_KEY_DIR` | Path to public key PEM files (e.g. `/keys/public`) |

---

## 🧪 Running Locally

```bash
git clone https://github.com/your-org/ms-authz.git
cd ms-authz

# Generate Swagger docs
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/main.go -o docs

# Run via Docker Compose
docker compose  up --build -d
```

Access Swagger UI at: [http://localhost:8000/swagger/index.html](http://localhost:8000/swagger/index.html)

---

## 🔐 JWT Authorization Flow

* `GET /authorize` endpoint
* Extracts `Authorization: Bearer <token>`
    * Optionally checks:

    * JWT signature
    * Blacklist presence
    * RBAC permission match (`check_rbac=true&privilege=DELETE_USER`)

    ---

    ## 🔄 RabbitMQ Events

    | Exchange             | Event Key           | Purpose                          |
    | -------------------- | ------------------- | -------------------------------- |
    | `auth.tokens.fanout` | `TOKEN_BLACKLISTED` | Add token to blacklist cache     |
    | `rbac.update.fanout` | `RBAC_CACHE_RELOAD` | Reload local RBAC permission map |

    ---

    ## 📚 API Endpoints (Summary)

    ### 🔒 Authorization

    | Method | Endpoint     | Description                  |
    | ------ | ------------ | ---------------------------- |
    | GET    | `/authorize` | JWT + Blacklist + RBAC check |

    ### 🧑‍💼 Roles

    | Method | Endpoint                           | Description            |
    | ------ | ---------------------------------- | ---------------------- |
    | GET    | `/roles`                           | Get all roles          |
    | POST   | `/roles`                           | Create new role        |
    | PUT    | `/roles/{id}`                      | Update existing role   |
    | DELETE | `/roles/{id}`                      | Delete role by ID      |
    | GET    | `/roles/{id}/permissions`          | Get role's permissions |
    | POST   | `/roles/{id}/permissions/{permID}` | Assign permission      |
    | DELETE | `/roles/{id}/permissions/{permID}` | Remove permission      |

    ### 🛡️ Permissions

    | Method | Endpoint                  | Description                      |
    | ------ | ------------------------- | -------------------------------- |
    | GET    | `/permissions`            | Get all permissions              |
    | POST   | `/permissions`            | Create new permission            |
    | PUT    | `/permissions/{id}`       | Update permission                |
    | DELETE | `/permissions/{id}`       | Delete permission by ID          |
    | GET    | `/permissions/{id}/roles` | Get roles assigned to permission |

    ### 🔁 Expanded Queries

    | Method | Endpoint                  | Description                    |
    | ------ | ------------------------- | ------------------------------ |
    | GET    | `/roles-with-permissions` | All roles with permission list |
    | GET    | `/permissions-with-roles` | All permissions with roles     |

    ---

    ## 🛠 Example Payloads

    ### Create Role

    ```json
    {
    "name": "admin",
    "description": "System administrator"
    }
    ```

    ### Create Permission

    ```json
    {
    "name": "DELETE_USER",
    "description": "Can delete users"
    }
    ```

    ---

    ## 🧑‍💻 Developer Notes

    * RBAC cache is stored in `sync.Map` and updated via MQ broadcast
    * Public keys for JWT must be stored in PEM format: `/keys/public/<kid>.pem`
        * Repositories are injected via `UnitOfWork`
        * Swagger annotations are located in handler files (e.g. `AuthorizeHandler`, `RBACAdminHandler`)