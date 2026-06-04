# Insomnia Collections

Colecciones de API para testear los endpoints de Qisur Challenge.

## Cómo importar en Insomnia

### Opción 1: Importar desde archivo
1. Abre **Insomnia**
2. Click en **File** → **Import**
3. Selecciona **From File**
4. Elige `Qisur-API.json`
5. Se importará la colección con todos los endpoints

### Opción 2: Crear manualmente
Si prefieres crear requests desde cero:

1. **New Request** → **HTTP Request**
2. Configura el método (GET, POST, etc.)
3. Usa variables de entorno para la URL base

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

### Categories (Categorías)
- **GET** `/api/categories` — Lista todas las categorías
- **GET** `/api/categories/{id}` — Obtiene una categoría
- **POST** `/api/categories` — Crea una nueva categoría
- **PUT** `/api/categories/{id}` — Actualiza una categoría
- **DELETE** `/api/categories/{id}` — Elimina una categoría

## Ejemplos de uso

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

## Notas

- Reemplaza `PASTE_CATEGORY_ID_HERE` con el ID real de una categoría
- Los IDs son UUID (ejemplo: `550e8400-e29b-41d4-a716-446655440000`)
- Para obtener IDs, usa primero el endpoint de `List all categories`
