# Qisur Challenge - Backend API Documentation

## 📚 Índice

- [Instalación](#instalación)
- [Configuración](#configuración)
- [Ejecución](#ejecución)
- [API REST](#api-rest)
- [WebSockets](#websockets)
- [Autenticación](#autenticación)
- [Arquitectura](#arquitectura)

---

## Instalación

### Requisitos

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 16 (o usar docker-compose)

### Clonar y configurar

```bash
git clone <repo>
cd qisurChallenge

# Instalar dependencias
go mod download

# Crear .env (o usar el existente)
cp .env.example .env
```

---

## Configuración

### Variables de entorno (.env)

```ini
# Database
DATABASE_URL=postgres://qisur:qisur@localhost:5433/qisur_db
POSTGRES_USER=qisur
POSTGRES_PASSWORD=qisur
POSTGRES_DB=qisur_db

# JWT
JWT_SECRET=change-me-in-production
JWT_EXPIRATION_HOURS=24

# Server
PORT=8080

# Seeders (true para poblar datos de prueba)
SEED=false
```

---

## Ejecución

### Con Docker (recomendado)

```bash
# Levanta PostgreSQL
docker-compose up -d postgres

# Espera a que esté healthy
docker-compose ps

# Ejecuta el server
go run ./cmd/server
```

### Localmente (sin Docker)

```bash
# Asegúrate de tener PostgreSQL corriendo en localhost:5433
go run ./cmd/server
```

El servidor inicia en `http://localhost:8080`

---

## API REST

### Base URL
```
http://localhost:8080/api
```

### Health Check
```http
GET /health
```

### Autenticación

#### Registrar usuario
```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "role": "client",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "role": "client",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Categorías

#### Listar todas
```http
GET /api/categories
```

#### Obtener por ID
```http
GET /api/categories/{id}
```

#### Crear (requiere admin)
```http
POST /api/categories
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Electrónica",
  "description": "Productos electrónicos"
}
```

#### Actualizar (requiere admin)
```http
PUT /api/categories/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Electrónica Premium",
  "description": "Productos electrónicos de alta gama"
}
```

#### Eliminar (requiere admin)
```http
DELETE /api/categories/{id}
Authorization: Bearer <token>
```

### Productos

#### Listar todos
```http
GET /api/products
```

**Parámetros de query (opcionales):**
- `page` — número de página (default: 1)
- `limit` — resultados por página (default: 10)

#### Obtener por ID
```http
GET /api/products/{id}
```

#### Crear (requiere admin)
```http
POST /api/products
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Laptop Dell XPS 13",
  "description": "Laptop ultraligera con pantalla OLED",
  "price": 1299.99,
  "stock": 15,
  "category_ids": ["015649b-f1c6-461f-8a01-da4dedc108f4"]
}
```

#### Actualizar (requiere admin)
```http
PUT /api/products/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Laptop Dell XPS 13 Plus",
  "description": "Laptop ultraligera - modelo mejorado",
  "price": 1399.99,
  "stock": 12,
  "category_ids": ["015649b-f1c6-461f-8a01-da4dedc108f4"]
}
```

#### Eliminar (requiere admin)
```http
DELETE /api/products/{id}
Authorization: Bearer <token>
```

#### Historial de cambios
```http
GET /api/products/{id}/history?start=2026-06-01&end=2026-06-05&limit=10&offset=0
```

**Parámetros:**
- `start` — fecha inicial (YYYY-MM-DD)
- `end` — fecha final (YYYY-MM-DD)
- `limit` — registros por página (default: 10)
- `offset` — desplazamiento (default: 0)

### Búsqueda

#### Buscar productos
```http
GET /api/search?type=product&q=laptop&min_price=500&max_price=2000&sort_by=price&order=ASC&page=1&limit=10
```

**Parámetros:**
- `type` — `product` o `category` (requerido)
- `q` — texto a buscar
- `min_price` — precio mínimo
- `max_price` — precio máximo
- `sort_by` — `name`, `price`, `stock`, `created_at`
- `order` — `ASC` o `DESC`
- `page` — número de página
- `limit` — resultados por página

#### Buscar categorías
```http
GET /api/search?type=category&q=electr&sort_by=name&order=ASC&page=1&limit=10
```

---

## WebSockets

### Conectar

```javascript
const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...";
const ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

ws.onopen = () => console.log("Conectado");

ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    console.log("Evento:", message.event);
    console.log("Datos:", JSON.parse(message.data));
};

ws.onerror = (err) => console.error("Error:", err);
ws.onclose = () => console.log("Desconectado");
```

### Eventos

#### Creación de producto
```json
{
  "event": "product_created",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Laptop",
    "price": 1299.99,
    "stock": 15,
    ...
  }
}
```

#### Actualización de producto
```json
{
  "event": "product_updated",
  "data": { ... }
}
```

#### Eliminación de producto
```json
{
  "event": "product_deleted",
  "data": { "id": "550e8400-e29b-41d4-a716-446655440000" }
}
```

Mismo patrón para categorías: `category_created`, `category_updated`, `category_deleted`

---

## Autenticación

### JWT Token

El token se obtiene en `/api/auth/login` o `/api/auth/register`.

**Enviar en requests protegidas:**
```http
Authorization: Bearer <token>
```

**Duración:** 24 horas (configurable en `.env` con `JWT_EXPIRATION_HOURS`)

### Roles

| Rol    | Permisos                                    |
|--------|---------------------------------------------|
| client | GET endpoints públicos, WS                  |
| admin  | GET endpoints, POST/PUT/DELETE              |

---

## Arquitectura

### Estructura del proyecto

```
qisurChallenge/
├── cmd/
│   └── server/
│       └── main.go              # Punto de entrada
├── internal/
│   ├── auth/
│   │   └── jwt.go               # JWT generation & validation
│   ├── config/
│   │   └── config.go            # Configuración de entorno
│   ├── db/
│   │   ├── postgres.go          # Pool de conexiones
│   │   └── migrations/          # SQL migrations
│   ├── handlers/
│   │   ├── auth_handler.go      # Auth endpoints
│   │   ├── product_handler.go   # Product CRUD
│   │   ├── category_handler.go  # Category CRUD
│   │   ├── search_handler.go    # Search endpoints
│   │   └── ws_handler.go        # WebSocket endpoint
│   ├── middleware/
│   │   ├── auth.go              # JWT validation
│   │   └── role.go              # Role-based access
│   ├── models/
│   │   ├── user.go
│   │   ├── product.go
│   │   └── category.go
│   ├── repository/
│   │   ├── user_repo.go
│   │   ├── product_repo.go
│   │   ├── category_repo.go
│   │   └── search_repo.go
│   ├── router/
│   │   └── router.go            # Rutas y configuración
│   ├── utils/
│   │   └── helpers.go           # Funciones auxiliares
│   └── ws/
│       ├── hub.go               # WebSocket hub
│       └── client.go            # WebSocket cliente
├── seeds/
│   └── seeder.go                # Datos de prueba
├── collections/
│   ├── Qisur-Products-API.json  # Colección Insomnia
│   └── README.md                # Docs de colección
├── docs/
│   └── README.md                # Esta documentación
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── .env
```

### Base de datos

**Tablas:**
- `users` — usuarios y autenticación
- `categories` — categorías de productos
- `products` — inventario de productos
- `product_category` — relación N:M productos-categorías
- `product_history` — historial de cambios (precio/stock)

### Flujo de una request

```
HTTP Request
    ↓
Router (chi)
    ↓
Middleware (Auth + Role)
    ↓
Handler (validación + lógica)
    ↓
Repository (queries a BD)
    ↓
Response JSON
```

---

## Testing

### Con curl

```bash
# Registrar
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}'

# Login y copiar token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}'

# Crear producto (requiere cambiar rol a admin primero)
curl -X POST http://localhost:8080/api/products \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","price":100,"stock":10,"category_ids":[]}'
```

### Con Insomnia

Importa la colección desde `collections/Qisur-Products-API.json`

### Con wscat

```bash
wscat -c "ws://localhost:8080/ws?token=<token>"
```

---

## Notas

- Los IDs son UUIDs v4
- Las fechas usan ISO 8601 (UTC)
- Los precios son decimales (NUMERIC 12,2)
- El stock es un entero no negativo
- Los cambios de precio/stock se registran automáticamente en el historial
