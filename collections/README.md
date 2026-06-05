# Insomnia Collections

Colección unificada de API para testear los endpoints de Qisur Challenge.

## Archivo principal

- **Qisur-Products-API.json** — Colección completa con carpetas por módulo

## Estructura de la colección

```
Qisur Products API/
├── 🏥 Health Check
├── 🔐 Auth
│   ├── POST - Register
│   └── POST - Login
├── 📁 Categories
│   ├── GET - List all categories
│   ├── GET - Get category by ID
│   ├── POST - Create category
│   ├── PUT - Update category
│   └── DELETE - Delete category
├── 📦 Products
│   ├── GET - List all products
│   ├── GET - Get product by ID
│   ├── POST - Create product
│   ├── PUT - Update product
│   ├── DELETE - Delete product
│   └── GET - Product history
├── 🔐 Auth (próximamente)
├── 🔍 Search
│   ├── GET - Search products (completo con filtros)
│   ├── GET - Search products (simple)
│   └── GET - Search categories
└── Base Environment (variables)
```

## Cómo importar en Insomnia

1. Abre **Insomnia**
2. Click en **File** → **Import**
3. Selecciona **From File**
4. Elige `Qisur-Products-API.json`
5. Se importará la colección completa con todas las carpetas

## Cómo agregar nuevos endpoints

Cuando agregues nuevas features (Auth, Search, etc.):
1. Simplemente actualiza este archivo agregando nuevas carpetas/requests
2. Reimporta en Insomnia para tener la última versión
3. Sin necesidad de múltiples archivos

## Variables de entorno

La colección usa variables para mayor flexibilidad:

```
base_url = http://localhost:8080
api_url = http://localhost:8080/api
```

Puedes cambiarlas según tu ambiente (desarrollo, staging, producción).

## Endpoints disponibles

### Health Check
- **GET** `/health` — Verifica que el servidor esté activo

### Auth (Autenticación)
- **POST** `/api/auth/register` — Registra un nuevo usuario
- **POST** `/api/auth/login` — Login y obtén JWT

### Categories (Categorías)
- **GET** `/api/categories` — Lista todas las categorías
- **GET** `/api/categories/{id}` — Obtiene una categoría
- **POST** `/api/categories` — Crea una nueva categoría
- **PUT** `/api/categories/{id}` — Actualiza una categoría
- **DELETE** `/api/categories/{id}` — Elimina una categoría

### Products (Productos)
- **GET** `/api/products` — Lista todos los productos
- **GET** `/api/products/{id}` — Obtiene un producto
- **POST** `/api/products` — Crea un nuevo producto
- **PUT** `/api/products/{id}` — Actualiza un producto
- **DELETE** `/api/products/{id}` — Elimina un producto
- **GET** `/api/products/{id}/history` — Obtiene historial de cambios (precio y stock)

## Ejemplos de uso

### Registrar un usuario
```json
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

### Login
```json
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

**Nota:** El token devuelto es un JWT válido por 24 horas (configurable en `.env` con `JWT_EXPIRATION_HOURS`).

### Crear una categoría
```json
POST /api/categories
Content-Type: application/json

{
  "name": "Electrónica",
  "description": "Productos electrónicos"
}
```

### Actualizar una categoría
```json
PUT /api/categories/{id}
Content-Type: application/json

{
  "name": "Electrónica Premium",
  "description": "Productos electrónicos de alta gama"
}
```

### Crear un producto (con categorías)
```json
POST /api/products
Content-Type: application/json

{
  "name": "Laptop Dell XPS 13",
  "description": "Laptop ultraligera con pantalla OLED",
  "price": 1299.99,
  "stock": 15,
  "category_ids": ["015649b-f1c6-461f-8a01-da4dedc108f4", "550e8400-e29b-41d4-a716-446655440000"]
}
```

**Nota sobre `category_ids`:**
- Campo requerido: sí
- Debe ser un array de UUID válidos
- Cada UUID debe corresponder a una categoría existente
- Si la categoría no existe, recibirás error: `"categoría inválida, no existe"`
- Puedes obtener los IDs de categorías con `GET /api/categories`

### Actualizar un producto (con categorías)
```json
PUT /api/products/{id}
Content-Type: application/json

{
  "name": "Laptop Dell XPS 13 Plus",
  "description": "Laptop ultraligera con pantalla OLED - modelo mejorado",
  "price": 1399.99,
  "stock": 12,
  "category_ids": ["015649b-f1c6-461f-8a01-da4dedc108f4"]
}
```

**Nota sobre actualización de categorías:**
- Si envías `category_ids`, se actualiza la relación N:M
- Las categorías se reemplazan completamente (no se agregan, se reemplazan)
- Si solo cambió 1 categoría de 10, solo esa se actualiza (optimizado con diff)
- Si la categoría no existe, devuelve error con el mensaje: `"categoría inválida, no existe"`
- Los cambios en precio/stock se registran en el historial automáticamente

### Obtener historial de cambios de un producto

**Sin filtro de fechas:**
```
GET /api/products/{id}/history?limit=10&offset=0
```

**Con rango de fechas (día/mes/año):**
```
GET /api/products/{id}/history?start=2026-06-01&end=2026-06-04&limit=10&offset=0
```

**Parámetros:**
- `start` — fecha inicial (YYYY-MM-DD) — opcional
- `end` — fecha final (YYYY-MM-DD) — opcional
- `limit` — cantidad de registros por página (default: 10)
- `offset` — desplazamiento para paginación (default: 0)

**Nota:** Las fechas se buscan sin considerar horas/minutos/segundos. 
Si pones `start=2026-06-01` y `end=2026-06-04`, busca todos los cambios ocurridos en esos días completos.

### Buscar productos

**Con todos los filtros:**
```
GET /api/search?type=product&q=laptop&min_price=500&max_price=2000&sort_by=price&order=ASC&page=1&limit=10
```

**Búsqueda simple:**
```
GET /api/search?type=product&q=cable&page=1&limit=10
```

**Parámetros de búsqueda de productos:**
- `type=product` — requerido
- `q` — texto a buscar en nombre y descripción (opcional)
- `min_price` — precio mínimo (opcional)
- `max_price` — precio máximo (opcional)
- `sort_by` — campo para ordenar: `name`, `price`, `stock`, `created_at` (default: `created_at`)
- `order` — `ASC` o `DESC` (default: `DESC`)
- `page` — número de página (default: 1)
- `limit` — resultados por página, máximo 100 (default: 10)

### Buscar categorías

```
GET /api/search?type=category&q=electr&sort_by=name&order=ASC&page=1&limit=10
```

**Parámetros de búsqueda de categorías:**
- `type=category` — requerido
- `q` — texto a buscar en nombre y descripción (opcional)
- `sort_by` — campo para ordenar: `name`, `created_at` (default: `created_at`)
- `order` — `ASC` o `DESC` (default: `DESC`)
- `page` — número de página (default: 1)
- `limit` — resultados por página, máximo 100 (default: 10)

**Respuesta de búsqueda:**
```json
{
  "items": [ ... ],
  "total": 45,
  "page": 1,
  "limit": 10,
  "total_pages": 5
}
```

## Notas

- Reemplaza `PASTE_CATEGORY_ID_HERE` con el ID real de una categoría
- Los IDs son UUID (ejemplo: `550e8400-e29b-41d4-a716-446655440000`)
- Para obtener IDs, usa primero el endpoint de `List all categories`
