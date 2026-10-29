# TODO — Unggulmart Marketing Workspace

Diperbarui: 29 Oktober 2026 (WITA) · Lokasi: `TODO.md` di root repository.

Pelacak pekerjaan lokal untuk [PRD](PRD.md), [panduan agent](AGENTS.md), dan [architecture](ARCHITECTURE.md). Seluruh dokumen Markdown proyek berada di root; tautan menggunakan lokasi aktual. Tracker mencatat progres; tidak mengubah requirement atau menggantikan acceptance criteria sumber.

## Baseline dan cara pakai

Baseline sebelum M1-01 hanya berisi dokumen proyek di root. Progres implementasi aktual mengikuti tabel dan bukti di bawah; preview/production tidak dianggap selesai dari hasil lokal. Pembuatan/perbaikan dokumen tidak dihitung sebagai progres implementasi.

1. Baca sumber dan pilih task pada milestone pertama yang belum selesai. Jangan mulai milestone berikut sebelum gate sebelumnya `DONE`.
2. Isi PIC pelaksana aktual dan ubah ke `IN_PROGRESS` saat mulai. PIC tracker adalah developer/agent pelaksana, bukan otomatis Sultan/Alvin sebagai PIC task produk. `—` berarti belum ditetapkan.
3. Perbarui baris task yang sama setelah perubahan: status, PIC, serta bukti/catatan. Jangan membuat salinan checklist status di tempat lain.
4. Gunakan `VERIFY` ketika implementasi tersedia tetapi acceptance/checks belum lengkap. Hanya beri `DONE` setelah Definition of Done terpenuhi dan bukti tersedia.
5. Jika terhambat, gunakan `BLOCKED` dan tulis penyebab, ID dependensi, penanggung jawab tindak lanjut serta langkah pembuka blokir. Dependensi yang belum dikerjakan saja tidak otomatis menjadikan semua backlog `BLOCKED`.
6. Perbarui ringkasan milestone, fokus aktif, issue terbuka dan log pada akhir sesi. Simpan perubahan tracker bersama perubahan kode bila repository sudah tersedia.

| Status | Arti |
|---|---|
| TODO | Belum mulai; termasuk menunggu urutan dependensi |
| IN_PROGRESS | Sedang dikerjakan oleh PIC yang terisi |
| VERIFY | Implementasi tersedia, sedang/menunggu verifikasi |
| BLOCKED | Ada hambatan konkret yang dicatat |
| DONE | Acceptance dan checks relevan lulus, bukti tercatat |
| CANCELLED | Dibatalkan dengan alasan dan keputusan scope yang eksplisit; bukan selesai |

Status tracker ini berbeda dari status task produk `todo/in_progress/review/done`. Task `VERIFY` juga tidak berarti review produk.

## Ringkasan progres

Angka `Selesai/Total` hanya menghitung baris implementasi Mx-NN, tidak termasuk gate, issue tambahan, maupun dokumen yang sudah disusun. Persentase = DONE ÷ total task aktif × 100, bukan estimasi waktu/usaha. CANCELLED dikeluarkan dari denominator hanya setelah keputusan scope tercatat; pengurangan scope tidak otomatis meluluskan gate. Ringkasan manual wajib dihitung ulang saat tabel berubah.

| Milestone | Status | Selesai/Total | Progres | Gate |
|---|---|---|---|---|
| M1 — Foundation, database, auth | IN_PROGRESS | 2/11 | 18,2% | M1-GATE: TODO |
| M2 — Task & Direct Completion | TODO | 0/8 | 0% | M2-GATE: TODO |
| M3 — Optional Review | TODO | 0/8 | 0% | M3-GATE: TODO |
| M4 — Marketing & Dual Calendar | TODO | 0/6 | 0% | M4-GATE: TODO |
| M5 — Event Reminder | TODO | 0/6 | 0% | M5-GATE: TODO |
| M6 — PWA & Push | TODO | 0/6 | 0% | M6-GATE: TODO |
| M7 — Hardening & Production | TODO | 0/8 | 0% | M7-GATE: TODO |
| **Total implementasi** | **IN_PROGRESS** | **2/53** | **3,8%** | **0/7 gate selesai** |

## Fokus sesi berikutnya

| Task | Aksi berikutnya | PIC | Hambatan aktual |
|---|---|---|---|
| M1-03 | Belum dimulai; menunggu instruksi untuk migration runner/database dasar | — | M1-02 selesai lokal; BUG-001 tetap terbuka dan terpisah |

Jangan mengisi tanggal selesai, hasil test, nama branch/commit, atau penerimaan pengguna yang belum terjadi. Tidak ada target tanggal rilis sampai kapasitas pelaksana diketahui.

## Definition of Done per task

