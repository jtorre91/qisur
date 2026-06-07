# Qisur Challenge - Backend API

> Backend REST API con WebSockets en tiempo real para gestión de productos y categorías

## 🚀 Quick Start

### Requisitos

- Go 1.21+
- Docker & Docker Compose

### Instalación

```bash
# Clonar
git clone https://github.com/jtorre91/qisur
cd qisur

# Descargar dependencias
go mod download


# Crear .env en la raiz del proyecto (pueden copiar directamente el env.example)

# Database
DATABASE_URL=postgres://qisur:qisur@localhost:5433/qisur_db

# JWT
JWT_SECRET=change-me-in-production
JWT_EXPIRATION_HOURS=24

# Server
PORT=8080

# Populate datos de prueba
SEED=True

# Levantar PostgreSQL
docker-compose up -d postgres

# Ejecutar servidor (puede tardar unos minutos la 1ra vez)
go run ./cmd/server
```

El servidor inicia en `http://localhost:8080`

## 📖 Documentación

- **[API Completa](./docs/README.md)** — Endpoints, ejemplos, WebSockets
- **[Swagger UI](http://localhost:8080/swagger/index.html)** — Documentación interactiva
- **[Colección Insomnia](./collections/README.md)** — Tests listos para usar

## ✨ Features

### Autenticación
- ✅ JWT tokens con expiración configurable
- ✅ Registro y login de usuarios
- ✅ Roles basados en acceso (admin/client)

### REST API
- ✅ CRUD completo de productos
- ✅ CRUD completo de categorías
- ✅ Relación N:M productos-categorías
- ✅ Historial de cambios (precio/stock)
- ✅ Búsqueda con filtros y ordenamiento
- ✅ Paginación

### Real-time
- ✅ WebSockets con eventos en tiempo real
- ✅ Broadcast automático en mutaciones
- ✅ Conexiones autenticadas

### Base de datos
- ✅ PostgreSQL 16
- ✅ Migrations automáticas
- ✅ Datos de prueba (seeders)

## 📊 Stack

- **Go 1.21+** — Lenguaje principal
- **Chi v5** — Router HTTP
- **pgx v5** — Driver PostgreSQL
- **gorilla/websocket** — WebSockets
- **golang-jwt** — Autenticación JWT
- **Swagger** — Documentación API

## 📂 Estructura del Proyecto

```
qisur/
├── cmd/server/              # Punto de entrada
├── internal/
│   ├── auth/               # Autenticación JWT
│   ├── config/             # Configuración
│   ├── db/                 # PostgreSQL + migrations
│   ├── handlers/           # HTTP handlers
│   ├── middleware/         # Auth + Role guards
│   ├── models/             # Data structures
│   ├── repository/         # Data access
│   ├── router/             # Rutas
│   ├── utils/              # Helpers
│   └── ws/                 # WebSocket hub + client
├── docs/                   # Documentación
├── collections/            # Insomnia collection
├── seeds/                  # Test data
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── README.md
```

## 🏗️ Arquitectura

```
Clients (HTTP + WebSocket)
        ↓
    Router (Chi)
        ↓
    Middleware (Auth + Role)
        ↓
    Handlers (Validación)
        ↓
    Repository (BD)
        ↓
    PostgreSQL
```

### Flujo de WebSocket

```
Cliente conecta → WS upgrades → Hub registra cliente
                                      ↓
                          POST/PUT/DELETE en REST
                                      ↓
                      hub.Broadcast() → todos reciben evento
```

## 📝 Ejemplos

### Registrar usuario

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Crear producto (requiere JWT + admin)

```bash
curl -X POST http://localhost:8080/api/products \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop Dell",
    "price": 1299.99,
    "stock": 15,
    "category_ids": ["uuid-categoria"]
  }'
```

### Conectar WebSocket

```javascript
const ws = new WebSocket('ws://localhost:8080/ws?token=<jwt>');

ws.onmessage = (event) => {
  const { event, data } = JSON.parse(event.data);
  console.log(`Evento: ${event}`, data);
};
```

## 🧪 Testing

### Con Insomnia

1. Abre `collections/Qisur-Products-API.json`
2. Registra usuario
3. Copia JWT del response
4. Usa en requests protegidas

### Con curl

```bash
# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"pass"}'

# Crear producto
curl -X POST http://localhost:8080/api/products \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","price":100,"stock":10,"category_ids":[]}'
```

## 📚 Permisos

| Endpoint | Public | Client | Admin |
|----------|:------:|:------:|:-----:|
| GET /products | ✓ | ✓ | ✓ |
| GET /categories | ✓ | ✓ | ✓ |
| GET /search | ✓ | ✓ | ✓ |
| POST/PUT/DELETE /products | | | ✓ |
| POST/PUT/DELETE /categories | | | ✓ |
| GET /ws | | ✓ | ✓ |

## 🔐 Seguridad

- ✅ JWT tokens firmados con HS256
- ✅ Contraseñas hasheadas con bcrypt
- ✅ Validación de entrada en handlers
- ✅ Role-based access control
- ✅ CORS configurado para navegador

## 🚀 Despliegue

### Con Docker

```bash
docker-compose up -d
```

Incluye:
- PostgreSQL container
- Go app container
- Network bridge

## 📖 Documentación Completa

- **[docs/README.md](./docs/README.md)** — API REST detallada
- **[Swagger](http://localhost:8080/swagger/index.html)** — Documentación interactiva
- **[collections/README.md](./collections/README.md)** — Ejemplos con Insomnia


## 👨‍💻 Autor

Javier A. Torre

---

**¿Necesitas ayuda?**
- 📖 Lee la [documentación completa](./docs/README.md)
- 🧪 Prueba con [Insomnia](./collections/Qisur-Products-API.json)
- 💬 Revisa los [ejemplos de curl](#ejemplos)
