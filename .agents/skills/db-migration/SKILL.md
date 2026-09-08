---
name: db-migration
description: Panduan perubahan skema database, migrasi aman (safe migrations), indexing, rollback strategy, dan optimasi query data.
---

# Skill: Database Migration & Schema Evolution

Gunakan skill ini ketika tugas melibatkan modifikasi skema tabel, relasi, index, atau migrasi data.

## Prinsip Migrasi Aman (Safe Migrations)
1. **Backward Compatibility:** Perubahan skema tidak boleh memutus aplikasi yang sedang berjalan (gunakan pendekatan *Expand & Contract* / *Two-phase migration*).
2. **Hindari Lock Berkepanjangan:** Hindari `ALTER TABLE` besar yang mengunci tabel di database production.
3. **Sediakan Mekanisme Rollback:** Setiap skrip migrasi naik (*up*) harus memiliki skrip pembalikan (*down/rollback*) jika memungkinkan.

---

## Prosedur Kerja

### 1. Modifikasi Skema
- Tentukan ORM / Migration tool yang digunakan proyek (Prisma, TypeORM, Drizzle, Knex, Alembic, Flyway, Go-Migrate).
- Buat file migrasi baru menggunakan generator CLI tool tersebut (hindari mengedit file migrasi yang sudah pernah di-*apply* sebelumnya).

### 2. Aturan Kolom & Indexing
- Saat menambahkan kolom `NOT NULL` baru, sertakan nilai `DEFAULT` agar migrasi pada data yang sudah ada tidak gagal.
- Tambahkan index pada kolom foreign key atau kolom yang sering dijadikan filter query (`WHERE`, `ORDER BY`, `JOIN`).
- Gunakan tipe data yang presisi (hindari tipe data terlalu luas/boros memori).

### 3. Uji Coba Migrasi Lokal
- Jalankan perintah migrasi ke database testing lokal.
- Verifikasi status migrasi dan uji apakah rollback berjalan mulus tanpa error data.
