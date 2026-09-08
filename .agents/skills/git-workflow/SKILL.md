---
name: git-workflow
description: Standar konvensi Git, penamaan branch ai/issue-<id>, Conventional Commits, penanganan git conflict, dan linking Pull Request ke Issue GitHub.
---

# Skill: Git & GitHub Workflow Standards

Gunakan skill ini sebagai acuan standar konvensi Git dalam seluruh alur kerja agen `pi-github-loop`.

## 1. Penamaan Branch
- Format branch wajib: `ai/issue-<id>`
- Contoh: `ai/issue-42`
- Branch selalu dibuat dari branch utama yang mutakhir (`main` atau `master`).

## 2. Standar Conventional Commits
Format pesan commit:
```
<type>(<scope>): <subject>

[optional body]

[optional footer]
```

### Tipe Commit (`<type>`):
- `feat`: Fitur baru untuk pengguna.
- `fix`: Perbaikan bug.
- `refactor`: Refactoring kode tanpa menambah fitur atau mengubah perilaku eksternal.
- `test`: Menambah atau memperbaiki unit test / integration test.
- `docs`: Dokumentasi file atau README.
- `chore`: Konfigurasi build, dependensi, atau tooling.

### Contoh Pesan Commit:
- `fix(auth): handle expired token error in auth interceptor`
- `feat(cart): add bulk item removal endpoint`
- `refactor(review): simplify conditional branching based on review feedback`

## 3. Menghubungkan Pull Request dengan Issue
Setiap kali membuka Pull Request, sertakan salah satu keyword penutup GitHub di bagian deskripsi PR:
- `Closes #<id>`
- `Fixes #<id>`
- `Resolves #<id>`

Hal ini memastikan GitHub Issue otomatis tertutup saat Pull Request di-merge ke branch utama.
