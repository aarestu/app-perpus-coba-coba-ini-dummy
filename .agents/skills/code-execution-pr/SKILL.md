---
name: code-execution-pr
description: Panduan eksekusi penulisan kode, pengetesan lokal (test & linter), commit rapi, dan pembukaan Pull Request otomatis setelah Action Plan disetujui.
---

# Skill: Code Execution & Pull Request Creation

Gunakan skill ini ketika tiket berpindah status menjadi `status:approved` dan `turn:ai`.

## Aturan Utama
1. **Patuhi Action Plan:** Terapkan perubahan kode secara presisi sesuai rencana yang telah disetujui pada fase sebelumnya.
2. **Minimal & Terisolasi:** Hanya ubah file dan baris yang relevan. Jangan melakukan format ulang (*reformat*) pada seluruh file yang tidak berkaitan.
3. **Verifikasi Lokal:** Selalu jalankan *test suite* atau *linter* lokal sebelum melakukan commit.
4. **Clean Git History:** Gunakan standar *Conventional Commits* dan sertakan referensi penutupan issue (`Closes #<id>`).

---

## Prosedur Eksekusi

### 1. Implementasi Kode
- Modifikasi file sesuai langkah-langkah di Action Plan.
- Tambahkan validasi error handling, edge cases, dan logging yang memadai.
- Jaga konsistensi gaya kode (naming convention, arsitektur folder, type safety).

### 2. Validasi & Pengujian Lokal
- Jalankan test runner lokal (misal: `npm test`, `pytest`, `cargo test`, `go test ./...`).
- Jalankan linter / typecheck (misal: `npm run lint`, `tsc --noEmit`, `flake8`).
- Jika terjadi error pada test/linter, perbaiki langsung sebelum commit.

### 3. Commit & Format Pesan
Format commit yang wajib diikuti:
```
<type>(<scope>): <deskripsi singkat perubahan>

[Penjelasan detail jika diperlukan]

Closes #<issue_number>
```
Contoh:
- `fix(auth): handle token expiration gracefully in refresh interceptor`
- `feat(rate-limit): implement redis token bucket middleware`

### 4. Rangkuman Pembuatan PR
Deskripsi PR harus mencakup:
- Ringkasan masalah dan solusi teknis.
- Daftar file yang diubah beserta penjelasannya.
- Hasil verifikasi pengujian lokal.
- Referensi issue penutup (`Closes #<id>`).
