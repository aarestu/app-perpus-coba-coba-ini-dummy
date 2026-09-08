---
name: action-plan-negotiation
description: Panduan untuk memproses feedback pengguna pada Action Plan di komentar GitHub Issue, menyesuaikan scope, memperbarui analisis teknis, dan menjaga konsistensi sesi percakapan sebelum eksekusi kode.
---

# Skill: Action Plan Negotiation & Feedback Loop

Gunakan skill ini ketika tiket berlabel `status:plan-review` dan pengguna memberikan komentar/koreksi pada rencana kerja yang telah diposting.

## Prinsip Kerja
1. **Pertahankan Konteks Sesi:** Agen bekerja dalam file sesi yang sama (`issue-<id>.json`). Jangan menghilangkan temuan investigasi sebelumnya.
2. **Fleksibel & Adaptif:** Terima masukan batasan teknis dari pengguna (misal: "Jangan gunakan library eksternal baru", "Gunakan helper existing di `utils/`", "Ubah signature fungsi").
3. **Jangan Tulis Kode:** Pada fase ini, kode aplikasi tetap **TIDAK BOLEH** diubah.

---

## Prosedur Kerja

### 1. Ekstraksi Feedback Pengguna
- Baca komentar terbaru pengguna dari issue.
- Identifikasi poin koreksi, preferensi arsitektur, atau batasan baru yang diminta.

### 2. Evaluasi Ulang Action Plan
- Sesuaikan file target atau langkah implementasi sesuai preferensi pengguna.
- Periksa apakah feedback menimbulkan dependensi baru atau potensi dampak samping pada modul lain.
- Perbarui daftar skill jika ada perubahan kebutuhan (misal: perlu tambahan `skill:db-migration`).

### 3. Posting Pembaruan Plan
Buat respons komentar di issue yang merangkum poin koreksi yang telah diakomodasi.

---

## Format Standar Respons Pembaruan

```markdown
## 📝 Pembaruan Action Plan (Berdasarkan Feedback)

Terima kasih atas masukannya. Berikut poin penyesuaian yang telah diakomodasi:

### 1. Poin Koreksi yang Diterapkan
- **[Koreksi 1]:** [Penjelasan bagaimana masukan pengguna diakomodasi]
- **[Koreksi 2]:** ...

### 2. Action Plan Terbaru
- [ ] **Langkah 1 (`path/to/file`):** ...
- [ ] **Langkah 2 (`path/to/file`):** ...
- [ ] **Langkah 3 (Pengujian & Validasi):** ...

Silakan periksa kembali rencana yang telah diperbarui. Jika sudah disetujui, silakan pasang label `status:approved` atau ketik `/approve`.
```