Task DONE harus memenuhi kriteria barisnya dan requirement sumber: kontrak/API types serta migration diperbarui bila relevan; aturan bisnis/permission ada di backend; audit/inbox atomik bila dipicu; UI memuat loading/empty/error/conflict sesuai kebutuhan; checks dan test sebanding risiko lulus; docs/env/runbook yang terdampak diperbarui; tidak ada secret atau fitur di luar scope.

Kolom bukti minimal mencatat tanggal WITA, lokasi file/commit/PR bila tersedia, command verifikasi nyata beserta hasil, dan bukti acceptance/manual bila relevan. Jika check tidak dapat dijalankan, catat command, penyebab dan risiko; task tetap VERIFY/BLOCKED bila bukti itu diperlukan untuk selesai. `Kode sudah dibuat` atau `build lulus` saja tidak membuktikan aturan concurrency/security. Tidak perlu test baru untuk perubahan copy kecil tanpa risiko behavior.

Setiap OQ PRD §26.1 dipetakan ke issue keputusan di bawah. Default tetap boleh dipakai untuk pembangunan, tetapi konfirmasi yang diwajibkan harus tercatat pada issue sebelum tenggat task/gate terkait. Jika implementasi tersedia tetapi konfirmasi wajib belum ada, gunakan VERIFY; jangan menutup issue keputusan hanya karena default sudah diimplementasikan. Issue keputusan tidak menambah jumlah 53 task implementasi.

## M1 — Foundation, database, auth

Sumber: PRD §24 dan bagian domain terkait; AGENTS serta architecture berlaku pada semua baris.

| ID | Status | Task | Dependensi | Kriteria selesai | PIC | Bukti / catatan |
|---|---|---|---|---|---|---|
| M1-01 | DONE | Monorepo dan tooling | — | apps/web dan apps/api berjalan lokal; pnpm/Go Modules serta versi kompatibel dikunci. | Codex | 2026-10-29 WITA: apps/web, apps/api/cmd/server, package.json, pnpm-workspace.yaml, pnpm-lock.yaml, go.mod/go.sum, runtime pins dan scripts tersedia. Node 24.22.2/pnpm 11.10.0/Go 1.27.3 terverifikasi; pnpm quality, frozen install, smoke dan Chrome desktop lulus. Unit test bisnis belum ada; detail bukti di bawah, README dan ADR-001. Tidak lanjut M1-02. |
| M1-02 | DONE | Boundary backend dan konfigurasi | M1-01 | Router Chi reusable di server lokal/adapter Vercel; config fail-closed, request ID, errors, timeout dan log teredaksi. | Codex | 2026-10-29 WITA: router Chi bersama server lokal/adapter Vercel, config fail-closed, request ID, JSON error, timeout/cancellation dan log allowlist tersedia. PASS pnpm quality, checks/build backend akhir, Go race dan smoke lokal; detail bukti di bawah serta ADR-002. Tanpa akses database/dependency baru; cloud belum diuji. |
| M1-03 | TODO | Migration runner dan database dasar | M1-02 | Users, sessions, audit, settings, rate limits tersedia; audit belum memiliki event_id sampai M4 sesuai staging PRD §13.5; direct migration terpisah pooled runtime; audit append-only. | — | — |
| M1-04 | TODO | OpenAPI baseline dan API client | M1-02 | Auth/users/settings/health terdokumentasi; generated types, envelope/errors/CSRF terpusat dan contract validation tersedia. | — | — |
| M1-05 | TODO | Bootstrap dan password | M1-03 | CLI bootstrap sekali jalan, Argon2id terukur, temporary password wajib diganti; tidak ada registrasi publik. | — | — |
| M1-06 | TODO | Session, login/logout dan CSRF | M1-04, M1-05 | Token hash, secure cookie, expiry/rotation/revoke, Origin dan CSRF multi-tab diuji; rate limit durable. | — | — |
| M1-07 | TODO | User management dan workspace settings | M1-06 | Supervisor membuat/mengubah/nonaktifkan akun, reset/role dan jam reminder; last-supervisor guard atomik, session revoke dan audit diuji. Guard tanggung jawab ditambahkan bersama domain: task M2, review M3, event M4. | — | — |
| M1-08 | TODO | Shell UI dan auth screens | M1-06 | Login/change-password, sidebar desktop, nav mobile, settings/profile/users; dua sesi terpisah; 401/403 dan return path aman. | — | — |
| M1-09 | TODO | CI dan test foundation | M1-01, M1-03, M1-04 | Scripts quality nyata, Svelte/TS/build, gofmt/vet/test, PostgreSQL integration, OpenAPI drift dan secret scan berjalan. | — | — |
| M1-10 | TODO | Spike deployment dan pairing environment | M1-06, M1-09 | Dua project Vercel; nested Go routes/query/body, Cookie/Set-Cookie dan CSRF melalui proxy terbukti; preview DB/API terpisah, cron/push OFF. | — | — |
| M1-11 | TODO | Catat keputusan runtime dan biaya | M1-10 | ADR adapter/PWA/version/limits/pooling tersedia; paket bisnis dan kebutuhan cron diverifikasi, biaya/budget dicatat tanpa asumsi gratis. | — | — |
| M1-GATE | TODO | Acceptance milestone | Seluruh task M1 | Dua sesi beda role; Go/Chi Vercel dan cookie/CSRF proxy terbukti; DB preview terpisah; reset/last-supervisor lulus; keputusan biaya tercatat. | — | — |

