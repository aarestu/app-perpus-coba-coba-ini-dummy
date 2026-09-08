---
name: deep-investigation
description: Panduan investigasi mendalam untuk melacak codebase, merekonstruksi skenario kronologis terjadinya bug, menganalisis root cause, dan merancang Action Plan terstruktur tanpa mengubah kode.
---

# Skill: Deep Investigation & Scenario Reconstruction

Gunakan skill ini pada fase investigasi tiket (ketika tiket berlabel `status:investigating` dan `turn:ai`).

## Tujuan Utama
1. Menemukan akar permasalahan (*Root Cause Analysis*) secara presisi di codebase.
2. Merekonstruksi skenario kegagalan secara kronologis langkah demi langkah (*Scenario Reconstruction*).
3. Merumuskan *Action Plan* yang terukur, terarah, dan aman.
4. **DILARANG** melakukan modifikasi/penulisan file kode apapun pada fase ini.

---

## Prosedur Pelaksanaan

### 1. Pindai Codebase (Code Tracing)
- Cari fungsi, rute API, middleware, atau file konfigurasi terkait dari kata kunci tiket.
- Telusuri alur pemanggilan (*call stack* & *data flow*) dari titik masuk (input/request) hingga titik terjadinya kegagalan.
- Periksa dependensi, penanganan error (`try/catch`, async rejection), dan interaksi database/layanan pihak ketiga.

### 2. Rekonstruksi Skenario (Kronologi Teknis)
Jelaskan urutan kejadian teknis secara terperinci:
1. **Langkah 1 (Trigger):** Request/aksi apa yang memicu alur kerja.
2. **Langkah 2 (State & Payload):** Data atau parameter apa yang dikirimkan.
3. **Langkah 3 (Failure Point):** Baris kode mana yang mengalami kegagalan dan mengapa (misal: state tidak tervalidasi, unhandled promise, token expired tanpa refresh handling).
4. **Langkah 4 (Impact):** Dampak langsung pada sistem (misal: HTTP 500, silent fail, data corruption).

### 3. Pemilihan Skill Pendukung
Pindai direktori skill dan deklarasikan skill apa saja yang akan diaktifkan untuk fase eksekusi (contoh: `skill:unit-test`, `skill:api-design`, `skill:db-migration`).

### 4. Rancang Usulan Action Plan
Susun rencana kerja dengan struktur:
- **Tujuan / Scope:** Apa yang akan diperbaiki dan apa yang berada di luar lingkup.
- **File Target:** Daftar file yang akan ditambah atau dimodifikasi.
- **Langkah Perubahan Teknis:** Rincian modifikasi pada masing-masing file.
- **Rencana Pengujian:** Unit test / integration test yang akan dibuat atau dijalankan untuk memvalidasi perbaikan.
- **Potensi Risiko & Mitigasi:** Kemungkinan *breaking changes* atau efek samping.

---

## Format Standar Laporan Komentar GitHub Issue

```markdown
## 🔍 Laporan Investigasi Teknis

### 1. Ringkasan Masalah
[Penjelasan singkat mengenai akar permasalahan]

### 2. Rekonstruksi Skenario Kronologis
1. **Pemicu:** ...
2. **Alur Eksekusi:** ...
3. **Titik Kegagalan (`file:baris`):** ...
4. **Dampak:** ...

### 3. Skill yang Dipilih
- `skill:<nama-skill-1>`: [Alasan pemilihan]
- `skill:<nama-skill-2>`: [Alasan pemilihan]

### 4. Usulan Action Plan
- [ ] **Langkah 1 (`path/to/file`):** [Deskripsi perubahan]
- [ ] **Langkah 2 (`path/to/file`):** [Deskripsi perubahan]
- [ ] **Langkah 3 (Testing):** [Rencana pengetesan]

### 5. Konfirmasi Pengguna
Mohon berikan masukan atau koreksi pada rencana di atas. Jika sudah sesuai, ubah label menjadi `status:approved` atau balas `/approve` untuk memulai eksekusi kode.
```
