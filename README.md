# marketing-workspace

Workspace internal Marketing Unggulmart untuk task, review opsional dua arah, event, kalender WITA, reminder, notifikasi, dan PWA. Stack yang direncanakan: SvelteKit/TypeScript, Go/Chi, dan PostgreSQL Neon.

Status per 29 Oktober 2026: spesifikasi tersedia; implementasi aplikasi belum dimulai. Status task/gate serta bukti terbaru dikelola hanya di [TODO.md](TODO.md).

Seluruh dokumen Markdown proyek berada di root repository:

| Dokumen | Peran |
|---|---|
| [PRD.md](PRD.md) | Kebutuhan produk, aturan bisnis, milestone, dan pertanyaan terbuka |
| [ARCHITECTURE.md](ARCHITECTURE.md) | Ownership, struktur kode, dan aliran data |
| [AGENTS.md](AGENTS.md) | Aturan agent dan invariant teknis |
| [TODO.md](TODO.md) | Task, dependensi, keputusan terbuka, progres, dan bukti |
| [SKILL.md](SKILL.md) | Panduan workflow berulang untuk proyek ini |
| [WORKFLOW.md](WORKFLOW.md) | Otonomi, izin, verifikasi, dan serah terima |

Dokumen pendukung architecture-decisions.md, deployment.md, runbook.md, dan CALENDAR-DATA.md juga akan dibuat di root saat dibutuhkan milestone. Kontrak API akan berada di openapi.yaml. File tersebut belum dibuat; tidak ada klaim command setup/test/build atau deployment sudah tersedia.

Kode tetap direncanakan di apps/web dan apps/api, SQL migration di db/migrations/, dan asset dataset Hijriyah di calendar-data/. Urutan berikutnya adalah M1-01; default produk tetap mengikuti PRD sampai dikonfirmasi sesuai tenggatnya.