### Bukti M1-01 — 29 Oktober 2026 WITA

Seluruh command memakai toolchain lokal yang diaktifkan sesuai [README](README.md). Source dan pilihan versi dijelaskan pada [ADR-001](architecture-decisions.md). Tidak ada commit, push, deploy, migration, akses DB, atau pekerjaan M1-02.

| Pemeriksaan | Command / metode aktual | Hasil dan batas verifikasi |
|---|---|---|
| Install reproducible | `pnpm install --frozen-lockfile --offline` | PASS; pin exact cocok lockfile, cache tersedia pada mesin ini. Clone baru memakai install frozen tanpa --offline. |
| Quality lokal | `pnpm quality` | PASS; Node/Go pins, svelte-check 0 error/0 warning, gofmt check, go vet ./..., Prettier/ESLint, Vitest, go test ./..., build web adapter Vercel dan binary Go. Vitest/go test belum memiliki unit test bisnis; tidak dihitung sebagai coverage domain. |
| Check setelah perapian akhir | `pnpm check:api`; `pnpm --filter @marketing-workspace/web exec prettier --check package.json` | PASS; gofmt memeriksa seluruh folder API, manifest web akhir sesuai format. |
| Integritas Go Modules | `go mod verify` dari apps/api | PASS; all modules verified. |
| Runtime lama ditolak | `/usr/local/bin/node scripts/check-tooling.mjs` pada Node 24.14.0 | PASS untuk uji negatif; keluar 1 dengan instruksi memakai 24.22.2. Build/check positif memakai Node 24.22.2. |
| Startup lokal | `pnpm dev:web --strictPort`, `pnpm dev:api`, lalu `pnpm test:smoke` | PASS; SSR web 200 berbahasa Indonesia pada 127.0.0.1:5173; server Go/Chi merespons 404 yang diharapkan pada 127.0.0.1:8080 karena belum ada endpoint produk. Memerlukan izin loopback di sandbox sesi ini. |
| Visual | Chrome desktop: buka web lokal dan periksa screenshot | PASS untuk halaman fondasi light mode, judul/copy Indonesia, tanpa data demo. Warning hydration favicon tercatat ketika href menjadi SVG dengan atribut data-codex-favicon-badge yang tidak berasal dari source proyek; bukan klaim runtime tanpa warning. Bukan UAT produk atau uji perangkat mobile. |
| Audit dependency frontend | `pnpm audit --audit-level high`; rincian via `pnpm audit --json` dan `pnpm view deepmerge version --json` | PASS untuk ambang high; 0 high/critical, 1 low pada dependency build, dilacak BUG-001. Patch belum diterbitkan di registry saat pemeriksaan. |
| Svelte autofixer / cloud / domain | Pemeriksaan ketersediaan tool; scope M1-01 | NOT RUN: autofixer tidak tersedia (diganti compiler/typecheck). CI, deployment, DB, auth, PWA/push dan domain belum diimplementasikan/diuji. |

### Bukti M1-02 — 29 Oktober 2026 WITA

PIC: Codex. State awal bersih; bukti M1-01 dan pin Node 24.22.2/pnpm 11.10.0/Go 1.27.3 diperiksa sebelum perubahan. Semua command memakai toolchain lokal dan cache proyek sesuai README.

