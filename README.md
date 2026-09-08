# App Perpustakaan - Modul Autentikasi

Aplikasi backend perpustakaan berbasis **Golang** dan **SQLite** yang menyediakan sistem autentikasi pengguna dan proteksi endpoint menggunakan JSON Web Token (JWT).

## Fitur
- **Pendaftaran Pengguna (`POST /api/auth/register`):** Registrasi akun baru dengan validasi format input dan enkripsi password menggunakan Bcrypt.
- **Login Pengguna (`POST /api/auth/login`):** Verifikasi kredensial email & password serta penerbitan token JWT.
- **Profil Pengguna (`GET /api/auth/me`):** Endpoint terproteksi menggunakan JWT Auth Middleware.
- **Health Check (`GET /api/health`):** Endpoint pengecekan status server.
- **Database SQLite:** Penyimpanan data persisten pada tabel `users` dengan index unik pada kolom email.

## Konvensi API Response

### Response Sukses (HTTP 200 / 201)
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOi...",
    "user": {
      "id": 1,
      "name": "Budi",
      "email": "budi@perpus.local",
      "role": "user",
      "height": 170,
      "created_at": "2026-09-08T16:00:00Z",
      "updated_at": "2026-09-08T16:00:00Z"
    }
  }
}
```

### Response Gagal (HTTP 4xx / 5xx)
```json
{
  "success": false,
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "email atau password yang dimasukkan salah"
  }
}
```

## Menjalankan Server
```bash
go run cmd/server/main.go
```

Variabel Lingkungan yang dapat dikonfigurasi:
- `PORT`: Port server (default: `8080`)
- `DB_PATH`: Lokasi file SQLite (default: `app.db`)
- `JWT_SECRET`: Secret key signing JWT

## Menjalankan Pengujian
```bash
go test -v ./...
```
