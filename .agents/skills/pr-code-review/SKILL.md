---
name: pr-code-review
description: Panduan untuk menangani feedback code review pada Pull Request (inline comments dan review body), menerapkan refactoring terarah, dan mengonfirmasi perbaikan di thread PR.
---

# Skill: PR Code Review & Inline Feedback Handling

Gunakan skill ini ketika Pull Request menerima feedback perbaikan dengan label `status:changes-requested` dan `turn:ai`.

## Prinsip Kerja
1. **Pertahankan Konteks Keseluruhan:** Pahami alasan arsitektural dari fase investigasi dan issue sebelum melakukan perubahan.
2. **Bedah Feedback Review:** Baca teliti catatan review reviewer, baik yang bersifat umum (review body) maupun *inline code comments* pada baris tertentu.
3. **Refactoring Bedah (Surgical Refactor):** Ubah hanya bagian yang diminta oleh reviewer tanpa merusak fungsionalitas yang sudah bekerja.
4. **Validasi Ulang:** Jalankan pengetesan lokal untuk memastikan tidak terjadi regresi (*regression*).

---

## Prosedur Pelaksanaan

### 1. Analisis Catatan Review
- Baca semua review comment yang dikumpulkan dari PR.
- Kelompokkan masukan ke dalam kategori:
  - *Bug / Edge-case fixes*
  - *Performance / Optimization*
  - *Clean code / Naming / Styling*
  - *Additional test cases*

### 2. Terapkan Koreksi pada Branch PR
- Modifikasi baris kode yang disorot oleh reviewer.
- Pastikan perubahan tetap selaras dengan modul lain.

### 3. Pengetesan Lokal
- Jalankan *unit test suite* dan *typecheck*.
- Tambahkan test case baru jika review meminta cakupan pengujian edge-case tambahan.

### 4. Commit & Konfirmasi
- Lakukan commit dengan format: `refactor(review): <ringkasan perbaikan>`.
- Kirim balasan komentar di PR yang merinci file apa saja yang diperbaiki dan bagaimana catatan reviewer telah dipenuhi.

---

## Format Standar Komentar Konfirmasi Review

```markdown
## ✅ Perbaikan Code Review Diterapkan

Berikut ringkasan perbaikan yang telah di-commit ke branch PR ini:

1. **[Topik Review 1]:** [Penjelasan perbaikan yang dilakukan pada `path/to/file:Lxx`]
2. **[Topik Review 2]:** ...

### Status Verifikasi Pengujian
- Unit Test: 🟢 Passed
- Linter / Typecheck: 🟢 Passed

Silakan dilakukan review ulang.
```