| Pemeriksaan | Command / metode aktual | Hasil dan batas verifikasi |
|---|---|---|
| Quality repository | `pnpm quality` | PASS: tooling, Svelte/TypeScript 0 error/0 warning, format/lint, gofmt/vet, tests, build web dan Go. Vitest belum memiliki test; Go kini memiliki test pada empat package foundation, bukan coverage bisnis. |
| Backend akhir | `pnpm check:api`; `pnpm build:api`; `go test -race ./...` dari apps/api | PASS setelah perbaikan 405/Allow pada route nested; empat package Go lulus race detector. Percobaan wrapper `pnpm --dir apps/api exec go test -race ./...` gagal karena apps/api bukan package pnpm; command Go langsung dari apps/api lulus. |
| Acceptance HTTP | Test di apps/api/internal/app, internal/platform/httpx, api dan cmd/server | PASS config invalid/fail-closed, env mismatch, push/cron OFF, request ID tidak percaya input, 404/405 JSON termasuk Allow, nested route/query/body dan Cookie/Set-Cookie, Origin negatif, body limit/JSON ketat, panic tanpa bocor partial response, cancellation/timeout/late write, log tanpa data sensitif. Route produk belum dibuat; fixture test memakai data sintetis. |
| Smoke build akhir | `pnpm dev:web --strictPort`; binary `.cache/marketing-api` dengan env development eksplisit; `pnpm test:smoke` | PASS web SSR Indonesia, API JSON404/request ID/no-store dan POST tanpa Origin ditolak403 pada loopback. Server Go mencatat server_started setelah bind. Akses localhost membutuhkan izin sandbox sesi ini. |
| Dependency/pin | `git diff --exit-code -- package.json apps/web/package.json pnpm-lock.yaml pnpm-workspace.yaml apps/api/go.mod apps/api/go.sum .node-version .go-version .nvmrc` | PASS tidak berubah. BUG-001 tetap TODO; audit dependency tidak diulang dan patch tidak diasumsikan tersedia. |
| Dokumentasi/config | .env.example, README, ARCHITECTURE, PRD common HTTP errors dan ADR-002 | Env wajib/default/limits, pola router dan batas verifikasi dicatat; semua Markdown proyek tetap di root. |
| Di luar scope | Cloud, DB, auth, OpenAPI/client generation, CI, browser visual baru | NOT RUN: Vercel routing/proxy/cookie/environment pairing nyata menunggu M1-10; DB M1-03, OpenAPI M1-04 dan auth M1-06 belum dimulai. Tidak ada UI berubah; autofixer tidak tersedia, compiler/typecheck lulus. Tidak commit/push/deploy atau mengubah konfigurasi global mesin. |

## M2 — Task & Direct Completion

Sumber: PRD §24 dan bagian domain terkait; AGENTS serta architecture berlaku pada semua baris.

| ID | Status | Task | Dependensi | Kriteria selesai | PIC | Bukti / catatan |
|---|---|---|---|---|---|---|
| M2-01 | TODO | Schema task, child, history dan inbox awal | M1-GATE | Migration task/checklist/comments/links/activity serta inbox dan tabel pendukung sesuai kebutuhan writer; completion CHECK dasar, version/content_revision dan ledger idempotency diuji. Kolom/FK menuju tabel event/round/reminder ditunda sesuai PRD §13.5; inbox task-only. Konfirmasi OQ-02 (DEC-006) sebelum schema final. | — | — |
| M2-02 | TODO | Policy, transaction dan idempotency task | M2-01 | Creator/PIC/admin override sesuai matrix; version/revision dan replay key aman; rollback audit/inbox serta SQL parameter diuji. | — | — |
| M2-03 | TODO | Task CRUD, assignment, archive/restore dan guard PIC | M2-02 | Create/read/update/filter/pagination, PIC aktif, review OFF; archive flag; assignment inbox atomik, self-assignment tanpa inbox. Deaktivasi PIC task terbuka ditolak 409 ACTIVE_RESPONSIBILITIES; assignment-versus-deaktivasi diserialisasi melalui lock pengguna dan diuji. | — | — |
| M2-04 | TODO | Checklist, komentar dan links | M2-03 | Limits/reorder/URL tervalidasi; checklist/link memakai expected_version dan menaikkan task version/content_revision. Komentar append-only memakai Idempotency-Key dan izin, tanpa expected_version maupun kenaikan task version/content_revision; history terisi. Konfirmasi kebijakan komentar OQ-10 (DEC-010) sebelum UAT M2. | — | — |
| M2-05 | TODO | Start, return dan direct completion | M2-04 | Hanya PIC; todo→in_progress→done; checklist selesai; completed_by/waktu server/method direct dan completion CHECK dasar; transisi ilegal ditolak. Konfirmasi checklist OQ-05 (DEC-009) sebelum UAT M2. | — | — |
| M2-06 | TODO | Task list/detail dan Kanban | M2-05 | Search/filter, URL detail/drawer, form/conflict states; drag memakai aksi backend, alternatif keyboard; refresh konsisten. | — | — |
| M2-07 | TODO | Dashboard dan aktivitas awal | M2-03, M2-05 | Mine/team, counts/overdue/recent activities sesuai PRD; pagination; polling ≥60 detik hanya visible/online. | — | — |
| M2-08 | TODO | Acceptance direct dan race tests | M2-06, M2-07 | Alvin dan Sultan create→start→done; forbidden API, stale version, retry timeout/same key, concurrent mutation, inbox dedup serta race assignment/deaktivasi lulus; OQ-05/OQ-10 terkonfirmasi sebelum UAT. | — | — |
| M2-GATE | TODO | Acceptance milestone | Seluruh task M2 | Kedua PIC menjalankan direct completion; completion CHECK, permission/concurrency/idempotency, guard deaktivasi PIC dan assignment inbox tanpa push lulus; OQ-02/OQ-05/OQ-10 memiliki bukti keputusan. | — | — |

## M3 — Optional Review

Sumber: PRD §24 dan bagian domain terkait; AGENTS serta architecture berlaku pada semua baris.

