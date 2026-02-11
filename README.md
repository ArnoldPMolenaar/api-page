# Page API

[![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go)](https://go.dev/)
[![Fiber](https://img.shields.io/badge/Fiber-v2-00ACD7?logo=fiber)](https://github.com/gofiber/fiber)
[![GORM](https://img.shields.io/badge/GORM-ORM-7E57C2?logo=sqlite)](https://gorm.io/)

## 📚 Overview
A microservice responsible for managing Pages, Menus, Modules, and Versions for applications. It exposes a clean set of public endpoints for published content and machine-protected private endpoints for authoring and management.

Key features:
- Page composition with Partials, Rows, and Columns per locale
- Menu and Module management tied to a Version
- Versioning with publish/restore flows per App
- Valkey (Redis-compatible) cache for performance
- Built on Fiber (HTTP), GORM (DB), and Go

## 🧩 Architecture
- Fiber app bootstrapped in `main.go`
- Database and cache connections initialized at startup
- Public and Private routes defined under `src/routes`
- Controllers orchestrate validation, services, and DTO responses
- Uses `api-utils` for middleware, errors, routes, and server lifecycle

## 🐳 Run with Docker
This service ships with a compose setup for development and production. Both use `network_mode: host` for simple local networking and include a Valkey cache.

### 🛠️ Development
Prerequisites:
- Docker and Docker Compose
- A `.env` file at repo root (see Environment section)

Steps:
1) Start Valkey and the API in dev mode with live reload (Air):

```zsh
docker compose up dev
```

Notes:
- The dev service mounts the repo into the container and runs Air for hot reload.
- API listens on port 5000.

### 🚀 Production (local build)
Build and run the production image (multi-stage, static binary):

```zsh
docker compose up prod valkey --build
```

Notes:
- The prod image is based on scratch and runs the compiled `/api` binary.
- API listens on port 5000.
- Depends on `valkey` service is ignored in production so you can wire it to an external cache if desired.

## 🔐 Environment
Provide a `.env` at the project root. Typical variables include:
- STAGE_STATUS=dev|prod (controls graceful shutdown)
- SERVER_PORT=5000
- DATABASE_* (driver, DSN, etc.)
- VALKEY_* (host, port)
- Any app-specific settings referenced by services

Tip: the production Dockerfile copies `.env` into the image; keep secrets scoped to your environment.

## 🌐 API Endpoints
Base path: `/v1`

### 🔓 Public
For consuming published content.

- GET `/v1/versions/published`
  - Query: `app=<appName>`
  - Returns the published Version of an App.

- GET `/v1/versions/:id/menus`
  - Query: `locale=<locale>`
  - Returns published Menus for a published Version. Requires the Version to be published.

- GET `/v1/pages/:menuItemId/:locale/published`
  - Returns the published Page for a Menu Item in a given locale.

### 🛡️ Private (Machine Protected)
All endpoints require machine authentication via `api-utils` middleware.

- Apps
  - POST `/v1/apps/`
    - Create a new App.

- Versions
  - GET `/v1/versions/`
    - Paginated list of Versions.
  - POST `/v1/versions/`
    - Create a Version for an App.
  - GET `/v1/versions/lookup`
    - Query: `app=<appName>` (required), `name=<partialName>` (optional)
    - Returns Version lookup list for an App.
  - GET `/v1/versions/name/available`
    - Query: `app=<appName>`, `name=<versionName>`, `ignore=<nameToIgnore>` (optional)
    - Checks if a Version name is available.
  - GET `/v1/versions/:id`
    - Get a Version by ID.
  - PUT `/v1/versions/:id`
    - Update a Version.
  - DELETE `/v1/versions/:id`
    - Soft-delete a Version.
  - PUT `/v1/versions/:id/publish`
    - Publish a Version.
  - PUT `/v1/versions/:id/restore`
    - Restore a previously deleted Version.

- Menus
  - GET `/v1/menus/`
    - Paginated list of Menus.
  - POST `/v1/menus/`
    - Create a Menu.
  - GET `/v1/menus/lookup`
    - Returns Menu lookup list.
  - GET `/v1/menus/name/available`
    - Checks if a Menu name is available.
  - GET `/v1/menus/:id`
    - Get a Menu by ID.
  - PUT `/v1/menus/:id`
    - Update a Menu.
  - DELETE `/v1/menus/:id`
    - Soft-delete a Menu.
  - PUT `/v1/menus/:id/restore`
    - Restore a previously deleted Menu.

- Pages
  - GET `/v1/pages/:menuItemId/:locale`
    - Get or create the draft Page for a Menu Item in a given locale.
  - PUT `/v1/pages/:menuItemId/:locale`
    - Update a Page.
  - DELETE `/v1/pages/:menuItemId/:locale`
    - Soft-delete a Page.
  - PUT `/v1/pages/:menuItemId/:locale/restore`
    - Restore a previously deleted Page.
  - POST `/v1/pages/:menuItemId/:locale/partials`
    - Create a Page Partial.
  - PUT `/v1/pages/:menuItemId/:locale/partials/:id`
    - Update a Page Partial.
  - DELETE `/v1/pages/:menuItemId/:locale/partials/:id`
    - Soft-delete a Page Partial.
  - PUT `/v1/pages/:menuItemId/:locale/partials/:id/restore`
    - Restore a previously deleted Page Partial.

- Modules
  - GET `/v1/modules/`
    - Paginated list of Modules.
  - POST `/v1/modules/`
    - Create a Module.
  - GET `/v1/modules/lookup`
    - Returns Module lookup list.
  - GET `/v1/modules/name/available`
    - Checks if a Module name is available.
  - GET `/v1/modules/:id`
    - Get a Module by ID.
  - PUT `/v1/modules/:id`
    - Update a Module.
  - DELETE `/v1/modules/:id`
    - Soft-delete a Module.
  - PUT `/v1/modules/:id/restore`
    - Restore a previously deleted Module.

## 🧪 Health and Errors
- 404 route is registered via `api-utils` to handle unknown endpoints.
- Consistent error responses through `api-utils/errors`.

## 📦 Tech Stack
- Go, Fiber, GORM, Valkey
- Dockerized development and production workflows
- Shared utilities via `github.com/ArnoldPMolenaar/api-utils`

---

## 🤝 Contributing
We welcome contributions! Please fork the repository and submit a pull request.

## 📝 License
This project is licensed under the MIT License.

## 📞 Contact
For any questions or support, please contact [arnold.molenaar@webmi.nl](mailto:arnold.molenaar@webmi.nl).

<hr />

Made with ❤️ by Arnold Molenaar