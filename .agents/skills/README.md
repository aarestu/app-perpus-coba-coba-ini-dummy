# 📚 Agent Skill Registry (`.agents/skills/`)

Katalog skill ini digunakan oleh **Antigravity IDE** dan **Pi Coding Agent** dalam ekosistem daemon `pi-github-loop`. Agen secara otomatis memindai direktori ini pada fase investigasi, memilih skill yang relevan, dan menyematkan label `skill:<nama>` pada GitHub Issue.

---

## Daftar Skill & Fase Terkait

| Nama Skill | Direktori | Fase Utama | Deskripsi Singkat |
| :--- | :--- | :--- | :--- |
| **`deep-investigation`** | [`deep-investigation/SKILL.md`](file:///d:/project/pi-github-loop/.agents/skills/deep-investigation/SKILL.md) | Fase 1 (Investigasi) | Tracing kode, rekonstruksi skenario kegagalan kronologis, dan perumusan Action Plan. |
| **`action-plan-negotiation`** | [`action-plan-negotiation/SKILL.md`](file:///d:/project/pi-github-loop/.agents/skills/action-plan-negotiation/SKILL.md) | Fase 2 (Loop Issue) | Penyesuaian rencana kerja berdasarkan komentar/masukan pengguna. |
| **`code-execution-pr`** | [`code-execution-pr/SKILL.md`](file:///d:/project/pi-github-loop/.agents/skills/code-execution-pr/SKILL.md) | Fase 3 (Eksekusi) | Implementasi kode, validasi pengujian lokal, dan pembukaan PR otomatis. |
| **`pr-code-review`** | [`pr-code-review/SKILL.md`](file:///d:/project/pi-github-loop/.agents/skills/pr-code-review/SKILL.md) | Fase 4 (Loop Review) | Penanganan catatan review (inline comments) dan refactor terarah. |
| **`unit-test`** | [`unit-test/SKILL.md`](file:///d:/project/pi-github-loop/.agents/skills/unit-test/SKILL.md) | Kapabilitas Teknis | Standar penulisan unit test, mocking, dan verifikasi test suite. |
| **`api-design`** | [`api-design/SKILL.md`](file:///d:/project/pi-github-loop/.agents/skills/api-design/SKILL.md) | Kapabilitas Teknis | Standar perancangan endpoint REST, format response, dan validasi input. |
| **`db-migration`** | [`db-migration/SKILL.md`](file:///d:/project/pi-github-loop/.agents/skills/db-migration/SKILL.md) | Kapabilitas Teknis | Panduan migrasi database yang aman, indexing, dan rollback. |
| **`git-workflow`** | [`git-workflow/SKILL.md`](file:///d:/project/pi-github-loop/.agents/skills/git-workflow/SKILL.md) | Standar Git | Konvensi branch `ai/issue-<id>`, Conventional Commits, dan PR linking. |