| ID | Status | Task | Dependensi | Kriteria selesai | PIC | Bukti / catatan |
|---|---|---|---|---|---|---|
| M3-01 | TODO | Schema review rounds dan invariant | M2-GATE | Tambah rounds, FK dua arah komentar/keputusan dan notifications.review_id sesuai PRD §13.5; pending round unik, final round immutable. Validasi approved round melengkapi completion CHECK dasar M2; review_sequence/content_revision dan upgrade migration diuji. | — | — |
| M3-02 | TODO | Requires Review dan request review | M3-01 | Default OFF; ON wajib reviewer aktif; voluntary review tetap boleh; PIC/pengaju bukan reviewer, checklist lengkap dan request inbox atomik. | — | — |
| M3-03 | TODO | Approve dan request changes | M3-02 | Actor/round/version/snapshot tervalidasi; approval method review-approved; perubahan wajib komentar dan inbox; jabatan tidak membatasi review. | — | — |
| M3-04 | TODO | Cancel, resubmit, reopen dan reviewer change | M3-03 | Alasan wajib; freeze Review/Done; cancel sebelum edit; resubmit round baru; reopen clear completion kini tanpa menghapus history. | — | — |
| M3-05 | TODO | Policy change dan administrative override | M3-04 | OFF←ON hanya creator/admin dengan alasan; no force approve; archive Review ditolak; aksi/allowed_actions konsisten. Konfirmasi kebijakan OQ-04 (DEC-008) sebelum UAT M3. | — | — |
| M3-06 | TODO | Perluasan guard deaktivasi untuk review | M3-04 | Guard PIC task terbuka dari M2 tetap lulus; review pending memblok deaktivasi sampai cancel/pemindahan eksplisit. Session/subscription revocation aman; race request-review/deaktivasi diuji. | — | — |
| M3-07 | TODO | Review UI dan history | M3-05 | Tandai Selesai/Minta Review/Setujui/Minta Revisi sesuai konteks; round history, stale notification dan draft conflict ditangani. | — | — |
| M3-08 | TODO | Acceptance review dan concurrency | M3-06, M3-07 | Sultan→Alvin dan Alvin→Sultan; self-approval/ON direct/stale round ditolak; approve-cancel dan submit bersamaan; completion history lulus. | — | — |
| M3-GATE | TODO | Acceptance milestone | Seluruh task M3 | Review dua arah, ON/OFF, self-approval, stale round, cancel race, serta history completion lulus; OQ-04 memiliki bukti keputusan sebelum UAT M3. | — | — |

## M4 — Marketing & Dual Calendar

Sumber: PRD §24 dan bagian domain terkait; AGENTS serta architecture berlaku pada semua baris.

| ID | Status | Task | Dependensi | Kriteria selesai | PIC | Bukti / catatan |
|---|---|---|---|---|---|---|
| M4-01 | TODO | Schema event dan event API | M3-GATE | Calendar_events beserta tasks.event_id/audit_logs.event_id/notifications.event_id dan FK/index; CHECK inbox diperluas menjadi tepat satu task/event. All-day/timed union, exclusive end, category/scale/PIC, version/audit, archive/restore; relasi task tanpa cascade archive; upgrade migration diuji. | — | — |
| M4-02 | TODO | Linked tasks, progress dan guard akun | M4-01 | Linked summary hanya non-archived; tanpa task bukan 0%; event mendatang/berlangsung mencegah deaktivasi PIC sampai dialihkan. | — | — |
| M4-03 | TODO | Calendar read projection dan WITA | M4-02 | Task deadline tidak membuat event record; range ≤93 hari, overlap/filter, date-only/timed/overdue dan ID projection benar. | — | — |
| M4-04 | TODO | Provider Hijriyah dan data provenance | M4-03 | Referensi Indonesia diverifikasi dengan URL/hash/cakupan/fixtures; dataset di calendar-data/, provenance di CALENDAR-DATA.md pada root; fallback tabular berlabel; dataset runtime terpaket, bukan HTTP per request. OQ-06 (DEC-003) terjawab sebelum gate M4. | — | — |
| M4-05 | TODO | Event UI dan FullCalendar | M4-03, M4-04 | Month/week, filters, event detail/forms/progress; lifecycle destroy; WITA benar pada timezone perangkat berbeda, metadata Hijriyah terlihat. | — | — |
| M4-06 | TODO | Acceptance kalender | M4-05 | Date-only/all-day/timed, akhir bulan/tahun, arsip, linked task; fallback/dataset QA dan perubahan referensi tidak menggeser Gregorian lulus. | — | — |
| M4-GATE | TODO | Acceptance milestone | Seluruh task M4 | Event dan task projection tanpa duplikat; WITA/date union benar; provenance/fallback Hijriyah serta Gregorian tetap lulus. | — | — |

## M5 — Event Reminder

Sumber: PRD §24 dan bagian domain terkait; AGENTS serta architecture berlaku pada semua baris.

