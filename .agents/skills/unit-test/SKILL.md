---
name: unit-test
description: Panduan penulisan unit test, integration test, mock dependencies, dan eksekusi test runner untuk memverifikasi bugfix dan fitur baru sebelum commit.
---

# Skill: Unit & Integration Testing

Gunakan skill ini ketika merancang dan mengimplementasikan pengujian otomatis pada codebase.

## Prinsip Pengujian
1. **Reproduksi Bug Terlebih Dahulu (Regression Test):** Saat memperbaiki bug, buat test case yang mereproduksi kegagalan (test gagal / red) sebelum menulis perbaikan kode (test lolos / green).
2. **Mocking Layanan Eksternal:** Selalu mock panggilan jaringan (HTTP API), disk IO yang tidak perlu, dan layanan pihak ketiga (email, payment gateway, OAuth).
3. **Cakupan Skenario Positif & Negatif:** Uji jalur sukses (*happy path*), jalur gagal (*unhappy path*), input kosong/null/undefined, boundary limits, dan exception handling.

---

## Panduan Berdasarkan Ekosistem

### Node.js / TypeScript (Jest, Vitest, Mocha)
- Gunakan `vi.fn()` / `jest.fn()` untuk spy dan mock function.
- Gunakan async/await assertion yang tepat (`await expect(...).rejects.toThrow(...)`).
- Pastikan membersihkan mock state setelah setiap test (`afterEach(() => { vi.clearAllMocks(); })`).

### Python (Pytest, Unittest)
- Gunakan fixtures untuk setup state reusable.
- Gunakan `unittest.mock.patch` untuk mocking dependensi.
- Pastikan exception ditangkap dengan `pytest.raises(CustomException)`.

### Go (testing package)
- Gunakan *table-driven tests* dengan subtests (`t.Run(name, func(t *testing.T) { ... })`).
- Gunakan interfaces untuk memudahkan dependency injection dan mocking.

---

## Verifikasi Wajib
Sebelum menandai pekerjaan selesai:
- Jalankan test suite lokal secara penuh.
- Pastikan tidak ada test case lain yang *flaky* atau *broken* akibat perubahan baru.
