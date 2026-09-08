---
name: api-design
description: Standar perancangan API RESTful, format response konsisten, HTTP status code yang tepat, skema validasi request, dan error handling terstruktur.
---

# Skill: API Design & Backend Best Practices

Gunakan skill ini ketika merancang, menambah, atau merefaktor endpoint API backend.

## 1. Konvensi Endpoint & HTTP Methods
- Gunakan kata benda jamak untuk resource (contoh: `GET /api/v1/users`, `POST /api/v1/orders`).
- Gunakan method yang sesuai:
  - `GET`: Mengambil data (Idempotent, tanpa side effect).
  - `POST`: Membuat data baru.
  - `PUT`: Mengganti resource secara penuh.
  - `PATCH`: Mengubah sebagian field resource.
  - `DELETE`: Menghapus resource.

## 2. Struktur Format Response Konsisten

### Response Sukses (200 / 201)
```json
{
  "success": true,
  "data": { ... },
  "meta": {
    "page": 1,
    "total": 100
  }
}
```

### Response Error (4xx / 5xx)
```json
{
  "success": false,
  "error": {
    "code": "RESOURCE_NOT_FOUND",
    "message": "User dengan ID 123 tidak ditemukan.",
    "details": []
  }
}
```

## 3. Validasi & Keamanan Input
- Validasi seluruh input payload (`body`, `query`, `params`) menggunakan schema validator (Zod, Joi, Pydantic, dsb).
- Lindungi endpoint dari SQL Injection, NoSQL Injection, dan XSS.
- Terapkan rate limiting dan sanitasi data sensitif (password, secret token) agar tidak bocor ke logging/response.