| ID | Status | Task | Dependensi | Kriteria selesai | PIC | Bukti / catatan |
|---|---|---|---|---|---|---|
| M5-01 | TODO | Schema dan generation reminder | M4-GATE | Event_reminders dan notifications.reminder_id beserta FK/index/upgrade migration; H-7/H-3/H-1 atau H-1 default 09.00 WITA; custom/rules=[]; masa lalu skipped; penerima creator+PIC dedup. OQ-03 (DEC-007) dikonfirmasi sebelum UAT M5. | — | — |
| M5-02 | TODO | Revision, reschedule dan cancellation | M5-01 | Jadwal/PIC/rules naikkan revision; cancel pending/delivery lama; scale custom tidak tertimpa; restore regenerasi eksplisit. | — | — |
| M5-03 | TODO | Cron authentication dan bounded dispatch | M5-02 | Bearer secret/flags/no-store; event→reminder locks; maksimum 100 kandidat/budget terverifikasi; catch-up ≤24 jam sebelum event mulai. | — | — |
| M5-04 | TODO | Atomic reminder inbox dan delivery writer | M5-03 | Unique dedupe key, issued+inbox+delivery commit bersama; crash sebelum/sesudah commit pulih; tanpa network I/O di transaksi. | — | — |
| M5-05 | TODO | Reminder UI dan dashboard | M5-04 | Edit rules, skipped/cancelled/scheduled states, upcoming reminder 7 hari; notifikasi snapshot membuka event terkini. | — | — |
| M5-06 | TODO | Scheduler environment dan acceptance | M5-05 | Cron nyata pada environment sesuai dan manual staging terotorisasi; secret, overlap/reschedule/archive/catch-up tests; preview tidak memicu production. OQ-03 terkonfirmasi sebelum UAT; OQ-09 (DEC-002) diputuskan sebelum cron M5 di production. | — | — |
| M5-GATE | TODO | Acceptance milestone | Seluruh task M5 | Cron nyata dan manual staging, secret, dedup, cancellation, catch-up dan race tests lulus; OQ-03 memiliki bukti keputusan sebelum UAT M5. | — | — |

## M6 — PWA & Push

Sumber: PRD §24 dan bagian domain terkait; AGENTS serta architecture berlaku pada semua baris.

| ID | Status | Task | Dependensi | Kriteria selesai | PIC | Bukti / catatan |
|---|---|---|---|---|---|---|
| M6-01 | TODO | Notification center dan preferences | M5-GATE | Inbox/read/unread count/filter/deep link, owner-only/read-all snapshot; preferences hanya mengatur push dan lima trigger PRD lengkap. | — | — |
| M6-02 | TODO | PWA manifest, icons dan satu worker | M6-01 | Manifest/icons/install, satu registrasi SW, cache asset publik saja; update prompt aman, offline notice dan save diblok. | — | — |
| M6-03 | TODO | Opt-in dan subscription lifecycle | M6-02 | Feature detection, denied tidak berulang, VAPID public/config; user+session/device ownership, transfer/logout/revoke; tanpa backlog perangkat baru. | — | — |
| M6-04 | TODO | Push transport dan keamanan endpoint | M6-03 | VAPID private backend; HTTPS/provider allowlist, destination/SSRF/no redirect; payload generik dan same-origin notificationclick. | — | — |
| M6-05 | TODO | Lease, retry dan recovery delivery | M6-04 | SKIP LOCKED, lease token/expiry, network setelah commit; lima percobaan/TTL/jitter/Retry-After; 404/410 disable, revalidasi relevansi/ownership. | — | — |
| M6-06 | TODO | Acceptance PWA/push/perangkat | M6-05 | Fake transport dan perangkat target; inbox bertahan saat gagal; opt-in/denied/logout, crash/retry, SW update/offline/private cache dan deep link diuji; batas perangkat dicatat. OQ-07 (DEC-005) diputuskan sebelum M6 selesai. | — | — |
| M6-GATE | TODO | Acceptance milestone | Seluruh task M6 | PWA/perangkat, opt-in/denied, inbox tahan kegagalan, logout, retry, invalid subscription, deep link dan SW update lulus. | — | — |

## M7 — Hardening & Production

Sumber: PRD §24 dan bagian domain terkait; AGENTS serta architecture berlaku pada semua baris.

| ID | Status | Task | Dependensi | Kriteria selesai | PIC | Bukti / catatan |
|---|---|---|---|---|---|---|
| M7-01 | TODO | Full contract dan requirement audit | M6-GATE | Seluruh endpoint/FR/US/DoD PRD dipetakan ke bukti; generated client tak drift; tidak ada placeholder atau scope liar. | — | — |
| M7-02 | TODO | Security dan cross-feature regressions | M7-01 | IDOR/escalation/CSRF/XSS/SQLi/SSRF/rate limit/revocation diuji; secrets tidak di repo/bundle/log; tidak ada vulnerability kritis terbuka. | — | — |
| M7-03 | TODO | Accessibility dan performance | M7-01 | Keyboard/focus/labels/status bukan warna saja; dataset 10k task/1k event diuji; target p95 warm <1s dan usable ≤3s dinilai, cold start terpisah. | — | — |
| M7-04 | TODO | Release workflow dan environment isolation | M7-02, M7-03 | Migration sekali/lock, backend→web→smoke, job activation; production auto-deploy tidak mendahului schema; secrets dan preview pairing terpisah. | — | — |
| M7-05 | TODO | Backup/restore dan rollback rehearsal | M7-04 | Restore database terpisah serta rollback app kompatibel dibuktikan; target RPO≤24h/RTO≤4h, retensi dan akses operasional disepakati. | — | — |
| M7-06 | TODO | Monitoring, maintenance dan budget | M7-04 | Cron >15m alert, errors/pool/retry/usage; cleanup aman dedup; paket bisnis/cron/backup, budget serta owner operasional dicatat. | — | — |
| M7-07 | TODO | Runbook, README dan UAT | M7-05, M7-06 | Setup/env/deploy/recovery/kalender dan keterbatasan terdokumentasi; Sultan/Alvin menjalankan skenario PRD; hasil penerimaan nyata dicatat. | — | — |
| M7-08 | TODO | Production rollout dan smoke | M7-07 | Gate produksi dan keputusan OQ-01/OQ-08/OQ-09 terpenuhi sebelum produksi/aktivasi cron berbayar sesuai tenggat PRD; deploy sesuai otorisasi, health/auth/read/write/cron/push smoke; versi release dan keputusan lanjut/rollback dicatat. | — | — |
| M7-GATE | TODO | Acceptance milestone | Seluruh task M7 | Semua AC kritis, biaya/sumber kalender, isolation, UAT, restore/rollback, release dan smoke produksi memiliki bukti. | — | — |

## Issue lokal: bug, blocker, dan keputusan

Tidak ada bug implementasi atau blocker yang telah diverifikasi pada baseline. Pertanyaan teknis dan produk di bawah adalah pekerjaan verifikasi/keputusan, bukan laporan kegagalan. Issue tidak menambah scope tanpa keputusan produk. PIC tetap belum ditetapkan; keputusan produk memerlukan konfirmasi pemilik produk sesuai PRD.

| ID | Status | Jenis | Ringkasan / keputusan yang dibutuhkan | Terkait | PIC | Resolusi / bukti |
|---|---|---|---|---|---|---|
| DEC-001 | TODO | Keputusan | Adapter, Go routing, cookie proxy, limits/region dan pasangan preview | M1-10, M1-11 | — | ADR-001 mencatat tooling; ADR-002 mencatat handler Go dan pengujian lokal M1-02. Spike Go/cloud, cookie proxy, limits dan preview nyata belum dibuktikan. |
| DEC-002 | TODO | Keputusan | OQ-01: paket bisnis, budget/pemilik akun dan biaya; OQ-09: interval cron dan keterlambatan reminder yang diterima | M1-11, M5-06, M7-06, M7-08 | — | Spike/biaya dicatat M1; OQ-01 sebelum production/cron berbayar; OQ-09 sebelum reminder M5 di production. Verifikasi aktual, tanpa nominal tebakan. |
| DEC-003 | TODO | Data | OQ-06: sumber/cakupan Hijriyah Indonesia dan fallback yang layak rilis | M4-04, M4-GATE, M7-07 | — | Putuskan pada gate M4; provenance/fixtures serta keterbatasan harus jelas. |
| DEC-004 | TODO | Operasional | OQ-08: retensi, akses recovery, backup storage serta target RPO/RTO | M7-05, M7-06, M7-08 | — | Default PRD dikonfirmasi dan dibuktikan sebelum production. |
| DEC-005 | TODO | Perangkat | OQ-07: matriks perangkat nyata untuk PWA/install/push | M6-06, M6-GATE | — | Putuskan sebelum M6 selesai; catat yang diuji dan yang belum dapat diverifikasi. |
| DEC-006 | TODO | Produk | OQ-02: kecukupan satu PIC/reviewer per task dan satu PIC per event | M2-01, M2-GATE | — | Default satu; konfirmasi sebelum schema M2 final, tanpa menganggap revisi dokumen sebagai persetujuan. |
| DEC-007 | TODO | Produk | OQ-03: penerima reminder dan jam default | M5-01, M5-06, M5-GATE | — | Default creator+PIC, 09.00 WITA; konfirmasi sebelum UAT M5. |
| DEC-008 | TODO | Produk | OQ-04: creator boleh mematikan Requires Review dengan alasan/history | M3-05, M3-GATE | — | Default ya; PIC non-creator tidak boleh. Konfirmasi sebelum UAT M3. |
| DEC-009 | TODO | Produk | OQ-05: checklist wajib selesai sebelum complete/submit, link tidak wajib | M2-05, M2-08, M2-GATE | — | Default semua checklist checked; konfirmasi sebelum UAT M2. |
| DEC-010 | TODO | Produk | OQ-10: seluruh pengguna aktif boleh komentar, append-only | M2-04, M2-08, M2-GATE | — | Konfirmasi sebelum UAT M2; tidak memberi izin edit/hapus komentar. |
| BUG-001 | TODO | Keamanan dependency | Expected: dependency build tanpa advisory deepmerge; actual: plugin Svelte membawa deepmerge 4.3.1, advisory low GHSA-gcqx-w48m-w2h4 (prototype poisoning). Reproduksi: pnpm audit --json; jalur melalui @sveltejs/vite-plugin-svelte 7.3.1. | M1-09, M7-02 | Codex | Pemeriksaan M1-01: registry latest 4.3.1; metadata audit menyebut rentang aman ≥4.3.2, tetapi advisory resmi menyatakan patched versions None. Tidak ada patch yang diasumsikan tersedia; audit tidak diulang pada M1-02. Update transitif saat patch tersedia dan ulang audit/check/build. Tidak ada high/critical; tidak mengubah dependency sendiri atau menyembunyikan temuan. Tidak menghalangi acceptance scaffold lokal M1-01. |

Issue baru memakai ID stabil `BUG-002`, `BLK-001`, atau `DEC-011` berikutnya; jangan daur ulang ID. Tautkan task terdampak. Bug menulis expected/actual, langkah reproduksi, dampak dan bukti. Blocker menulis penyebab/aksi pembuka/owner; keputusan menulis alternatif dan hasil di ADR. Status menggunakan vocabulary tracker yang sama; VERIFY untuk fix/keputusan yang belum dibuktikan. Jangan menyimpan credential, data produksi, atau screenshot sensitif dalam tracker.

Jika task dipecah, pertahankan ID induk dan gunakan suffix `M2-03a`/`M2-03b`, dengan status `induk → anak` dijelaskan. Hitung hanya leaf task pada ringkasan agar tidak double count; gate tetap tidak dihitung. Task baru harus punya milestone, sumber, acceptance dan dependensi. Pembatalan scope harus merujuk keputusan eksplisit, bukan dipakai untuk menghapus kegagalan acceptance.

## Log pembaruan

| Tanggal WITA | ID | Perubahan | Bukti / langkah berikutnya |
|---|---|---|---|
| 2026-10-28 | TRACKER | Baseline TODO.md dibuat dari tiga dokumen; implementasi belum diverifikasi | Mulai M1-01; seluruh task/gate implementasi TODO |
| 2026-10-29 | DOCS | Codex menyelaraskan tujuh dokumen: seluruh Markdown/rujukan di root sesuai permintaan pengguna; guard PIC M2/review M3/event M4, staging migration/FK/completion CHECK, pengecualian version komentar, dan pemetaan seluruh OQ ke DEC-002–DEC-010. | PASS pemeriksaan Python inline: 10 tautan lokal, tabel/fence/whitespace, 53 task/7 gate/dependensi dan 10 OQ terpetakan. PASS validasi YAML/frontmatter SKILL via Ruby/Psych; validator quick_validate.py tidak dapat berjalan karena PyYAML tidak tersedia. Implementasi tetap 0/53, seluruh gate dan keputusan produk tetap TODO; berikutnya M1-01. |
| 2026-10-29 | M1-01 | Codex menyelesaikan monorepo/tooling lokal: TODO → IN_PROGRESS → DONE. Node awal 24.14.0 diganti untuk eksekusi proyek dengan toolchain lokal 24.22.2 setelah verifikasi advisory resmi; Go 1.27.3 disiapkan lokal. | PASS quality/frozen install/smoke/Chrome desktop; unit bisnis belum ada. BUG-001 mencatat satu advisory low build. Progres 1/53 (1,9%); M1 1/11 (9,1%); 0/7 gate. M1-02 belum dimulai; tidak commit/push/deploy. |
| 2026-10-29 | M1-02 | Codex menyelesaikan boundary HTTP lokal: TODO → IN_PROGRESS → DONE; config fail-closed, router bersama, errors/request ID, timeout/cancellation dan log aman. | PASS quality, checks/build backend akhir, race dan smoke; README/env/ADR-002 diperbarui. Progres 2/53 (3,8%); M1 2/11 (18,2%); 0/7 gate. BUG-001 terbuka; M1-03 belum dimulai. Tidak commit/push/deploy/akses DB. |

Tambahkan satu baris ringkas per perubahan bermakna; status terkini tetap di tabel task. Log tidak menggantikan bukti. Saat bug membuka kembali pekerjaan DONE, ubah task terdampak ke TODO/IN_PROGRESS/BLOCKED sesuai kondisi, tautkan BUG, hitung ulang ringkasan, dan buka kembali gate bila acceptance milestone tidak lagi terpenuhi. Gate yang dibuka kembali menahan milestone lanjutan sampai dampaknya dinilai.
